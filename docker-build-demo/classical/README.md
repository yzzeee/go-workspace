* 
```
GO_ENABLED=0 go build -o main .
```

* 이미지 생성 및 실행 확인 
```shell
docker build . -t demo:classic
# docker build . -t <태그명:버전>

docker run -d  --name demo-classic demo:classic
# docker run -d  --name <컨테이너명> <태그명:버전>

docker logs -f demo-classic
# docker logs -f <컨테이너명>
```