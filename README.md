# 图床API

提供一个上传图片到Github仓库的API

通过下面的步骤来运行这个程序 
```bash
# 填写配置项 
$ cp .env.exmaple .env # 然后修改必要的配置项 

# 安装依赖  
$ go get .  # 仅安装依赖，不更新依赖包  
$ go get -u . # 安装依赖，并且更新依赖包 
# -u 标志表示更新所有的依赖项到最新的次要版本或修补版本（minor or patch version）

# 运行程序 
$ go run main.go 

# 访问服务状态 
$ curl http://localhost:8000/health

# 打包编译程序 
$ go build main.go 
$ ./main
```


# 部署 
```sh
# 创建网络 
docker network create --driver bridge --subnet=172.20.0.0/16 --gateway=172.20.0.1 docker20

# mysql 数据库允许 172.20.% 访问 
# redis bind 127.0.0.1 172.20.0.1 
# 防火墙设置 172.17.0.0/12 允许访问 6379,3306 端口 

$ cp .env.production .env 
# 然后修改 数据库密码等等 

# 启动 
$ make production 

# 关闭
$ make clean 
```