FROM selenium/standalone-chrome:latest
COPY build_hust_pass_linux docker/config.json /
WORKDIR /
RUN sudo mkdir "spider" && sudo cp /bin/chromedriver /spider/chromedriver.exe
CMD sudo ./build_hust_pass_linux