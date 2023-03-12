package all

// 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序
import (
	_ "github.com/opengoats/cmdb/apps/book/impl"
)
