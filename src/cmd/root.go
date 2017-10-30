package cmd

import(
	"github.com/spf13/cobra"
	"log"
	"github.com/idahoakl/go-i2c"
)

var bus int
var address uint8

var RootCmd = &cobra.Command{
	Use: "htu21d-util",
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&bus, "bus", "b", 1, "i2c bus")
	RootCmd.PersistentFlags().Uint8VarP(&address, "address", "a", 0x40, "device address")
}

func createi2c() *i2c.I2C {
	if conn, e := i2c.NewI2C(bus); e != nil {
		log.Fatal(e)
		return nil
	} else {
		return conn
	}
}