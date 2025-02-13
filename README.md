## 项目目录

- `dist` 打包后的文件，直接使用即可. `RapidOCR模型要和main.exe`放在一起
    - `windows` 可执行文件 main.exe
    - `linux` 可执行文件 main  `超级慢比window慢10倍+`

### RapidOCR模型下载地址     [https://vip.123pan.cn/1821560246/%E7%83%AD%E6%9B%B4/%E5%85%B6%E4%BB%96/model/RapidOCR.7z](https://vip.123pan.cn/1821560246/%E7%83%AD%E6%9B%B4/%E5%85%B6%E4%BB%96/model/RapidOCR.7z)

### RapidOCR官网     [https://github.com/RapidAI/RapidOCR](https://github.com/RapidAI/RapidOCR)

### 图片在线转base64   [https://www.jyshare.com/front-end/59/](https://www.jyshare.com/front-end/59/)

### 部署目录

``` yaml
- main.go
- RapidOCR
  - models
  - RapidOCR-json.exe
```

### 模型路径不能存在中文

## 请求示例

| URL                   | 请求                     |
|-----------------------|------------------------|
| `http://127.0.0.1`    | POST `80端口`            |
| `http://ocr.tcbyj.cn` | 测试地址`指不定哪天就失联了 要用自己部署` |

| 请求参数     | Body   | 类型     | 默认值  | 可选参数                              | 描述                   |
|----------|--------|--------|------|-----------------------------------|----------------------|
| `model`  | `json` | string | v4   | `v3` `v4`                         | PP模型版本               |
| `dict`   | `json` | string | 中英综合 | `ppocr_keys_v1.txt` `dict_en.txt` | dictTXT文件            |
| `x`      | `json` | int    | 0    |                                   | 起始点 x                |
| `y`      | `json` | int    | 0    |                                   | 起始点 y                |
| `base64` | `json` | string | ""   | `一个base64或者部署在本地,一个本地图片路径`        | 示例  `iVBORw0KGgo...` |

## 返回值

- `code`: 状态码0为正常返回  `请求超过10秒自动超时`
- `message`：响应信息

## 返回示例

```json
{
  "code": 0,
  "message": "model:v4",
  "item": [
    "10"
  ],
  "data": [
    {
      "x": 12,
      "y": 14
    }
  ]
}
```

[Auto.js调用例子](/调用例子/Auto.js调用例子.md)<br>
[按键精灵调用例子](/调用例子/按键精灵调用例子.md)<br>
[c++本地RapidOCR例子](/调用例子/c++本地RapidOCR例子.cpp)<br>