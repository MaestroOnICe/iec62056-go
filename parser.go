package iec62056

import "strings"

func parseDataBlock(line string) (*DataBlock, error) {
	readWrite := strings.HasPrefix(line, "W")
	if readWrite {
		line = line[1:]
	}

	parts := strings.Split(line, "(")
	if len(parts) != 2 {
		return nil, ErrInvalidDataBlock
	}

	address := parts[0]
	valueStr := strings.TrimSuffix(parts[1], ")")
	if valueStr == parts[1] {
		return nil, ErrInvalidDataBlock
	}

	valueParts := strings.Split(valueStr, "*")
	block := &DataBlock{
		Address:   address,
		Value:     valueParts[0],
		ReadWrite: readWrite,
	}

	if len(valueParts) > 1 {
		block.Unit = valueParts[1]
	}

	if strings.Contains(block.Value, ",") {
		block.ValuesCount = len(strings.Split(block.Value, ","))
	} else {
		block.ValuesCount = 1
	}

	return block, nil
}

func parseBaudRate(char byte) (BaudRate, error) {
	switch char {
	case '0':
		return Baud300, nil
	case '1':
		return Baud600, nil
	case '2':
		return Baud1200, nil
	case '3':
		return Baud2400, nil
	case '4':
		return Baud4800, nil
	case '5':
		return Baud9600, nil
	case '6':
		return Baud19200, nil
	default:
		return 0, ErrInvalidBaudRate
	}
}
