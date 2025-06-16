

# 原项目 https://github.com/lanyulei/ferry

## 基于Gin + Vue + Element UI前后端分离的工单系统

文档: [https://www.fdevops.com/docs/ferry](https://www.fdevops.com/docs/ferry-tutorial-document/introduction)

# 开发：
```aiexclude
git clone https://github.com/lanyulei/ferry.git
cd ferry
nano /config/settings.dev.yml #修改 
go get           # 安装依赖

<SQL>create database ferry charset 'utf8mb4';

go run main.go init -c config/settings.dev.yml   # 初始化数据结构

go run main.go server -c config/settings.dev.yml #启动

```

## 部署编译：
```aiexclude
cd ferry 
env GOOS=linux GOARCH=amd64 go build   #更多交叉编译内容，请访问 https://fdevops.com/2020/03/08/go-locale-configuration

## 修改 config/settings.yml

mkdir -p log static/uploadfile static/scripts static/template

将本地项目下 static/template 目录下的所有文件上传的到，服务器对应的项目目录下 static/template

<SQL> create database ferry charset 'utf8mb4';

./ferry init -c config/settings.yml

nohup ./ferry server -c=config/settings.yml > /dev/null 2>&1 & 
```
## git老提交不上：
```
git add .
git commit -m "记录"
git push origin main
```
## License

[MIT](https://github.com/lanyulei/ferry/blob/master/LICENSE)

Copyright (c) 2024 lanyulei


## 补充安装方法
[1.部署docker.md](1.%E9%83%A8%E7%BD%B2docker.md)  
[2.部署数据库.MD](2.%E9%83%A8%E7%BD%B2%E6%95%B0%E6%8D%AE%E5%BA%93.MD)   
[3.部署ferry.md](3.%E9%83%A8%E7%BD%B2ferry.md)


