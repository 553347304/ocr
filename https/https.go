package https

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func Post(r Request) string {
	r.Param()
	data, _ := json.Marshal(r.Request)
	body := strings.NewReader(string(data))
	request, err := http.NewRequest("POST", r.URL, body)
	if err != nil {
		return "request: " + err.Error()
	}
	for key, value := range r.Header {
		request.Header.Add(key, value)
	}
	// 超时时间
	client := http.Client{
		Timeout: time.Second * r.Time,
	}
	res, err := client.Do(request)
	if err != nil {
		return "client: " + err.Error()
	}

	// 转换JSON
	jsons, err := io.ReadAll(res.Body)
	if err != nil {
		return "json: " + err.Error()
	}
	if r.Response != nil {
		err = json.Unmarshal(jsons, &r.Response)
		if err != nil {
			return "绑定失败: " + err.Error()
		}
	}

	return string(jsons)
}
