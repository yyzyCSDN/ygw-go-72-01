package mapper

import "powergw/internal/model"

func PendingAddrs(snap *model.TableSnapshot) []uint16 {
	if snap == nil {
		return nil
	}
	addrs := make([]uint16, 0, snap.Len())
	for _, point := range snap.Points {
		if !point.Uploaded {
			addrs = append(addrs, point.RawAddr)
		}
	}
	return addrs
}

func LatestUploaded(snap *model.TableSnapshot) []uint16 {
	if snap == nil {
		return nil
	}
	addrs := make([]uint16, 0, snap.Len())
	for _, point := range snap.Points {
		if point.Uploaded {
			addrs = append(addrs, point.RawAddr)
		}
	}
	return addrs
}

func ApplyUploaded(m *Mapper, addrs []uint16, version string) error {
	for _, addr := range addrs {
		if err := m.MarkUploaded(addr, version); err != nil {
			return err
		}
	}
	return nil
}

func SnapshotTable(m *Mapper) *model.TableSnapshot {
	return m.Snapshot()
}
