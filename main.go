package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"ocr/middleware"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Response struct {
	Code int `json:"code"`
	Data []struct {
		Box   [][]int `json:"box"`
		Score float64 `json:"score"`
		Text  string  `json:"text"`
	} `json:"data"`
}

type Request struct {
	Model  string `json:"model"`
	Dict   string `json:"dict"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Base64 string `json:"base64"`
	Dir    string `json:"dir"`
}

type Data struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Result struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Item    []string `json:"item"`
	Data    []Data   `json:"data"`
}

func InList(imageExtensions []string, ext string) bool {
	ext = strings.ToLower(ext)
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// go build -ldflags="-H windowsgui" main.go
// RapidOCR-json.exe --models=models --image_path=1.png
func ocr(c *gin.Context) {
	var cr Request
	var data = make([]Data, 0)
	var item = make([]string, 0)
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.JSON(200, Result{Code: 7, Message: "参数错误: " + string(body)})
		return
	}

	ext := path.Ext(cr.Base64)
	is := InList([]string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp"}, ext)

	imagePath := cr.Base64
	if !is {
		imageBytes, err := base64.StdEncoding.DecodeString(cr.Base64)
		if err != nil {
			c.JSON(200, Result{Code: 7, Message: "base64编码错误: " + cr.Base64})
			return
		}

		tmpFile, _ := ioutil.TempFile("", "img-")
		tmpFile.Write(imageBytes)
		tmpFile.Close()
		defer os.Remove(tmpFile.Name())
		imagePath = tmpFile.Name()
	}

	dir, _ := os.Getwd()

	switch cr.Model {
	case "":
		cr.Model = "v4"
	case "v4", "v3":
		break
	default:
		c.JSON(200, Result{Code: 7, Message: "不存在此模型版本"})
		return
	}

	if cr.Dict == "" {
		cr.Dict = "ppocr_keys_v1.txt"
	}

	cmd := exec.Command(
		//"wine", // Linux执行时
		path.Join(dir, "RapidOCR", "RapidOCR-json.exe"),
		"--models=models",
		"--cls=ch_ppocr_mobile_v2.0_cls_infer.onnx",
		fmt.Sprintf("--det=ch_PP-OCR%s_det_infer.onnx", cr.Model),
		fmt.Sprintf("--ch_PP-OCR%s_rec_infer.onnx", cr.Model),
		fmt.Sprintf("--keys=%s", cr.Dict),
		"--image_path="+imagePath,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Dir = "RapidOCR"
	output, err := cmd.Output()
	if err != nil {
		c.JSON(200, Result{Code: 7, Message: fmt.Sprint("内部路径错误", err.Error())})
		return
	}
	re := regexp.MustCompile(`\{.*\}`)
	matches := re.FindAllString(string(output), -1)
	var response Response
	err = json.Unmarshal([]byte(matches[0]), &response)
	//logs.Structs(response)
	if err == nil {
		for _, datum := range response.Data {

			item = append(item, datum.Text)
			data = append(data, Data{
				X: cr.X + datum.Box[0][0] + (datum.Box[2][0]-datum.Box[0][0])/2,
				Y: cr.Y + datum.Box[0][1] + (datum.Box[2][1]-datum.Box[0][1])/2,
			})
		}
	}
	c.JSON(200, Result{Code: 0, Message: "model:" + cr.Model, Item: item, Data: data}) // 返回全部
}

func main() {
	r := gin.Default()
	r.Use(middleware.Http().Timeout(10 * time.Second))
	r.POST("/", ocr)
	r.Run(":80")
}
