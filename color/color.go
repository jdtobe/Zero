package color

import (
	colorful "github.com/lucasb-eyer/go-colorful"
)

// Value can be read in multiple formats.
type Value interface {
	UInt32() uint32
}

// Common Colors
var (
	Red    = HSV{H: 0, S: 1.0, V: 1.0}
	Blue   = HSV{H: 215, S: 1.0, V: 1.0}
	Green  = HSV{H: 125, S: 1.0, V: 0.75}
	Purple = HSV{H: 285, S: 1.0, V: 1.0}
	Pink   = HSV{H: 305, S: 1.0, V: 1.0}
	White  = HSV{H: 0, S: 0.0, V: 1.0}
)

// HSV is a Color defined bye its Hue, Saturation, and Value values.
type HSV struct {
	H, S, V float64
}

// UInt32 returns this HSV Value as a uint32
func (c *HSV) UInt32() uint32 {
	r, g, b := colorful.Hsv(c.H, c.S, c.V).RGB255()
	return (uint32(r)<<8+uint32(g))<<8 + uint32(b)
}
