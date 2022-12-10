//use https://github.com/spf13/viper
package config

import (
	"time"

	"github.com/spf13/viper"
)

var (
	// AppPath is the absolute path to the app
	AppPath string
	// appConfigPath is the path to the config files
	appConfigPath string
	//AppConfig app.yaml配置对象
	AppConfig *AppConf
	//CacheConfig cache.yaml配置对象
	CacheConfig *CacheConf
	//DatabaseConfig database.yaml配置对象
	DatabaseConfig *DataBaseConf
	//LogConfig log.yaml配置对象
	LogConfig *LogConf
	// 存放每个配置文件的viper对象
	configs = make(map[string]*viper.Viper, 0)
	//specialConfigFileList 固定配置文件名,4个配置文件
	specialConfigFileList = make(map[string]interface{}, 4)
)

//AppConf app.yaml struct
type AppConf struct {
	AppName    string //应用名称
	Debug      bool   //是否调试模式
	HttpAddr   string //监听地址
	HttpPort   int    //监听端口号
	ServerName string //在请求的时候输出 server 为 该字段值
}

//CacheConf cache.yaml
type CacheConf struct {
	Redis        map[string]RedisConf
	RedisCluster map[string]RedisClusterConf
}

//DataBaseConf database.yaml
type DataBaseConf struct {
	Mysql map[string]MysqlConf
}

//LogConf log.yaml
type LogConf struct {
	AccessLog             string        //AccessLog 访问日志存放路径，off关闭
	AccessLogFormat       string        //AccessLogFormat 访问日志格式
	AccessLogTimeInterval time.Duration //TimeInterval 更新时间戳周期
	ErrorLog              string        //服务错误日志以及程序运行过程调用logger模块方法打印日志都会存放在此文件,默认存在应用根目录下data/log/error.log
	ErrorLogFullCaller    bool          //错误日志打印完整调用栈，默认为false
	CmdErrorLog           string        //命令行日志文件，替换模板：{%app%}-应用名、{%process%}-进程名
	ServerLog             string        //服务启动/关闭/panic日志存放,默认存在应用根目录下data/log/server.log
	ServerLogLevel        string        //server log 日志级别,支持 debug/info/warn/error/dpanic/panic/fatal 共7种日志级别,级别从左往右为从小到大,默认为warn
	TimeZone              string        //TimeZone can be specified, such as "UTC" and "America/New_York" and "Asia/Chongqing", etc ,Optional. Default: "Local"
	TimeFormat            string        //TimeFormat 时间格式，默认为2006-01-02 15:04:05 ,针对errorlog和accesslog生效
	ErrorLogLevel         string        `mapstructure:"level"` //日志等级,支持 debug/info/warn/error/dpanic/panic/fatal 共7种日志级别,级别从左往右为从小到大,用于控制error log
	MaxSize               int           //日志轮转配置, 单位MB,表示最大文件大,超出则会新生成一个日志文件,默认为100MB
	MaxAge                int           //日志轮转配置, 文件最多保存多少天,单位天,默认不移除
	MaxBackups            int           //日志轮转配置, 日志文件最多保存多少个备份,默认保留所有日志文件
	Compress              bool          //日志轮转配置, 决定是否压缩日志文件存放
	OutputStdout          int           //设置输出到标准输出的日志，二进制标识，从左到右：第一位标识accessLog,第二位标识errorLog，第三位标识serverLog，0关1开,仅调试用，线上建议关闭
}

