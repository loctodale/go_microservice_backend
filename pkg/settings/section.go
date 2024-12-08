package settings

type Config struct {
	Server     ServerSetting `mapstructure:"server"`
	Mysql      MySQLSetting  `mapstructure:"mysql"`
	Redis      RedisSetting  `mapstructure:"redis"`
	SMTP       SMTPSetting   `mapstructure:"smtp"`
	JWTSetting JWTSetting    `mapstructure:"jwt"`
	MailTrap   MailTrap      `mapstructure:"mailtrap"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type MySQLSetting struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DbName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
}

type SMTPSetting struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type JWTSetting struct {
	TokenHourLifeSpan string `mapstructure:"TOKEN_HOUR_LIFESPAN"`
	JWTExpiration     string `mapstructure:"JWT_EXPIRATION"`
	APISecret         string `mapstructure:"API_SECRET"`
}

type MailTrap struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}
