package https

import "time"

const timeout = 10 // 超时时间

type Request struct {
	URL      string            `json:"url"`      // URL
	Header   map[string]string `json:"header"`   // Header
	Request  any               `json:"request"`  // 传入数据
	Response any               `json:"response"` // 绑定结构体指针
	Time     time.Duration     `json:"time"`     // 超时时间
}

func (r *Request) Param() {
	// 默认超时时间
	if r.Time == 0 {
		r.Time = timeout
	}
}

type Response struct {
	Code int    `json:"code"`
	Body string `json:"body"`
}
