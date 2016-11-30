# max-server 局域网游戏的消息服务器

## 功能说明
  - 配置玩家数目，并等待连接
  - 处理指定的消息，如选人 
  - 对于非指定的消息，进行广播
  - 处理了粘包、分包等通信问题

## 基本架构
  基于Golang进行开发,通信协议采用Protobuf

## 如何搭建开发环境
  - 安装Go(>= v1.7rc6)
  - 安装Protobuf
  ```sh
    go get github.com/golang/protobuf/proto
  ```
  - 安装Config模块
  ```sh
    go get github.com/robfig/config
  ```
## 如何测试

  ```sh
    go run server.go
    go run client.go
  ```
## 如何发布
  ```sh
    go build server.go
  ```

Copyright (c) 2016 KeenVision.cn All Rights Reserved.
