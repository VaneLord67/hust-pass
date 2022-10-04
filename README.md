# hust-pass

基于golang的selenium+colly+百度数字OCR的HUST统一身份登录脚本

# 配置文件示例
文件名: config.json，放置于根目录下，与main.go同一级
```json
{
  "ak": "百度数字识别appId",
  "sk": "百度数字识别appSecret",
  "chrome_driver_service_port": 8833,
  "username": "统一身份系统账号，一般是学号",
  "password": "统一身份系统密码",
  "chrome_driver_path": "chromedriver的路径",
  "caller_token": "调用OCR服务的token",
  "ocr_service_port": 8754,
  "AmMeter_ID": "电费查询接口的AmMeter_ID",
  "sms_app_id": "榛子云短信平台AppID",
  "sms_app_secret": "榛子云短信平台AppSecret",
  "sms_template_id": "榛子云短信平台短信模板ID",
  "phone_number": "电话号码",
  "room_id": "寝室号",
  "elec_threshold": 30.0
}
```

# Dockerfile部署
```shell
1. 编写配置文件，放置于根目录，与Dockerfile同级
2. 使用Go交叉编译，生成Linux平台可用的二进制可执行文件，命名为build_hust_pass_linux
   [Goland交叉编译教程](https://blog.csdn.net/a6652162/article/details/121084093)
3. docker build -t hust-pass . # 生成镜像
4. docker -itd --name hust-pass hust-pass # 运行容器
5. docker logs hust-pass # 查看日志
```

# Docker镜像部署
```shell
1. docker pull vanelord/hust-pass:v1
2. docker run -itd --name hust-pass vanelord/hust-pass:v1
3. 编写config.json，使用docker cp命令将config.json放到容器根目录下
4. docker restart hust-pass
5. docker logs -f hust-pass
```