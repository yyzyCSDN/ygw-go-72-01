package model

import "testing"

func TestProtocolString(t *testing.T) {
	cases := map[Protocol]string{
		ProtoIEC104: "iec104",
		ProtoModbus: "modbus",
		Protocol(9): "unknown",
	}
	for proto, want := range cases {
		if got := proto.String(); got != want {
			t.Fatalf("proto %d string = %q, want %q", proto, got, want)
		}
	}
}

func TestProtocolFromString(t *testing.T) {
	if proto, ok := ProtocolFromString("modbus"); !ok || proto != ProtoModbus {
		t.Fatalf("modbus lookup failed: %v %v", proto, ok)
	}
	if _, ok := ProtocolFromString("tcp"); ok {
		t.Fatal("unexpected protocol accepted")
	}
}

func TestMessageField(t *testing.T) {
	message := &Message{
		Fields: []Field{
			{Name: "ioa", Value: 4001, Valid: true},
			{Name: "value", Value: 22.5, Valid: true},
		},
	}
	field, ok := message.Field("value")
	if !ok || field.Value != 22.5 {
		t.Fatalf("field lookup failed: %v %v", field, ok)
	}
	if _, ok := message.Field("missing"); ok {
		t.Fatal("missing field reported present")
	}
	if got := message.FieldValue("ioa"); got != 4001 {
		t.Fatalf("field value = %v", got)
	}
	if got := message.FieldValue("missing"); got != 0 {
		t.Fatalf("missing field value = %v", got)
	}
	names := message.FieldNames()
	if len(names) != 2 || names[0] != "ioa" {
		t.Fatalf("field names = %v", names)
	}
	var nilMessage *Message
	if _, ok := nilMessage.Field("value"); ok {
		t.Fatal("nil message field lookup succeeded")
	}
}

func TestMessageKindString(t *testing.T) {
	if KindTelemetry.String() != "telemetry" || KindStatus.String() != "status" || KindEvent.String() != "event" {
		t.Fatalf("kind strings wrong: %s %s %s", KindTelemetry, KindStatus, KindEvent)
	}
}
