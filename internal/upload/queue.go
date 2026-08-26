package upload

import "powergw/internal/model"

type Queue struct {
	items []model.Telemetry
	limit int
}

func NewQueue() *Queue {
	return &Queue{limit: 4096}
}

func (q *Queue) Push(tel model.Telemetry) error {
	if len(q.items) >= q.limit {
		return ErrQueueFull
	}
	q.items = append(q.items, tel)
	return nil
}

func (q *Queue) Pop() (model.Telemetry, bool) {
	if len(q.items) == 0 {
		return model.Telemetry{}, false
	}
	first := q.items[0]
	q.items = q.items[1:]
	return first, true
}

func (q *Queue) Peek() (model.Telemetry, bool) {
	if len(q.items) == 0 {
		return model.Telemetry{}, false
	}
	return q.items[0], true
}

func (q *Queue) Len() int {
	return len(q.items)
}

func (q *Queue) Clear() {
	q.items = nil
}