//MysqlConf mysql config struct
type MysqlConf struct {
	//SQL driver configs ,see https://github.com/go-sql-driver/mysql
	Host              string        //数据库host/ip
	DataBase          string        //数据库
	Username          string        //数据库用户名
	Password          string        //数据库密码，密码为空则填空即可
	Charset           string        //字符集
	Location          string        //地区
	Port              int           //数据库端口号
	MaxIdleConns      int           //设置空闲连接池中连接的最大数量 gorm使用database/sql包维护连接池
	MaxOpenConns      int           //设置打开数据库连接的最大数量
	ConnMaxLifetime   int           //设置连接可复用的最大时间，即链接的最长生命时间,单位秒
	MaxAllowedPacket  int           //Max packet size allowed in bytes. 设置客户端包大小，当maxAllowedPacket=0时，会设置为4MB；当maxAllowedPacket=-1时自动使用服务端的maxAllowedPacket配置
	ConnMaxIdleTime   time.Duration //设置链接空闲最大等待时间
	WriteTimeout      time.Duration //写超时时间，0代表不限制
	ReadTimeout       time.Duration //读超时时间，0代表不限制
	Timeout           time.Duration //连接超时
	ParseTime         bool          //正确的处理 time.Time，需要设置为true
	ColumnsWithAlias  bool          //是否允许返回别名，false时如select u.id会返回id
	InterpolateParams bool          //https://github.com/go-sql-driver/mysql#interpolateparams This can not be used together with the multibyte encodings BIG5, CP932, GB2312, GBK or SJIS. These are rejected as they may introduce a SQL injection vulnerability!

	//GORM configs see https://gorm.io/zh_CN/docs/gorm_config.html
	SlowThreshold             int    //慢查询阈值，单位毫秒
	LogLevel                  string //日志等级，会自动打印sql调试
	TablePrefix               string //表前缀
	SkipDefaultTransaction    bool   //是否跳过默认事务，为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除），默认为true
	SingularTable             bool   //使用单数表名，启用该选项，`User` 的表名应该是 `t_user`,而非`users`，默认为复数
	NoLowerCase               bool   //snake_casing of names,默认为FALSE
	DisableAutomaticPing      bool   //GORM 会自动 ping 数据库以检查数据库的可用性，若要禁用该特性，可将其设置为 true
	AllowGlobalUpdate         bool   //启用全局 update/delete
	IgnoreRecordNotFoundError bool   //日志输出是否忽略记录未找到错误，默认打印
	PrepareStmt               bool   //为true时执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
}

//RedisConf cache.yaml redis struct. see redis.Options{}
type RedisConf struct {
	Addr               string //节点地址加端口，如127.0.0.1:6379
	Password           string //密码,无则为空
	UserName           string //账号,无则填空
	Db                 int    //选择db
	MaxRetries         int    //命令执行失败时，最多重试多少次，默认为3次,-1表示不重试
	PoolSize           int    //连接池大小，默认值为10*CPU个数
	MinIdleConns       int    //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
	MaxConnAge         int    //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
	DialTimeout        int    //连接建立超时时间，默认5秒
	ReadTimeout        int    //socket读取超时时间，-1为不限制超时，0为默认值，默认为3s,单位为秒
	WriteTimeout       int    //socket写超时时间，默认值跟readtimeout一致
	PoolTimeout        int    //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
	IdleTimeout        int    //闲置超时，默认5分钟，-1表示取消闲置超时检查
	IdleCheckFrequency int    //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理
}

//RedisClusterConf cache.yaml redis cluster struct. see redis.ClusterOptions{}
type RedisClusterConf struct {
	Addrs              []string //集群节点地址 ip:port
	Password           string   //集群密码
	UserName           string   //集群账号
	MaxRetries         int      //命令执行失败时，最多重试多少次，默认为0即不重试
	RouteByLatency     bool     //默认false,为true则ReadOnly自动置为true,表示在处理只读命令时，可以在一个slot对应的主节点和所有从节点中选取Ping()的响应时长最短的一个节点来读数据
	RouteRandomly      bool     //默认false,为true则ReadOnly自动置为true,表示在处理只读命令时，可以在一个slot对应的主节点和所有从节点中随机挑选一个节点来读数据
	PoolSize           int      //连接池最大socket连接数，默认为5倍CPU数， 5 * runtime.NumCPU
	MinIdleConns       int      //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；
	MaxConnAge         int      //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
	DialTimeout        int      //连接建立超时时间，默认5秒
	ReadTimeout        int      //读超时，默认3秒， -1表示取消读超时
	WriteTimeout       int      //写超时，默认等于读超时，-1表示取消读超时
	PoolTimeout        int      //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
	IdleTimeout        int      //闲置超时，默认5分钟，-1表示取消闲置超时检查
	IdleCheckFrequency int      //闲置连接检查的周期，无默认值，由ClusterClient统一对所管理的redis.Client进行闲置连接检查。初始化时传递-1给redis.Client表示redis.Client自己不用做周期性检查，只在客户端获取连接时对闲置连接进行处理
}

type ZipkinTraceConf struct {
	URL   string
	Heads map[string]string
}

type JaegerEndpointTraceConf struct {
	URL      string
	Username string
	Password string
}

type JaegerAgentTraceConf struct {
	Host                        string
	Port                        string
	AttemptReconnectingInterval time.Duration
	MaxPacketSize               int
	DisableAttemptReconnecting  bool
}
