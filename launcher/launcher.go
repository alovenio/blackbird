// launcher application initializes and starts Blackbird's SFU
package main

import (
	"alovenio.com/blackbird/logger"
	"alovenio.com/blackbird/sfu"
	"flag"
	"log"
)

var address = flag.String("address", "localhost:8000", "server address")
var logLevel = flag.String("logLevel", "info", "log level (debug, info, warn, error)")

func main() {
	flag.Parse()
	logLevel, err := logger.ParseLogLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger.LogLevel = logLevel
	server := new(sfu.Server)
	handler := new(sfu.MockSessionHandler)
	if err = server.Start(*address, handler); err != nil {
		logger.LogFatalF(err)
	}
}
