FROM spryker/chromedriver:latest
COPY build_hust_pass_linux docker/config.json /home/webdriver/
WORKDIR /home/webdriver
RUN mkdir "spider" && cp /usr/bin/chromedriver ./spider/chromedriver.exe
ENTRYPOINT ./build_hust_pass_linux