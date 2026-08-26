package parse

import "powergw/internal/model"

const (
	iec104Start     = 0x68
	iec104TypeME    = 0x09
	iec104TypeMS    = 0x0B
	iec104TypeMEB   = 0x0D
	iec104MinLength = 8
	iec104MaxLength = 253
)

func parseIEC104(raw []byte) (*model.Message, error) {
	if len(raw) < 6 {
		return nil, ErrBadFrame
	}
	if raw[0] != iec104Start {
		return nil, ErrBadFrame
	}
	length := int(raw[1])
	if length < iec104MinLength || length > iec104MaxLength {
		return nil, ErrBadFrame
	}
	if length != len(raw)-2 {
		return nil, ErrBadFrame
	}
	sendSeq := uint32(raw[2]>>1) | uint32(raw[3]>>1)<<7
	asdu := raw[6:]
	message, err := parseASDU(asdu)
	if err != nil {
		return nil, err
	}
	message.Seq = sendSeq
	return message, nil
}

func parseASDU(asdu []byte) (*model.Message, error) {
	if len(asdu) < 6 {
		return nil, ErrBadFrame
	}
	switch asdu[0] {
	case iec104TypeME:
		return parseMeasuredNormalized(asdu)
	case iec104TypeMS:
		return parseSinglePoint(asdu)
	case iec104TypeMEB:
		return parseMeasuredScaled(asdu)
	default:
		return nil, ErrBadFrame
	}
}

func parseMeasuredNormalized(asdu []byte) (*model.Message, error) {
	if len(asdu) < 10 {
		return nil, ErrBadFrame
	}
	count := int(asdu[1] & 0x7F)
	if count != 1 {
		return nil, ErrBadFrame
	}
	ioa := uint16(asdu[6]) | uint16(asdu[7])<<8
	rawValue := int16(uint16(asdu[9]) | uint16(asdu[10])<<8)
	return &model.Message{
		Kind: model.KindTelemetry,
		Addr: ioa,
		Fields: []model.Field{
			{Name: "ioa", Value: float64(ioa), Valid: true},
			{Name: "value", Value: float64(rawValue) / 100.0, Valid: true},
		},
	}, nil
}

func parseSinglePoint(asdu []byte) (*model.Message, error) {
	if len(asdu) < 9 {
		return nil, ErrBadFrame
	}
	count := int(asdu[1] & 0x7F)
	if count != 1 {
		return nil, ErrBadFrame
	}
	ioa := uint16(asdu[6]) | uint16(asdu[7])<<8
	state := float64(asdu[8] & 0x01)
	return &model.Message{
		Kind: model.KindStatus,
		Addr: ioa,
		Fields: []model.Field{
			{Name: "ioa", Value: float64(ioa), Valid: true},
			{Name: "state", Value: state, Valid: true},
		},
	}, nil
}

func parseMeasuredScaled(asdu []byte) (*model.Message, error) {
	if len(asdu) < 10 {
		return nil, ErrBadFrame
	}
	count := int(asdu[1] & 0x7F)
	if count != 1 {
		return nil, ErrBadFrame
	}
	ioa := uint16(asdu[6]) | uint16(asdu[7])<<8
	rawValue := int16(uint16(asdu[9]) | uint16(asdu[10])<<8)
	return &model.Message{
		Kind: model.KindTelemetry,
		Addr: ioa,
		Fields: []model.Field{
			{Name: "ioa", Value: float64(ioa), Valid: true},
			{Name: "value", Value: float64(rawValue), Valid: true},
		},
	}, nil
}

func BuildIEC104Frame(addr uint16, value int16, seq uint16) []byte {
	asdu := make([]byte, 0, 12)
	asdu = append(asdu, iec104TypeME, 0x01, 0x03, 0x00, 0x01, 0x00)
	asdu = append(asdu, byte(addr), byte(addr>>8), 0x00)
	asdu = append(asdu, byte(value), byte(value>>8), 0x00)
	frame := make([]byte, 0, 6+len(asdu))
	frame = append(frame, iec104Start, byte(len(asdu)+4), byte((seq&0x7F)<<1), byte(seq>>7), 0x00, 0x00)
	frame = append(frame, asdu...)
	return frame
}
