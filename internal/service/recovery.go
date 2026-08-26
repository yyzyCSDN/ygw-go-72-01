package service

import "powergw/internal/model"

func Restore(g *Gateway, persisted *model.TableSnapshot) error {
	if persisted != nil {
		g.Recovery.Store(persisted)
	}
	ResetParserState(g)
	return g.Recover()
}

func SnapshotForRecovery(g *Gateway) *model.TableSnapshot {
	return g.Mapper.Snapshot()
}

func ResetParserState(g *Gateway) {
	g.Parser.ResetAll()
}
