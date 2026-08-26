package upload

import (
	"testing"

	"powergw/internal/model"
)

func TestQueuePushPop(t *testing.T) {
	queue := NewQueue()
	tel := model.NewTelemetry("电压", 1, 22.5, 1, 100)
	if err := queue.Push(tel); err != nil {
		t.Fatal(err)
	}
	if queue.Len() != 1 {
		t.Fatalf("len = %d", queue.Len())
	}
	got, ok := queue.Pop()
	if !ok || got.PointName != "电压" {
		t.Fatalf("pop failed: %v %v", got, ok)
	}
	if queue.Len() != 0 {
		t.Fatalf("len after pop = %d", queue.Len())
	}
	if _, ok := queue.Pop(); ok {
		t.Fatal("pop on empty succeeded")
	}
}

func TestQueuePeek(t *testing.T) {
	queue := NewQueue()
	_ = queue.Push(model.NewTelemetry("电压", 1, 22.5, 1, 100))
	peeked, ok := queue.Peek()
	if !ok || peeked.RawAddr != 1 {
		t.Fatalf("peek failed: %v %v", peeked, ok)
	}
	if queue.Len() != 1 {
		t.Fatalf("peek mutated queue: %d", queue.Len())
	}
}

func TestQueueClear(t *testing.T) {
	queue := NewQueue()
	_ = queue.Push(model.NewTelemetry("电压", 1, 22.5, 1, 100))
	queue.Clear()
	if queue.Len() != 0 {
		t.Fatalf("len after clear = %d", queue.Len())
	}
}

func TestRecoveryStore(t *testing.T) {
	recovery := NewRecovery()
	if recovery.Has() {
		t.Fatal("empty recovery has snapshot")
	}
	table := model.NewPointTable("t1", "v1")
	table.Add(model.Point{RawAddr: 1, Name: "电压"})
	snap := model.NewTableSnapshot(table, 1)
	recovery.Store(snap)
	if !recovery.Has() {
		t.Fatal("stored recovery missing")
	}
	if recovery.Latest() != snap {
		t.Fatal("latest mismatch")
	}
	recovery.Clear()
	if recovery.Has() {
		t.Fatal("cleared recovery still has snapshot")
	}
}
