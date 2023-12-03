package config

type Config struct {
	App      Application `mapstructure:"app"`
	Postgres Postgres    `mapstructure:"database"`
	Redis    Redis       `mapstructure:"redis"`
	Constant Constant    `mapstructure:"constant"`
}

type Application struct {
	Env  string `mapstructure:"env"`
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type Postgres struct {
	MaxIdleCons    int  `mapstructure:"maxIdleCons"`
	MaxOpenCons    int  `mapstructure:"maxOpenCons"`
	ConMaxIdleTime int  `mapstructure:"conMaxIdleTime"`
	ConMaxLifetime int  `mapstructure:"conMaxLifeTime"`
	Slave          PSQL `mapstructure:"slave"`
	Master         PSQL `mapstructure:"master"`
}

type PSQL struct {
	DBName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Debug    bool   `mapstructure:"debug"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Constant struct {
	TrxTTL int `mapstructure:"trxTtl"`
}
