# Auto.js调用例子

``` javascript
let img = captureScreen();
let result = http.postJson("http://ocr.tcbyj.cn", {
    "model": "v4",
    "base64": images.toBase64(img),
}).body.json();
console.log(result);
```