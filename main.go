package main

import (
	"flag"
	"os"

	"github.com/ibrokethecloud/enforcer/pkg/enforcer"

	"github.com/sirupsen/logrus"
)

func initFlags() *enforcer.Config {
	cfg := &enforcer.Config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.CertFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.KeyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.Port, "port", "8080", "Port to start webhook on. Default 8080")
	fl.StringVar(&cfg.IgnoreFile, "ignorefile", "/mnt/trivyignore", "Ignore file to be passed to enforcer")
	fl.Parse(os.Args[1:])
	return cfg
}

func main() {

	c := initFlags()

	logrus.Info("Booting up Webhook...")
	if err := c.Serve(); err != nil {
		logrus.Error(err)
	}

}
