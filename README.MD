## Run!

``` bash
docker build -t fetch-app .

docker run -it --rm fetch-app https://autify.com/ https://baidu.com

docker run -it --rm fetch-app --metadata https://autify.com/ https://baidu.com
```