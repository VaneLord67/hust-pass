# hust-pass

基于golang的selenium+colly+百度数字OCR的HUST统一身份登录脚本

# 配置文件示例
```json
{
  "ak": "百度数字识别appId",
  "sk": "百度数字识别appSecret",
  "chrome_driver_service_port": 8833,
  "username": "统一身份系统账号，一般是学号",
  "password": "统一身份系统密码",
  "chrome_driver_path": "chromedriver的路径",
  "caller_token": "调用OCR服务的token",
  "ocr_service_port": 8754
}
```