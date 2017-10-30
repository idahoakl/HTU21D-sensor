package cmd

import (
	"github.com/spf13/cobra"
	"github.com/idahoakl/HTU21D-sensor/src"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

func init() {
	temperatureCmd.Flags().BoolVarP(&celsius, "celsius", "c", false, "output in Celsius")
	RootCmd.AddCommand(temperatureCmd)
}

var celsius bool

var temperatureCmd = &cobra.Command{
	Use: "temp",
	Run: func(cmd *cobra.Command, args []string) {
		if sensor, e := htu21d.New(address, createi2c()); e != nil {
			log.Fatal(e)
		} else {
			if(celsius) {
				if t, e := sensor.ReadTemperatureC(); e != nil {
					log.Fatal(e)
				} else {
					fmt.Printf("%f C\n", t)
				}
			} else {
				if t, e := sensor.ReadTemperatureF(); e != nil {
					log.Fatal(e)
				} else {
					fmt.Printf("%f F\n", t)
				}
			}
		}
	},
}
