package models

import (
	"github.com/astaxie/beego/orm"
)

type UploadRecord struct {
	Id        int    `orm:"column(id)" json:"-"`
	Title     string `orm:"column(title)" json:"title"`
	ImgPath   string `orm:"column(img_path)" json:"img_path"`
	VideoPath string `orm:"column(video_path)" json:"video_path"`
	CreateDt  string `orm:"column(create_dt)" json:"create_dt"`
}

type AdvRecord struct {
	Id  int    `orm:"column(id)" json:"-"`
	Adv string `orm:"column(adv)" json:"adv"`
	//Name     float64 `orm:"column(name)" json:"name"`
	Num      int    `orm:"column(num)" json:"num"`
	CreateDt string `orm:"column(create_dt)" json:"create_dt" `
}

func (t *UploadRecord) TableName() string {
	return "upload_record"
}

func (t *AdvRecord) TableName() string {
	return "adv_record"
}

func AddUploadRecord(video *UploadRecord) error {
	o := orm.NewOrm()
	_, err := o.Insert(video)
	return err
}

func AddAdvRecord(adv *AdvRecord) error {
	o := orm.NewOrm()
	_, err := o.Insert(adv)
	return err
}

//func GetAdvRecordByAdvName(adv string, name float64) (*AdvRecord, error) {
//	o := orm.NewOrm()
//
//	advRecord := &AdvRecord{Adv: adv, Name: name}
//	err := o.Read(advRecord, "adv", "name")
//	if err != nil {
//		return nil, err
//	}
//	return advRecord, err
//}

func GetAdvRecordByAdvName(adv string) (*AdvRecord, error) {
	o := orm.NewOrm()

	advRecord := &AdvRecord{Adv: adv}
	err := o.Read(advRecord, "adv")
	if err != nil {
		return nil, err
	}
	return advRecord, err
}

//func UpdateAdvRecordByAdvName(adv string, name float64, num int) error {
//	o := orm.NewOrm()
//	advRec := new(AdvRecord)
//	_, err := o.QueryTable(advRec).Filter("adv", adv).Filter("name", name).Update(orm.Params{"num": num})
//	if err != nil {
//		return err
//	}
//	return nil
//}

func UpdateAdvRecordByAdvName(adv string, num int) error {
	o := orm.NewOrm()
	advRec := new(AdvRecord)
	_, err := o.QueryTable(advRec).Filter("adv", adv).Update(orm.Params{"num": num})
	if err != nil {
		return err
	}
	return nil
}
