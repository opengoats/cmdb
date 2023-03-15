## cmdb

多云资产管理平台


## 开发环境

* grpc 环境准备

```
# 1.安装protoc编译器,  项目使用版本: v3.19.1
# 下载预编译包安装: https://github.com/protocolbuffers/protobuf/releases

# 2.protoc-gen-go go语言查询, 项目使用版本: v1.27.1   
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 3.安装protoc-gen-go-grpc插件, 项目使用版本: 1.1.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 4.安装自定义proto tag插件
go install github.com/favadi/protoc-go-inject-tag@latest
```

## 启动服务
```sh
# 编译protobuf文件, 生成代码
$ make gen
# 如果是MySQL, 执行SQL语句(docs/schema/tables.sql)
$ make init
# 下载项目的依赖
$ make dep
# 运行程序
$ make run
```
