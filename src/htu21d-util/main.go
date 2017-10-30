package main

import (
	"fmt"
	"os"
	"github.com/idahoakl/HTU21D-sensor/src/cmd"
	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
