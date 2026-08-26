package service

import (
	"sort"
)

type ChannelStatus struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Session   string `json:"session"`
	SyncSeq   uint64 `json:"syncSeq"`
	Frames    uint64 `json:"frames"`
	Errors    uint64 `json:"errors"`
	Processed int    `json:"processed"`
}

type VersionStatus struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Proto string `json:"proto"`
	Table string `json:"table"`
}

type Status struct {
	Channels     []ChannelStatus `json:"channels"`
	Versions     []VersionStatus `json:"versions"`
	ActiveTable  string          `json:"activeTable"`
	Pending      int             `json:"pending"`
	Uploaded     int             `json:"uploaded"`
	UploadSeq    uint64          `json:"uploadSeq"`
	DedupSeen    int             `json:"dedupSeen"`
	SinkCount    int             `json:"sinkCount"`
	OpenSessions int             `json:"openSessions"`
}

func (g *Gateway) Status() *Status {
	status := &Status{
		ActiveTable:  g.Mapper.ActiveID(),
		Pending:      g.Upload.PendingCount(),
		Uploaded:     g.Mapper.UploadedCount(),
		UploadSeq:    g.Upload.StateSeq(),
		DedupSeen:    g.Dedup.SeenCount(),
		SinkCount:    g.Sink.Count(),
		OpenSessions: g.Sessions.OpenCount(),
	}
	records := g.Channels.Records()
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	for _, record := range records {
		session := "closed"
		if g.Sessions.IsOpen(record.ID) {
			session = "open"
		}
		status.Channels = append(status.Channels, ChannelStatus{
			ID:        record.ID,
			State:     record.State.String(),
			Session:   session,
			SyncSeq:   record.SyncSeq,
			Frames:    record.Frames,
			Errors:    record.Errors,
			Processed: g.Channels.ProcessedCount(record.ID),
		})
	}
	versions := g.Controller.Snapshot().History
	for _, item := range versions {
		status.Versions = append(status.Versions, VersionStatus{
			ID:    item.ID,
			State: item.State.String(),
			Proto: item.Proto.String(),
			Table: item.TableID,
		})
	}
	return status
}
