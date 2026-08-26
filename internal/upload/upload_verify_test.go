package upload

import (
	"testing"

	"powergw/internal/map"
	"powergw/internal/model"
)

type okWriter struct{}

func (w *okWriter) Send(tel model.Telemetry) error {
	return nil
}

func TestUploadRecoveryUsesLatestSnapshot(t *testing.T) {
	m := mapper.New()
	table := mapper.BuildTable("t1", "v1", []model.Point{
		{RawAddr: 1, Name: "电压"},
		{RawAddr: 2, Name: "电流"},
		{RawAddr: 3, Name: "功率"},
	})
	if err := m.Load(table); err != nil {
		t.Fatal(err)
	}
	if err := m.ApplyTable("t1"); err != nil {
		t.Fatal(err)
	}
	writer := &okWriter{}
	uploader := NewUploader(m, writer)
	if err := uploader.Submit(model.NewTelemetry("电压", 1, 22.5, 1, 100)); err != nil {
		t.Fatal(err)
	}
	if err := uploader.Submit(model.NewTelemetry("电流", 2, 10.0, 2, 100)); err != nil {
		t.Fatal(err)
	}
	if _, err := uploader.Flush(m.Snapshot()); err != nil {
		t.Fatal(err)
	}
	if err := uploader.Recover(m.Snapshot()); err != nil {
		t.Fatal(err)
	}
	if uploader.PendingCount() != 1 {
		t.Fatalf("pending = %d, want 1", uploader.PendingCount())
	}
	if uploader.DoneCount() != 2 {
		t.Fatalf("done = %d, want 2", uploader.DoneCount())
	}
}
