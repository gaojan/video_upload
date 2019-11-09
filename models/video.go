package models

type UploadVideo struct {
	Id         int    `json:"-"`
	Title      string `json:"title"`
	ImgPath    string `json:"img_path"`
	VideoPath  string `json:"video_path"`
	CreateUser string `json:"create_user"`
	CreateDt   string `json:"create_dt"`
}

type PalyVideo struct {
	id int
}
