```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o main main.go
```
```shell
docker build . -t <태그명:버전>
docker run -d  --name <컨테이너명> <태그명:버전>
docker logs -f <컨테이너명>
```

```shell
docker images

REPOSITORY                        TAG       IMAGE ID       CREATED          SIZE
태그명                              2.0       76a670597d92   2 seconds ago    1.22MB <--- 엄청 작아졌다!
태그명                              1.0       b289971a9452   16 minutes ago   798MB
```
