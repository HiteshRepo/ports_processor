package db

type DatabaseConfig struct {
	Dbname        string         `mapstructure:"name"`
	Username      string         `mapstructure:"user"`
	Password      string         `mapstructure:"password"`
	Host          string         `mapstructure:"host"`
	Schema        string         `mapstructure:"schema"`
	Port          int            `mapstructure:"port"`
	LogMode       bool           `mapstructure:"logMode"`
	SslMode       string         `mapstructure:"sslMode"`
	Connection    ConnectionPool `mapstructure:"connectionPool"`
}

type ConnectionPool struct {
	MaxOpenConnections int `mapstructure:"maxOpenConnections"`
	MaxIdleConnections int `mapstructure:"maxIdleConnections"`
	MaxIdleTime        int `mapstructure:"maxIdleTime"`
	MaxLifeTime        int `mapstructure:"maxLifeTime"`
	TimeOut            int `mapstructure:"timeout"`
}
