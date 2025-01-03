## 项目目录

- `dist` 打包后的文件，直接使用即可. `RapidOCR模型要和main.exe`放在一起
    - `windows` 可执行文件 main.exe

### RapidOCR模型下载地址     [https://vip.123pan.cn/1821560246/11196279](https://vip.123pan.cn/1821560246/11196279)

### RapidOCR官网     [https://github.com/hiroi-sora/RapidOCR-json?tab=readme-ov-file](https://github.com/hiroi-sora/RapidOCR-json?tab=readme-ov-file)

### 图片在线转base64   [https://www.jyshare.com/front-end/59/](https://www.jyshare.com/front-end/59/)

## 请求示例

| URL                   | 请求                     |
|-----------------------|------------------------|
| `http://127.0.0.1`    | POST `80端口`            |
| `http://ocr.tcbyj.cn` | 测试地址`指不定哪天就失联了 要用自己部署` |

| 请求参数     | Body   | 类型     | 默认值     | 可选参数                | 描述                               |
|----------|--------|--------|---------|---------------------|----------------------------------|
| `model`  | `json` | string | v4      | `v3` `v4`           | PP模型版本                           |
| `key`    | `json` | string | "" 返回全部 | `指定文本` `ALL` `TEXT` | TEXT返回第一个找到的文本 <br/>ALL返回全部文本    |
| `base64` | `json` | string | ""      | 无                   | iV开头的base64字符串  `iVBORw0KGgo...` |

## 返回值

- `code`: 状态码0为正常返回  `请求超过10秒自动超时`
- `message`：响应信息
- `data`内容为数组。数组每一项为字典，含四个元素： 当Key=ALL返回string数组
    - `x`：找到文本中心点x坐标，整数。
    - `y`：找到文本中心点y坐标，整数。
    - `text`：文本内容，字符串。
    - `score`：识别置信度，浮点数。

## 返回示例

```json
{
  "code": 0,
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

key为`ALL`时

```json
{
  "code": 0,
  "data": [
    ""
  ]
}
```

### Auto.js调用例子

```js
let img = captureScreen();
let result = http.postJson("http://ocr.tcbyj.cn", {
    "model": "v4",
    "key": "",
    "base64": images.toBase64(img),
}).body.json();
console.log(result);
```