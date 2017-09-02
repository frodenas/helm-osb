package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"

	"github.com/frodenas/helm-osb/broker"
	"github.com/frodenas/helm-osb/helm"
)

var (
	configFilePath = flag.String("config-file", "", "Location of the configuration file")
	listenAddress  = flag.String("listen-address", ":3000", "Address to listen on")

	logLevels = map[string]lager.LogLevel{
		"DEBUG": lager.DEBUG,
		"INFO":  lager.INFO,
		"ERROR": lager.ERROR,
		"FATAL": lager.FATAL,
	}
)

func buildLogger(logLevel string) lager.Logger {
	laggerLogLevel, ok := logLevels[strings.ToUpper(logLevel)]
	if !ok {
		log.Fatalf("Log level `%s` is invalid", logLevel)
	}

	logger := lager.NewLogger("helm-osb")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, laggerLogLevel))

	return logger
}

func main() {
	flag.Parse()

	config, err := LoadConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Error loading configuration file: %s", err)
	}

	logger := buildLogger(config.LogLevel)

	helmClient := helm.New(config.HelmConfig, logger)

	serviceBroker := broker.New(config.BrokerConfig, helmClient, logger)

	credentials := brokerapi.BrokerCredentials{
		Username: config.BrokerConfig.Username,
		Password: config.BrokerConfig.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)

	fmt.Println("Starting Kubernetes Helm Open Service Broker...")
	if config.BrokerConfig.TLSCertFile != "" && config.BrokerConfig.TLSKeyFile != "" {
		fmt.Println("Listening TLS on", *listenAddress)
		log.Fatal(http.ListenAndServeTLS(*listenAddress, config.BrokerConfig.TLSCertFile, config.BrokerConfig.TLSKeyFile, nil))
	} else {
		fmt.Println("Listening on", *listenAddress)
		log.Fatal(http.ListenAndServe(*listenAddress, nil))
	}
}
