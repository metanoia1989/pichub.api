#!/bin/bash

# 设置编译环境变量
export CGO_ENABLED=0
export GOOS=linux

# 设置目标目录，可以根据需要修改
TARGET_DIR="./dist"
APP_NAME="pichub-api"

echo "开始构建 ${APP_NAME}..."

# 创建目标目录
mkdir -p ${TARGET_DIR}

# 下载依赖
echo "下载依赖..."
go mod download

# 编译
echo "编译程序..."
go build -mod=readonly -v -o ${TARGET_DIR}/${APP_NAME} main.go

# 复制配置文件
echo "复制配置文件..."
cp .env ${TARGET_DIR}/.env

echo "构建完成！"
echo "程序位置: ${TARGET_DIR}/${APP_NAME}"
echo "配置文件: ${TARGET_DIR}/.env"

# 显示使用说明
echo ""
echo "运行方式："
echo "cd ${TARGET_DIR} && ./${APP_NAME}"

# 添加 PM2 部署命令
echo "正在使用 PM2 部署..."
if pm2 list | grep -q "pichub-api"; then
    pm2 reload pichub-api
else
    pm2 start ecosystem.config.js
fi

echo "部署完成！"
echo ""
echo "PM2 管理命令："
echo "查看状态: pm2 status"
echo "查看日志: pm2 logs pichub-api"
echo "停止服务: pm2 stop pichub-api"
echo "重启服务: pm2 restart pichub-api" 