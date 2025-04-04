package api

import "fmt"

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
	if pos >= len(d.Entries) {
		panic(fmt.Sprintf("out of bounds: %v / %v -> %v %v", x, y, pos, d.w))
	}
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
