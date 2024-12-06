package main

import (
	"github.com/gin-gonic/gin"
	"ocr/https"
	"ocr/logs"
	"strings"
)

type Request struct {
	Key    string `json:"key"`
	Base64 string `json:"base64"`
}
type Response struct {
	Code int `json:"code"`
	Data []struct {
		Box   [][]int `json:"box"`
		Score float64 `json:"score"`
		Text  string  `json:"text"`
		End   string  `json:"end"`
	} `json:"data"`
	Score     float64 `json:"score"`
	Time      float64 `json:"time"`
	Timestamp float64 `json:"timestamp"`
}

type Data struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}
type Result[T any] struct {
	Code int     `json:"code"`
	Time float64 `json:"time"`
	Data []T     `json:"data"`
}

func ocr(c *gin.Context) {
	var cr Request
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		c.JSON(200, "参数错误")
		return
	}

	var response Response
	https.Post(https.Request{
		URL:     "http://127.0.0.1:1224/api/ocr",
		Request: cr,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Response: &response,
	})

	logs.Structs(response)

	var list = make([]Data, 0)
	var stringList []string
	for _, datum := range response.Data {
		x := datum.Box[0][0] + (datum.Box[2][0]-datum.Box[0][0])/2
		y := datum.Box[0][1] + (datum.Box[2][1]-datum.Box[0][1])/2
		data := Data{
			X:     x,
			Y:     y,
			Text:  datum.Text,
			Score: datum.Score,
		}
		list = append(list, data)
		stringList = append(stringList, datum.Text)
		if cr.Key != "" {
			// 返回找到的第一个
			if cr.Key == "TEXT" {
				break
			}
			// 返回指定文本
			if strings.Contains(datum.Text, cr.Key) {
				c.JSON(200, Result[Data]{
					Code: response.Code,
					Time: response.Time,
					Data: []Data{data},
				})
				return
			}
		}
	}

	// 返回找到的第一个
	if cr.Key == "ALL" {
		c.JSON(200, Result[string]{
			Code: response.Code,
			Time: response.Time,
			Data: stringList,
		})
		return
	}

	c.JSON(200, Result[Data]{
		Code: response.Code,
		Time: response.Time,
		Data: list,
	})
}

func main() {
	r := gin.Default()
	r.POST("/ocr", ocr)
	r.Run(":8080")
}
