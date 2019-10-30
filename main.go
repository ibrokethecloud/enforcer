package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

func initFlags() *Config {
	cfg := &Config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.CertFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.KeyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.Port, "port", "8080", "Port to start webhook on. Default 8080")
	fl.BoolVar(&cfg.Prune, "prune-images", false, "Prune images after each validation")
	fl.StringVar(&cfg.Severity, "severity", "CRITICAL", "Severity level to check in images")
	fl.Parse(os.Args[1:])
	return cfg
}

func main() {

	c := initFlags()

	logrus.Info("Booting Webhook...")
	if err := c.Serve(); err != nil {
		logrus.Error(err)
	}

}
