package parse

import "powergw/internal/model"

const (
	modbusReadHolding  = 0x03
	modbusReadInput    = 0x04
	modbusMaxRegisters = 16
)

func parseModbus(raw []byte) (*model.Message, error) {
	if len(raw) < 5 {
		return nil, ErrBadFrame
	}
	fn := raw[1]
	switch fn {
	case modbusReadHolding, modbusReadInput:
		return parseModbusRead(raw)
	default:
		return nil, ErrBadFrame
	}
}

func parseModbusRead(raw []byte) (*model.Message, error) {
	byteCount := int(raw[2])
	if byteCount == 0 || byteCount%2 != 0 || byteCount > modbusMaxRegisters*2 {
		return nil, ErrBadFrame
	}
	if byteCount != len(raw)-5 {
		return nil, ErrBadFrame
	}
	registers := byteCount / 2
	fields := make([]model.Field, 0, registers)
	for index := 0; index < registers; index++ {
		offset := 3 + index*2
		reg := uint16(raw[offset])<<8 | uint16(raw[offset+1])
		fields = append(fields, model.Field{
			Name:  "reg",
			Value: float64(reg),
			Valid: true,
		})
	}
	return &model.Message{
		Kind:   model.KindTelemetry,
		Addr:   uint16(registers - 1),
		Fields: fields,
	}, nil
}

func BuildModbusFrame(fn byte, values ...uint16) []byte {
	frame := make([]byte, 0, 5+len(values)*2)
	frame = append(frame, 0x01, fn, byte(len(values)*2))
	for _, value := range values {
		frame = append(frame, byte(value>>8), byte(value))
	}
	crc := uint16(0xFFFF)
	for _, item := range frame {
		crc ^= uint16(item)
		for bit := 0; bit < 8; bit++ {
			if crc&0x0001 != 0 {
				crc = crc>>1 ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	frame = append(frame, byte(crc), byte(crc>>8))
	return frame
}
