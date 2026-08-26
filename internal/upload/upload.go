package upload

import (
	"powergw/internal/map"
	"powergw/internal/model"
)

type Writer interface {
	Send(tel model.Telemetry) error
}

type Uploader struct {
	mapper *mapper.Mapper
	writer Writer
	queue  *Queue
	state  *model.UploadState
	seq    uint64
}

func NewUploader(m *mapper.Mapper, w Writer) *Uploader {
	return &Uploader{
		mapper: m,
		writer: w,
		queue:  NewQueue(),
		state:  model.NewUploadState("", 0),
	}
}

func (u *Uploader) Submit(tel model.Telemetry) error {
	return u.queue.Push(tel)
}

func (u *Uploader) Flush(snap *model.TableSnapshot) (int, error) {
	if snap == nil {
		return 0, ErrNoSnapshot
	}
	u.seq++
	u.state = model.NewUploadState(snap.TableID, u.seq)
	count := 0
	for {
		tel, ok := u.queue.Pop()
		if !ok {
			break
		}
		if err := u.writer.Send(tel); err != nil {
			return count, err
		}
		if err := u.mapper.MarkUploaded(tel.RawAddr, snap.Version); err != nil {
			return count, err
		}
		u.state.Done[tel.RawAddr] = true
		count++
	}
	return count, nil
}

func (u *Uploader) Writeback(snap *model.TableSnapshot) error {
	if snap == nil {
		return ErrNoSnapshot
	}
	for _, point := range snap.Points {
		if point.Uploaded {
			if err := u.mapper.MarkUploaded(point.RawAddr, snap.Version); err != nil {
				return err
			}
		}
	}
	u.state = u.state.Merge(model.NewUploadState(snap.TableID, u.seq))
	return nil
}

func (u *Uploader) Recover(snap *model.TableSnapshot) error {
	if snap == nil {
		return ErrNoSnapshot
	}
	u.state = model.NewUploadState(snap.TableID, u.seq+1)
	for _, point := range snap.Points {
		if point.Uploaded {
			u.state.Done[point.RawAddr] = true
		} else {
			u.state.Pending[point.RawAddr] = true
		}
	}
	return nil
}

func (u *Uploader) PendingCount() int {
	return u.state.PendingCount()
}

func (u *Uploader) DoneCount() int {
	return u.state.DoneCount()
}

func (u *Uploader) QueueLen() int {
	return u.queue.Len()
}

func (u *Uploader) SnapshotID() string {
	return u.state.SnapshotID
}

func (u *Uploader) StateSeq() uint64 {
	return u.state.Seq
}
