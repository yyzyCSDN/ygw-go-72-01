package model

type Protocol int

const (
	ProtoIEC104 Protocol = iota + 1
	ProtoModbus
)

func (p Protocol) String() string {
	switch p {
	case ProtoIEC104:
		return "iec104"
	case ProtoModbus:
		return "modbus"
	default:
		return "unknown"
	}
}

func ProtocolFromString(value string) (Protocol, bool) {
	switch value {
	case "iec104":
		return ProtoIEC104, true
	case "modbus":
		return ProtoModbus, true
	default:
		return 0, false
	}
}

type MessageKind int

const (
	KindTelemetry MessageKind = iota + 1
	KindStatus
	KindEvent
)

func (k MessageKind) String() string {
	switch k {
	case KindTelemetry:
		return "telemetry"
	case KindStatus:
		return "status"
	case KindEvent:
		return "event"
	default:
		return "unknown"
	}
}

type Field struct {
	Name  string
	Value float64
	Valid bool
}

type Message struct {
	ChannelID string
	Seq       uint32
	Proto     Protocol
	Addr      uint16
	Kind      MessageKind
	Fields    []Field
	VersionID string
}

func (m *Message) Field(name string) (Field, bool) {
	if m == nil {
		return Field{}, false
	}
	for _, field := range m.Fields {
		if field.Name == name {
			return field, true
		}
	}
	return Field{}, false
}

func (m *Message) FieldValue(name string) float64 {
	field, ok := m.Field(name)
	if !ok {
		return 0
	}
	return field.Value
}

func (m *Message) FieldNames() []string {
	if m == nil {
		return nil
	}
	names := make([]string, 0, len(m.Fields))
	for _, field := range m.Fields {
		names = append(names, field.Name)
	}
	return names
}
