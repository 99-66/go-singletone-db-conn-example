package repositories

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

func LoadConfig(path string) {
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", ","))
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No such config file.")
		} else {
			log.Printf("reading config error.")
		}
		panic(err)
	}
}
