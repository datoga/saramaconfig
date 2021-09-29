package main

import (
	"github.com/datoga/saramaconfig"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg, err := saramaconfig.NewFromViper(viper.GetViper())

	if err != nil {
		panic(err)
	}

	spew.Dump(cfg)
}
