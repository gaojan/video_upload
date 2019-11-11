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

	advNo := u.GetString("adv_no")
	if advNo == "" {
		u.Ctx.WriteString("adv参数不能为空")
		return
	}

	// 获取设备信息
	var data = make(map[string]string)
	data["regadv"] = advNo
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

	// 解析返回值
	respBody := make([]map[string]interface{}, 1)
	json.Unmarshal(resp, &respBody)

	num := 0
	var adv string
	var nameList []float64
	for _, val := range respBody {
		adv = val["adv"].(string)
		name := val["name"].(float64)
		nameList = append(nameList, name)
		// adv为空 或name为0 不能上传
		if adv == "" {
			u.Ctx.WriteString("adv为空，不能上传")
			return
		}
	}

	//查询记录是否存在
	advRecord, _ := models.GetAdvRecordByAdvName(adv)
	// 3在name中 无限次
	if utils.SearchSlice(nameList, 3) {
		// 该条数据不存在 就保存
		num = -1
		if advRecord == nil {
			advRec := models.AdvRecord{
				Adv:      adv,
				Num:      -1,
				CreateDt: utils.GetCurrentTime(),
			}
			if err := models.AddAdvRecord(&advRec); err != nil {
				u.Ctx.WriteString("Save Database Filed")
				return
			}
		}
	} else {
		// 没有3 就是1或2  该条数据不存在就保存
		n := len(nameList)
		if advRecord == nil {
			num = n*5 - 1
			advRec := models.AdvRecord{
				Adv:      adv,
				Num:      num,
				CreateDt: utils.GetCurrentTime(),
			}
			if err := models.AddAdvRecord(&advRec); err != nil {
				u.Ctx.WriteString("Save Database Filed")
				return
			}
		} else {
			// 存在就更新数据
			num = advRecord.Num - 1
			if err := models.UpdateAdvRecordByAdvName(adv, num); err != nil {
				u.Ctx.WriteString("Update Database Filed")
				return
			}
		}
	}

	// 上传视频
	video, videoInfo, err := u.GetFile("video")
	if err != nil {
		u.Ctx.WriteString(fmt.Sprintf("%v", err))
		return
	}
	maxSize, _ := beego.AppConfig.Int64("maxSize")
	if videoInfo.Size > maxSize {
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
	newVideoName := utils.GetRandomString(16)
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
		// 没有上传封面 使用默认封面 defaultImg

		defaultImgPath := basePath + beego.AppConfig.String("defaultImg")
		imgPath := path.Join(dirPath, newVideoName+".jpg")
		if _, err := utils.CopyFile(defaultImgPath, imgPath); err != nil {
			u.Ctx.WriteString("默认封面不存在")
			return
		}
		uploadVideo := models.UploadRecord{
			ImgPath:   strings.Split(imgPath, "/static")[1],
			VideoPath: strings.Split(videoPath, "/static")[1],
			Key:       utils.GetRandomString(64),
			CreateDt:  utils.GetCurrentTime(),
		}

		// 保存上传记录
		if err := models.AddUploadRecord(&uploadVideo); err != nil {
			u.Ctx.WriteString("Save Database Filed")
			return
		}

		//组装返回参数
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

		// 保存上传记录
		uploadVideo := models.UploadRecord{
			ImgPath:   "/static" + strings.Split(imgPath, "/static")[1],
			VideoPath: "/static" + strings.Split(videoPath, "/static")[1],
			Key:       utils.GetRandomString(64),
			CreateDt:  utils.GetCurrentTime(),
		}
		if err := models.AddUploadRecord(&uploadVideo); err != nil {
			u.Ctx.WriteString("Save Database Filed")
			return
		}

		defer img.Close()
		// 组装返回参数
		respMap, err := utils.Struct2SSMap(uploadVideo)
		if err != nil {
			u.Ctx.WriteString("Struct to map error")
			return
		}

		respMap["video_path"] = "/snsvideodownload?filekey=" + respMap["key"].(string)
		respMap["num"] = num
		u.Data["json"] = respMap
		u.ServeJSON()
	}
}

type VideoDownloadController struct {
	beego.Controller
}

func (u *VideoDownloadController) Get() {
	fileKey := u.GetString("filekey")
	if fileKey == "" {
		u.Ctx.WriteString("filekey参数不能为空")
		return
	}
	video, err := models.GetUploadRecordByKey(fileKey)
	if err != nil {
		u.Ctx.WriteString("Video doesn't exist")
		return
	}
	videoPath := u.Ctx.Request.Host + video.VideoPath
	html := fmt.Sprintf("<html><head><title>Title</title><video src='http://%s' controls='controls'></video></head><body></body></html>", videoPath)
	u.Ctx.WriteString(html)
	return
}

type AdvController struct {
	beego.Controller
}

func (u *AdvController) Post() {
	adv := u.GetString("adv")
	if adv == "" {
		u.Ctx.WriteString("adv参数不能为空")
		return
	}
}
