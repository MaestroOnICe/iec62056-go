package iec62056

import "time"

type ProtocolMode byte

const (
	ModeA ProtocolMode = 'A'
	ModeB ProtocolMode = 'B'
	ModeC ProtocolMode = 'C'
	ModeD ProtocolMode = 'D'
)

type BaudRate int

const (
	Baud300   BaudRate = 300
	Baud600   BaudRate = 600
	Baud1200  BaudRate = 1200
	Baud2400  BaudRate = 2400
	Baud4800  BaudRate = 4800
	Baud9600  BaudRate = 9600
	Baud19200 BaudRate = 19200
)

type Message struct {
	Identification string
	Manufacturer   string
	BaudRate       BaudRate
	Mode           ProtocolMode
	DataBlocks     []DataBlock
	Password       string
}

type DataBlock struct {
	Address     string
	Value       string
	Unit        string
	ValuesCount int
	ReadWrite   bool
}

type DeviceConfig struct {
	InitialBaudRate BaudRate
	Mode            ProtocolMode
	Password        string
	Timeout         time.Duration
	ReadOnly        bool
}
