package service

import "powergw/internal/model"

func ProcessChannel(g *Gateway, channelID string, frames [][]byte) (int, error) {
	count := 0
	for _, raw := range frames {
		if err := g.Ingest(channelID, raw); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func MapMessage(g *Gateway, message *model.Message) (model.Point, bool) {
	if message == nil {
		return model.Point{}, false
	}
	return g.Mapper.Resolve(message.Addr)
}

func BuildTelemetry(g *Gateway, message *model.Message, point model.Point) model.Telemetry {
	tel := model.NewTelemetry(point.Name, point.RawAddr, messageValue(message, point), message.Seq, g.Clock())
	tel.VersionID = message.VersionID
	return tel
}

func SubmitMessage(g *Gateway, message *model.Message) error {
	point, ok := MapMessage(g, message)
	if !ok {
		return ErrUnmappedPoint
	}
	return g.Upload.Submit(BuildTelemetry(g, message, point))
}

func messageValue(message *model.Message, point model.Point) float64 {
	if message == nil {
		return 0
	}
	if value, ok := message.Field("value"); ok && value.Valid {
		return value.Value
	}
	if len(message.Fields) > 0 && message.Fields[0].Valid {
		return message.Fields[0].Value*point.Scale + point.Offset
	}
	return 0
}
