package controllers

import (
	"fmt"
	"new_project/models"
	"new_project/utils"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

// Operations about Users
type UploadController struct {
	beego.Controller
}

func (u *UploadController) Post() {
	// 视频
	video, videoInfo, err := u.GetFile("video")
	if err != nil {
		u.Ctx.WriteString(fmt.Sprintf("%v", err))
		return
	}
	//验证视频后缀
	videoName := videoInfo.Filename
	v := strings.Split(videoName, ".")
	layout := strings.ToLower(v[len(v)-1])
	if layout != "mp4" && layout != "avi" && layout != "mkv" && layout != "wmv" && layout != "mov" && layout != "rm" && layout != "3gp" {
		u.Ctx.WriteString("请上传符合格式的视频(mp4、avi、mkv、wmv、mov、rm、3gp)")
		return
	}
	// 保存视频
	videoPath := path.Join("static/upload/video", videoName)
	err = u.SaveToFile("video", videoPath)
	if err != nil {
		u.Ctx.WriteString("视频上传失败")
	}
	defer video.Close()

	// 封面图片
	img, imgInfo, err := u.GetFile("img")
	if imgInfo == nil || err != nil {
		// 没有上传封面 使用默认封面
		fmt.Println("没有上传封面")
		imgPath := "static/upload/image/default.jpg"

		uploadVideo := models.UploadVideo{
			ImgPath:   imgPath,
			VideoPath: videoPath,
			CreateDt:  utils.GetCurrentTime(),
		}

		u.Data["json"] = uploadVideo
		u.ServeJSON()

	} else {
		// 验证后缀
		imgName := imgInfo.Filename
		i := strings.Split(imgName, ".")
		suffix := strings.ToLower(i[len(i)-1])
		if suffix != "jpg" && suffix != "png" && suffix != "gif" {
			u.Ctx.WriteString("请上传符合格式的封面(jpg、png、gif)")
			return
		}
		// 保存封面
		imgPath := path.Join("static/upload/image", imgName)
		err := u.SaveToFile("img", imgPath)
		if err != nil {
			u.Ctx.WriteString("封面上传失败")
		}
		uploadVideo := models.UploadVideo{
			ImgPath:   imgPath,
			VideoPath: videoPath,
			CreateDt:  utils.GetCurrentTime(),
		}

		defer img.Close()
		u.Data["json"] = uploadVideo
		u.ServeJSON()

	}
}
