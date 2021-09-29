package main

import (
	"encoding/base64"
	"os"

	"github.com/datoga/saramaconfig"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"

	_ "embed"
)

//go:embed ca.crt
var ca []byte

//go:embed client.pem
var clientPem []byte

//go:embed client.key
var clientKey []byte

func main() {
	os.Setenv("TLS_CA", base64.StdEncoding.EncodeToString(ca))
	os.Setenv("TLS_CLIENTPEM", base64.StdEncoding.EncodeToString(clientPem))
	os.Setenv("TLS_CLIENTKEY", base64.StdEncoding.EncodeToString(clientKey))
	os.Setenv("NET_TLS_ENABLE", "true")

	viper.AutomaticEnv()

	cfg, err := saramaconfig.NewFromViper(viper.GetViper())

	if err != nil {
		panic(err)
	}

	spew.Dump(cfg.Net.TLS)
}
