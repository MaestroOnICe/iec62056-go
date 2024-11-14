package iec62056

import (
	"fmt"
	"time"
)

func (p *Protocol) SwitchBaudRate(newRate BaudRate) error {
	rateChar, err := getBaudRateChar(newRate)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(p.writer, "%c\r\n", rateChar)
	if err != nil {
		return err
	}

	time.Sleep(300 * time.Millisecond)
	return nil
}

func getBaudRateChar(rate BaudRate) (byte, error) {
	switch rate {
	case Baud300:
		return '0', nil
	case Baud600:
		return '1', nil
	case Baud1200:
		return '2', nil
	case Baud2400:
		return '3', nil
	case Baud4800:
		return '4', nil
	case Baud9600:
		return '5', nil
	case Baud19200:
		return '6', nil
	default:
		return 0, ErrInvalidBaudRate
	}
}
