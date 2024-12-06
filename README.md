# 项目目录

- `ocr` ocr RapidOCR 版本模型
- `main.go` 源文件
- `dist` 打包后的文件，直接使用即可。
    - `linux`   可执行文件 main
    - `windows` 可执行文件 main.exe

## 请求示例

### 后端服务需要开放8080端口

### umi-ocr项目源码   [https://github.com/hiroi-sora/RapidOCR-json](https://github.com/hiroi-sora/RapidOCR-json)

### 图片在线转base64   [https://www.jyshare.com/front-end/59/](https://www.jyshare.com/front-end/59/)

| URL                            | 请求                     |
|--------------------------------|------------------------|
| `http://127.0.0.1:8080/`       | POST                   |
| `http://ocr.tcbyj.cn:8080/ocr` | 测试地址`指不定哪天就失联了 要用自己部署` |

| 请求参数     | 类型     | 默认值  | 可选参数                | 描述                               |
|----------|--------|------|---------------------|----------------------------------|
| `Key`    | string | 返回全部 | `指定文本` `ALL` `TEXT` | TEXT返回第一个找到的文本 <br/>ALL返回全部文本    |
| `base64` | string | ""   | 无                   | iV开头的base64字符串  `iVBORw0KGgo...` |

data内容为数组。数组每一项为字典，含三个元素：
text ：文本内容，字符串。
box ：文本包围盒，长度为4的数组，分别为左上角、右上角、右下角、左下角的[x,y]。整数。
score ：识别置信度，浮点数。

- `code`: 状态码100为正常返回
- `time`：识别时间
- `data`内容为数组。数组每一项为字典，含四个元素： 当Key=ALL返回string数组
    - `x`：找到文本中心点x坐标，整数。
    - `y`：找到文本中心点y坐标，整数。
    - `text`：文本内容，字符串。
    - `score`：识别置信度，浮点数。

```json
{
  "code": 100,
  "time": 1.7801780700683594,
  "data": [
	{
	  "x": 104,
	  "y": 504,
	  "text": "text",
	  "score": 0.27060410380363464
	}
  ]
}
```