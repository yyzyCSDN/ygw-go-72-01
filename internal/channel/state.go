package channel

import "powergw/internal/model"

func CanConnect(record model.ChannelRecord) bool {
	return record.State == model.ChannelIdle
}

func CanSync(record model.ChannelRecord) bool {
	return record.State == model.ChannelConnected
}

func CanRun(record model.ChannelRecord) bool {
	return record.State == model.ChannelSyncing
}

func CanFault(record model.ChannelRecord) bool {
	return record.State == model.ChannelConnected ||
		record.State == model.ChannelSyncing ||
		record.State == model.ChannelRunning
}

func NextState(record model.ChannelRecord, next model.ChannelState) bool {
	switch next {
	case model.ChannelConnected:
		return CanConnect(record)
	case model.ChannelSyncing:
		return CanSync(record)
	case model.ChannelRunning:
		return CanRun(record)
	case model.ChannelFault:
		return CanFault(record)
	default:
		return false
	}
}

func ApplyTransition(record model.ChannelRecord, next model.ChannelState) (model.ChannelRecord, error) {
	if !NextState(record, next) {
		return record, ErrInvalidState
	}
	record.State = next
	return record, nil
}
