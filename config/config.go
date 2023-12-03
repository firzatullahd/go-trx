package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	"strconv"
)

var searchPath = []string{
	".",
}

func initialiseFileAndEnv(v *viper.Viper, configName string) error {
	v.SetConfigName(configName)
	v.SetConfigType("yaml")

	for _, path := range searchPath {
		v.AddConfigPath(path)
	}

	v.AutomaticEnv()

	return v.ReadInConfig()
}

func Load(ConfigName string) (conf Config) {
	v := viper.New()

	err := initialiseFileAndEnv(v, ConfigName)
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No '" + ConfigName + "' file found on search paths. Will either use environment variables or defaults")
		} else {
			log.Fatalf("Error occured during loading config: %s", err.Error())
		}
	}

	err = v.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Error occured during unrmashal config: %s", err.Error())
	}

	return
}

func (p PSQL) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		p.Host, strconv.Itoa(p.Port), p.User, p.Password, p.DBName, p.Schema, "disable")
}
