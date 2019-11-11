package utils

import (
	"encoding/json"
	f "fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"time"
)

// 获取当前时间
func GetCurrentTime() string {
	UnixTime := time.Now().Unix()
	dataTimeStr := time.Unix(UnixTime, 0).Format("2006-01-02 15:04:05")
	return dataTimeStr
}

// 结构图转字典
func Struct2SSMap(s interface{}) (map[string]interface{}, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	dict := make(map[string]interface{})
	err = json.Unmarshal(j, &dict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}

// 获取项目当前路径
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// 获取当前月份
func GetCurrentMonth() string {
	UnixTime := time.Now().Unix()
	dataTimeStr := time.Unix(UnixTime, 0).Format("20060102 15:04:05")

	return dataTimeStr[0:6]
}

// 判断文件夹是否存在
func PathExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	e := os.MkdirAll(path, os.ModePerm)
	if e != nil {
		return e
	}
	return nil
}

// 生成随机字符串16位
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 复制图片为默认封面
func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		f.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		f.Println(err)
	}
	defer desFile.Close()
	return io.Copy(desFile, srcFile)
}

func SearchSlice(s []float64, x float64) bool {
	a := sort.SearchFloat64s(s, x)
	if a > len(s)-1 {
		return false
	}
	return true
}

//
//func main()  {
//	//fmt.Println(GetCurrentMonth())
//	err := PathExists("./upload/hello")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//}
