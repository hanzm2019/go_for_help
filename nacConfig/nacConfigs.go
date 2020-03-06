package nacConfig

//radius服务配置
type NacConfig struct {
	CiscoReauthenticateAVPair string //radius返回AC的重认证常量
	CiscoLastAVPair           string //radius返回AC的重认证常量
	//SrcIpAddr                 string //nac安装的服务器ip:41340
	//DisconnectDstIpAddr       string //断开AC开关连接无线的配置:(无线AC的ip:3799)
	ReauthDstIpAddr string //重连接AC开关连接无线的配置:(无线AC的IP:1700)
	AcIp            string
	Secret          string //radius服务器与AC之间数据包传递使用的共享密码，之后多个wifi做成动态的
	//ServerIpAddr              string //nac服务器所在机器ip端口配置，一般设置为0.0.0.0:1812
	RedirectUrl   string //AC发送终端连网请求后，AC告诉终端重定向的url路径，进行认证，通过连接wifi
	PreAuthAcl    string
	PermissionAcl string
	LetGoMAcs     string
}

//mysql数据库连接配置
type MySqlDb struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

//redis连接配置
type RedisDb struct {
	DbBase   int
	Protocol string
	Ip       string
	Port     string
	Auth     string
}

//web服务配置
type Login struct {
	LoginCasServerUrl  string //CAS认证登录请求路径，返回CAS登录界面的html
	LogoutCasServerUrl string //CAS登出请求路径
	LoginSuccSign      string //CAS登录成功表示符
	LogoutSuccSign     string //CAS登出成功表示符
	LoginRedirect      string //认证提交路径
	ReLoginTimes       int    //重复登录次数
}

//日志路径配置
type LogPath struct {
	LogDir      string //日志输出目录
	LogFileName string //日志输出文件+YYYY-MM-dd.log
}

type TimeConfig struct {
	ExpireTime int
	OnlineTime int
	MaxIdle    int
	MaxActive  int
}

type Config struct {
	Env        string
	Nac        NacConfig
	MysqlDb    MySqlDb
	Redisdb    RedisDb
	Login      Login
	LogPath    LogPath
	TimeConfig TimeConfig
}
