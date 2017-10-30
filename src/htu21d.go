package htu21d

import (
	"sync"
	"github.com/idahoakl/go-i2c"
	"time"
	"github.com/sigurn/crc8"
	log "github.com/Sirupsen/logrus"
)

/*
http://www.te.com/usa-en/models/4/00/028/790/CAT-HSC0004.html

Command 						Code 	Comment
Trigger Temperature Measurement 0xE3 	Hold master
Trigger Humidity Measurement 	0xE5 	Hold master
Trigger Temperature Measurement 0xF3 	No Hold master
Trigger Humidity Measurement 	0xF5 	No Hold master
Write user register 			0xE6
Read user register 				0xE7
Soft Reset 						0xFE
*/

type command struct {
	code  byte
	delay time.Duration
}

var (
	tempMeasure     = command{code: byte(0xF3), delay: 50 * time.Millisecond}
	humidityMeasure = command{code: byte(0xF5), delay: 16 * time.Millisecond}
	writeUserReg    = command{code: byte(0xE6), delay: 0 * time.Millisecond}
	readUserReg     = command{code: byte(0xE7), delay: 0 * time.Millisecond}
	softReset       = command{code: byte(0xFE), delay: 15 * time.Millisecond}
)

var crc8Table = crc8.MakeTable(crc8.Params{
	Poly:   0x31,
	Init:   0x0,
	RefIn:  false,
	RefOut: false,
	XorOut: 0x0,
	Check:  0x0,
})

type HTU21D struct {
	Connection *i2c.I2C
	Address    uint8
	Mtx        sync.Mutex
}

// Create new instance for provided i2c bus and sensor address
func New(address uint8, connection *i2c.I2C) (*HTU21D, error) {
	sensor := &HTU21D{
		Connection: connection,
		Address:    address,
	}

	if e := sensor.Reset(); e != nil {
		return nil, e
	}

	return sensor, nil
}

// Soft reset sensor
func (sensor *HTU21D) Reset() error {
	sensor.Mtx.Lock()
	defer sensor.Mtx.Unlock()

	log.Debug("Sensor reset")

	if _, e := sensor.Connection.Write(sensor.Address, []byte{softReset.code}); e != nil {
		return e
	}

	time.Sleep(softReset.delay)

	return nil
}

// Read temperature in degrees C
func (sensor *HTU21D) ReadTemperatureC() (float32, error) {
	sensor.Mtx.Lock()
	defer sensor.Mtx.Unlock()

	if data, e := sensor.readAndValidate(tempMeasure); e != nil {
		return -1, e
	} else {
		var temp = float32(data)
		temp *= 175.72
		temp /= 65536
		temp -= 46.85

		return temp, nil
	}
}

// Read temperature in degrees F
func (sensor *HTU21D) ReadTemperatureF() (float32, error) {
	if v, e := sensor.ReadTemperatureC(); e != nil {
		return -1, e
	} else {
		return (v * 1.8) + 32, nil
	}
}

// Read percent humidity
func (sensor *HTU21D) ReadHumidity() (float32, error) {
	sensor.Mtx.Lock()
	defer sensor.Mtx.Unlock()

	if data, e := sensor.readAndValidate(humidityMeasure); e != nil {
		return -1, e
	} else {
		var humidity = float32(data)
		humidity *= 125
		humidity /= 65536
		humidity -= 6

		return humidity, nil
	}
}

func validateCrc(data []byte, expectedCrc byte) error {
	checksum := crc8.Checksum(data, crc8Table)

	if expectedCrc != checksum {
		return &crcError{expected: expectedCrc, actual: checksum}
	}
	return nil
}

// Initiate operation, read result and perform CRC validation
func (sensor *HTU21D) readAndValidate(cmd command) (uint16, error) {
	if _, e := sensor.Connection.Write(sensor.Address, []byte{cmd.code}); e != nil {
		return 0, e
	}

	time.Sleep(cmd.delay)

	rawData := make([]byte, 3)

	if _, e := sensor.Connection.Read(sensor.Address, rawData); e != nil {
		return 0, e
	} else {
		log.WithFields(
			log.Fields{
				"rawData": rawData,
				"cmdCode": cmd.code,
			}).Debug("Data returned from HTU21D")

		if e := validateCrc(rawData[:2], rawData[2]); e != nil {
			return 0, e
		}

		var data = uint16(rawData[0])
		data <<= 8
		data |= uint16(rawData[1])

		return data, nil
	}
}
