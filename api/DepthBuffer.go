package api

type DepthBuffer struct {
	Entries []float64
	w       int
}

func (d *DepthBuffer) Clear() {
	for i := range d.Entries {
		d.Entries[i] = 0
	}
}

func (d *DepthBuffer) At(x, y int) float64 {
	pos := y*d.w + x
	return d.Entries[pos]
}

func (d *DepthBuffer) Set(x, y int, w float64) {
	d.Entries[y*d.w+x] = w
}

func NewDepthBuffer(w, h int) *DepthBuffer {
	return &DepthBuffer{
		Entries: make([]float64, w*h),
		w:       w,
	}
}
