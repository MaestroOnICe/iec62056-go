package iec62056

import (
	"errors"
	"fmt"

	"go.bug.st/serial"
)

type SerialClient struct {
	port   serial.Port
	config DeviceConfig
	proto  *Protocol
}

type SerialConfig struct {
	PortName string
	DeviceConfig
}

func NewSerialClient(config SerialConfig) (*SerialClient, error) {
	mode := &serial.Mode{
		BaudRate: int(config.InitialBaudRate),
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(config.PortName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	// Set initial timeouts
	port.SetReadTimeout(config.Timeout)

	client := &SerialClient{
		port:   port,
		config: config.DeviceConfig,
		proto:  NewProtocol(port, port, config.DeviceConfig),
	}

	return client, nil
}

func (c *SerialClient) Close() error {
	return c.port.Close()
}

func (c *SerialClient) ReadMeter() (*Message, error) {
	// Send initial request
	request, err := c.proto.RequestMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	_, err = c.port.Write(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read identification
	msg, err := c.proto.ReadIdentification()
	if err != nil {
		return nil, fmt.Errorf("failed to read identification: %w", err)
	}

	// Switch baud rate if needed
	if msg.BaudRate != c.config.InitialBaudRate {
		err = c.switchBaudRate(msg.BaudRate)
		if err != nil {
			return nil, fmt.Errorf("failed to switch baud rate: %w", err)
		}
	}

	// Read data blocks
	blocks, err := c.proto.ReadData()
	if err != nil {
		return nil, fmt.Errorf("failed to read data blocks: %w", err)
	}

	msg.DataBlocks = blocks
	return msg, nil
}

func (c *SerialClient) switchBaudRate(newRate BaudRate) error {
	err := c.proto.SwitchBaudRate(newRate)
	if err != nil {
		return err
	}

	// Update port baud rate
	return c.port.SetMode(&serial.Mode{
		BaudRate: int(newRate),
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	})
}

func (c *SerialClient) WriteValue(address string, value string) error {
	if c.config.ReadOnly {
		return errors.New("device is in read-only mode")
	}

	cmd := fmt.Sprintf("W%s(%s)\r\n", address, value)
	_, err := c.port.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send write command: %w", err)
	}

	// Read acknowledgment with timeout
	buf := make([]byte, 1)
	n, err := c.port.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read acknowledgment: %w", err)
	}

	if n == 0 {
		return ErrTimeout
	}

	if buf[0] != byte(0x06) { // ACK
		return errors.New("write not acknowledged")
	}

	return nil
}
