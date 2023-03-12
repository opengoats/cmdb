# cmdb

多云资产管理平台


## 快速开发
make脚手架
```sh
➜  cmdb git:(master) ✗ make help
dep                            Get the dependencies
lint                           Lint Golang files
vet                            Run go vet
test                           Run unittests
test-coverage                  Run tests with coverage
build                          Local build
linux                          Linux build
run                            Run Server
clean                          Remove previous build
help                           Display this help screen
```

1. 使用安装依赖的Protobuf库(文件)
```sh
# 把依赖的probuf文件复制到/usr/local/include

# 创建protobuf文件目录
$ make -pv /usr/local/include/github.com/opengoats/goat/pb

# 找到最新的goat protobuf文件
$ ls `go env GOPATH`/pkg/mod/github.com/opengoats/

# 复制到/usr/local/include
$ cp -rf pb  /usr/local/include/github.com/opengoats/goat/pb
```

2. 添加配置文件(默认读取位置: etc/cmdb.toml)
```sh
$ 编辑样例配置文件 etc/cmdb.toml.book
$ mv etc/cmdb.toml.book etc/cmdb.toml
```

3. 启动服务
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
