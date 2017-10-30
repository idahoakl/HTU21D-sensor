package cmd

import (
	"github.com/spf13/cobra"
	"github.com/idahoakl/HTU21D-sensor/src"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func init() {
	RootCmd.AddCommand(humidityCmd)
}

var humidityCmd = &cobra.Command{
	Use: "humidity",
	Run: func(cmd *cobra.Command, args []string) {
		if sensor, e := htu21d.New(address, createi2c()); e != nil {
			log.Fatal(e)
		} else {
			if t, e := sensor.ReadHumidity(); e != nil {
				log.Fatal(e)
			} else {
				fmt.Printf("%f%%\n", t)
			}
		}
	},
}
