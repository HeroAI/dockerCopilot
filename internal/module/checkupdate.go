package module

import (
	"encoding/json"
	"github.com/onlyLTY/oneKeyUpdate/UGREEN/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"strings"
)

// 检查更新处理后的镜像列表
type ImageCheckList struct {
	NeedUpdate bool
}
type ImageUpdateData struct {
	Data map[string]ImageCheckList
}

func NewImageCheck() *ImageUpdateData {
	return &ImageUpdateData{
		Data: map[string]ImageCheckList{},
	}
}
func (i *ImageUpdateData) CheckUpdate(imageList []types.Image) {
	for _, images := range imageList {
		imagename := removeProxy(images.Image_Name)
		baseURL := os.Getenv("hubURL")
		r, err := http.Get(baseURL + "/v2/repositories/" + imagename +
			"/tags/" + images.Image_Tag)
		if err != nil || r.StatusCode != 200 {
			logx.Error("获取远程镜像信息失败" + images.Image_Name + ":" + images.Image_Tag)
			continue
		}
		defer r.Body.Close()
		hubimage := types.HubImageInfo{}
		err = json.NewDecoder(r.Body).Decode(&hubimage)
		if err != nil {
			logx.Error("解析远程镜像信息失败" + images.Image_Name + ":" + images.Image_Tag)
			continue
		}
		remoteSHA256 := hubimage.Digest
		localSHA256 := strings.Split(images.RepoDigests[0], "@")[1]
		if remoteSHA256 != localSHA256 {
			logx.Info(images.Image_Name + ":" + images.Image_Tag + " need update")
			i.Data[images.ID] = ImageCheckList{NeedUpdate: true}
		} else {
			logx.Info(images.Image_Name + ":" + images.Image_Tag + " not need update")
		}

	}

}

func removeProxy(imageName string) string {
	imageNames := strings.Split(imageName, "/")
	if len(imageNames) == 3 {
		//fmt.Println("image_name: " + imageNames[1] + "/" + imageNames[2])
		return imageNames[1] + "/" + imageNames[2]
	} else if len(imageNames) == 2 {
		//fmt.Println("image_name: " + imageNames[0] + "/" + imageNames[1])
		return imageNames[0] + "/" + imageNames[1]
	} else {
		//fmt.Println("image_name: " + imageNames[0])
		return imageNames[0]
	}
}
