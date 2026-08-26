package dedup

const DefaultWindow = 1024

type Window struct {
	size int
}

func NewWindow(size int) *Window {
	return &Window{size: size}
}

func (w *Window) Size() int {
	if w == nil || w.size <= 0 {
		return DefaultWindow
	}
	return w.size
}

func (w *Window) Half() int {
	size := w.Size()
	if size < 2 {
		return 1
	}
	return size / 2
}
