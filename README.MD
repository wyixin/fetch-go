## Run!

``` bash
docker build -t fetch-app .

docker run -it --rm fetch-app https://autify.com/ https://baidu.com

docker run --mount type=bind,source=[YOUR LOCAL PATH],target=/app/pages -it --rm fetch-app https://autify.com/ https://baidu.com
```