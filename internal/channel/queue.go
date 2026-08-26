package channel

import "powergw/internal/model"

type Queue struct {
	items []*model.Message
	limit int
}

func NewMessageQueue() *Queue {
	return &Queue{limit: 2048}
}

func (q *Queue) Push(message *model.Message) error {
	if len(q.items) >= q.limit {
		return ErrQueueFull
	}
	q.items = append(q.items, message)
	return nil
}

func (q *Queue) Pop() (*model.Message, bool) {
	if len(q.items) == 0 {
		return nil, false
	}
	first := q.items[0]
	q.items = q.items[1:]
	return first, true
}

func (q *Queue) Len() int {
	return len(q.items)
}

func (q *Queue) Clear() {
	q.items = nil
}
