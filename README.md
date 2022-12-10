# api_server
api_server demo

目录结构
config -- 写项目的配置文件
controller -- 控制器层，验证提交的数据，将验证完成的数据传递给 service
crontab -- 定时任务
service -- 业务层，只完成业务逻辑的开发，不进行操作数据库
repository -- 数据库操作层，比如写，多表插入，多表查询等，不写业务代码
model -- 数据库的ORM
entity -- 写返回数据的结构体, controller 层方法参数验证的结构体
framework -- 框架代码,以后会抽到一个公共仓库
router -- 路由注册
proto -- 写 gRPC 的 *.pb.go 文件
router -- 写路由配置及路由的中间件（鉴权、日志、异常捕获）
utils -- 写项目通用工具类.