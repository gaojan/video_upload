package controllers

import (
	"encoding/json"
	"fmt"
	"new_project/models"
	"new_project/utils"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type UploadController struct {
	beego.Controller
}

func (u *UploadController) Post() {

	adv := u.GetString("adv")
	if adv == "" {
		u.Ctx.WriteString("adv参数不能为空")
		return
	}

	// 获取设备信息
	var data = make(map[string]string)
	data["regadv"] = adv
	req := &utils.RequestInfo{
		Url:           beego.AppConfig.String("aut_url"),
		Data:          data,
		DataInterface: nil,
	}

	resp, err := req.PostUrlEncoded()
	if err != nil {
		u.Ctx.WriteString("获取adv错误")
		return
	}

	respBody := make([]map[string]interface{}, 1)
	json.Unmarshal(resp, &respBody)

	num := 0
	for _, val := range respBody {
		adv := val["adv"].(string)
		name := val["name"].(float64)
		// adv为空 或name为0 不能上传
		if adv == "" {
			u.Ctx.WriteString("adv为空，不能上传")
			return
		}

		// name 中有3 直接无限次
		if name == 3 {
			num = -1
			goto breakHere
		}

		//name为1和2 记录num次数 且总次数不能大于5次
		if name == 1 || name == 2 {
			advRecord, _ := models.GetAdvRecordByAdvName(adv, name)
			// 还没有该adv数据 初始化一次 并保存数据库
			if advRecord == nil {
				advRec := models.AdvRecord{
					Adv:      adv,
					Name:     name,
					Num:      4,
					CreateDt: utils.GetCurrentTime(),
				}
				if err := models.AddAdvRecord(&advRec); err != nil {
					u.Ctx.WriteString("Save Database Filed")
					return
				}
				// 第一次上传算一次 所以5次要减掉一次
				num += 4
			} else { // 数据库中有该adv 查询num次数
				if advRecord.Num <= 0 {
					u.Ctx.WriteString("没有权限不能上传")
					return
				} else {
					// 更新数据库记录
					num = advRecord.Num - 1
					if err := models.UpdateAdvRecordByAdvName(adv, name, num); err != nil {
						u.Ctx.WriteString("Update Database Filed")
						return
					}

				}
			}
		}

	}

breakHere:

	// 视频
	video, videoInfo, err := u.GetFile("video")
	if err != nil {
		u.Ctx.WriteString(fmt.Sprintf("%v", err))
		return
	}
	if videoInfo.Size > 20971520 {
		u.Ctx.WriteString("Video max size cannot exceed 20MB")
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
	// 新建文件夹
	basePath := utils.GetCurrentPath()
	fmt.Println(basePath)
	currentMonth := utils.GetCurrentMonth()
	dirPath := basePath + "/static/upload/" + currentMonth
	if err := utils.PathExists(dirPath); err != nil {
		fmt.Println(err)
	}

	// 保存视频
	newVideoName := utils.GetRandomString()
	videoPath := path.Join(dirPath, newVideoName+"."+layout)
	err = u.SaveToFile("video", videoPath)
	if err != nil {
		u.Ctx.WriteString("视频上传失败")
		return
	}
	defer video.Close()

	// 封面图片
	img, imgInfo, err := u.GetFile("img")
	if imgInfo == nil || err != nil {
		// 没有上传封面 使用默认封面
		defaultImgPath := basePath + "/static/default.jpg"
		imgPath := path.Join(dirPath, newVideoName+".jpg")
		w, err := utils.CopyFile(defaultImgPath, imgPath)
		fmt.Println(w)
		fmt.Println(err)
		uploadVideo := models.UploadRecord{
			ImgPath:   imgPath,
			VideoPath: videoPath,
			CreateDt:  utils.GetCurrentTime(),
		}

		// 保存数据库
		if err := models.AddUploadRecord(&uploadVideo); err != nil {
			u.Ctx.WriteString("Save Database Filed")
			return
		}

		respMap, err := utils.Struct2SSMap(uploadVideo)
		if err != nil {
			u.Ctx.WriteString("Struct to map error")
			return
		}
		respMap["num"] = num
		u.Data["json"] = respMap
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
		imgPath := path.Join(dirPath, newVideoName+"."+suffix)
		err := u.SaveToFile("img", imgPath)
		if err != nil {
			u.Ctx.WriteString("封面上传失败")
			return
		}
		uploadVideo := models.UploadRecord{
			ImgPath:   imgPath,
			VideoPath: videoPath,
			CreateDt:  utils.GetCurrentTime(),
		}

		// 保存数据库
		if err := models.AddUploadRecord(&uploadVideo); err != nil {
			u.Ctx.WriteString("Save Database Filed")
			return
		}

		defer img.Close()
		respMap, err := utils.Struct2SSMap(uploadVideo)
		if err != nil {
			u.Ctx.WriteString("Struct to map error")
			return
		}
		respMap["num"] = num
		u.Data["json"] = respMap
		u.ServeJSON()

	}
}
