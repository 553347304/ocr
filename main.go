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
	Key    string `json:"key"`
	Base64 string `json:"base64"`
}

type Data struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}
type Result[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
}

func ocr(c *gin.Context) {
	var cr Request
	var list = make([]Data, 0)
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.JSON(200, Result[Data]{Code: 7, Message: "参数错误: " + string(body)})
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(cr.Base64)
	if err != nil {
		c.JSON(200, Result[Data]{Code: 7, Message: "base64编码错误: " + cr.Base64})
		return
	}

	tmpFile, _ := ioutil.TempFile("", "img-")
	tmpFile.Write(imageBytes)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	dir, _ := os.Getwd()

	switch cr.Model {
	case "":
		cr.Model = "v4"
	case "v4", "v3":
		break
	default:
		c.JSON(200, Result[Data]{Code: 7, Message: "不存在此模型版本"})
		return
	}

	cmd := exec.Command(
		//"wine", // Linux执行时
		path.Join(dir, "RapidOCR", "RapidOCR-json.exe"),
		"--models=models",
		fmt.Sprintf("--det=ch_PP-OCR%s_det_infer.onnx", cr.Model),
		"--cls=ch_ppocr_mobile_v2.0_cls_infer.onnx",
		fmt.Sprintf("--ch_PP-OCR%s_rec_infer.onnx", cr.Model),
		"--keys=ppocr_keys_v1.txt",
		"--image_path="+tmpFile.Name(),
	)
	cmd.Dir = "RapidOCR"
	output, err := cmd.Output()
	if err != nil {
		c.JSON(200, Result[Data]{Code: 7, Message: fmt.Sprint("内部路径错误", err)})
		return
	}
	re := regexp.MustCompile(`\{.*\}`)
	matches := re.FindAllString(string(output), -1)

	var response Response
	err = json.Unmarshal([]byte(matches[0]), &response)
	if err != nil {
		c.JSON(200, Result[Data]{Code: 7, Message: "内部错误"})
		return
	}

	//logs.Structs(response)

	var all []string
	for _, datum := range response.Data {
		data := Data{
			X:     datum.Box[0][0] + (datum.Box[2][0]-datum.Box[0][0])/2,
			Y:     datum.Box[0][1] + (datum.Box[2][1]-datum.Box[0][1])/2,
			Text:  datum.Text,
			Score: datum.Score,
		}
		all = append(all, datum.Text)
		if cr.Key != "" && cr.Key != "ALL" {
			// 返回指定文本
			if cr.Key == "TEXT" || strings.Contains(datum.Text, cr.Key) {
				c.JSON(200, Result[Data]{Code: 0, Message: "ok 模型:" + cr.Model, Data: []Data{data}})
				return
			}
		} else {
			list = append(list, data)
		}
	}
	if cr.Key == "ALL" {
		c.JSON(200, Result[string]{Code: 0, Message: "ok 模型:" + cr.Model, Data: all}) // 返回全部文本
		return
	}
	c.JSON(200, Result[Data]{Code: 0, Message: "ok 模型:" + cr.Model, Data: list}) // 返回全部
}

func main() {
	r := gin.Default()
	r.Use(middleware.Http().Timeout(10 * time.Second))
	r.POST("/", ocr)
	r.Run(":80")
}
