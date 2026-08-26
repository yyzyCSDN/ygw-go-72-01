package service

import (
	"powergw/internal/model"
	"powergw/internal/parse"
)

func FeedDemo(g *Gateway, rounds int) (int, error) {
	count := 0
	for i := 0; i < rounds; i++ {
		proto, err := g.Controller.ProtocolFor("sta-a")
		if err != nil {
			return count, err
		}
		if proto == model.ProtoModbus {
			if err := g.Ingest("sta-a", parse.BuildModbusFrame(0x03, uint16(220+i))); err != nil {
				return count, err
			}
			count++
			if err := g.Ingest("sta-b", parse.BuildModbusFrame(0x04, uint16(50+i))); err != nil {
				return count, err
			}
			count++
			continue
		}
		if err := g.Ingest("sta-a", parse.BuildIEC104Frame(4001, int16(2200+i), uint16(1+i))); err != nil {
			return count, err
		}
		count++
		if err := g.Ingest("sta-b", parse.BuildIEC104Frame(4002, int16(100+i), uint16(1+i))); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
