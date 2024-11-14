package iec62056

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Protocol struct {
	config DeviceConfig
	reader io.Reader
	writer io.Writer
}

func NewProtocol(reader io.Reader, writer io.Writer, config DeviceConfig) *Protocol {
	return &Protocol{
		config: config,
		reader: reader,
		writer: writer,
	}
}

func (p *Protocol) RequestMessage() ([]byte, error) {
	switch p.config.Mode {
	case ModeA:
		return []byte("/?!\r\n"), nil
	case ModeB:
		return []byte("/?\r\n"), nil
	case ModeC:
		return []byte("/?"), nil
	case ModeD:
		return []byte(fmt.Sprintf("/?%s!\r\n", p.config.Password)), nil
	default:
		return nil, ErrInvalidMode
	}
}

func (p *Protocol) ReadIdentification() (*Message, error) {
	reader := bufio.NewReader(p.reader)
	msg := &Message{}

	idLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(idLine) < 7 {
		return nil, ErrInvalidFormat
	}

	msg.Mode = ProtocolMode(idLine[4])
	msg.Manufacturer = idLine[1:4]

	if msg.Mode != ModeA {
		baudChar := idLine[5]
		msg.BaudRate, err = parseBaudRate(baudChar)
		if err != nil {
			return nil, err
		}
	} else {
		msg.BaudRate = Baud300
	}

	return msg, nil
}

func (p *Protocol) ReadData() ([]DataBlock, error) {
	var blocks []DataBlock
	reader := bufio.NewReader(p.reader)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "!" {
			break
		}

		block, err := parseDataBlock(line)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, *block)
	}

	return blocks, nil
}
