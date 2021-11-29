package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	fmt.Println(1)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
