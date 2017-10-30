package htu21d

import "fmt"

type crcError struct {
	expected uint8
	actual   uint8
}

func (error *crcError) Error() string {
	return fmt.Sprintf("CRC validation failure.  expected=0x%X  actual=0x%X", error.expected, error.actual)
}
