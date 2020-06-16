package main

import (
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func initFlags() *Config {
	cfg := &Config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.CertFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.KeyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.Port, "port", "8080", "Port to start webhook on. Default 8080")
	fl.StringVar(&cfg.Severity, "severity", "CRITICAL", "Severity level to check in images")
	fl.StringVar(&cfg.IgnoreFile, "ignorefile", "/mnt/trivyignore", "Ignore file to be passed to enforcer")
	fl.Parse(os.Args[1:])
	return cfg
}

func main() {

	c := initFlags()

	go func() {
		for {
			dbUpdate()
			time.Sleep(12 * time.Hour)
		}
	}()

	logrus.Info("Booting up Webhook...")
	if err := c.Serve(); err != nil {
		logrus.Error(err)
	}

}
