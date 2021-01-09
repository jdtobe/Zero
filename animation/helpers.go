package animation

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jdtobe/Zero/color"
)

// FadeHSV fades the Value of an HSV slice of pixels by m every d ms.
func ahFadeHSV(ctx context.Context, pixels []color.HSV, d int, m float64) {
	ticker := time.NewTicker(time.Duration(d) * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Exiting: Fader")
			return
		case <-ticker.C:
			for i := range pixels {
				pixels[i].V = pixels[i].V * m
			}
		}
	}
}

// SparkleHSV "sparkles" an HSV slice of pixels, relative to o in Hue, Saturation and Value, by the provided limits (hl, sl, vl), every d ms.
func ahSparkleHSV(ctx context.Context, pixels []color.HSV, o color.HSV, hl, sl, vl float64, d int) {
	ticker := time.NewTicker(time.Duration(d) * time.Millisecond)

	var h, s, v float64
	for {
		select {
		case <-ctx.Done():
			log.Printf("Exiting: Sparkle")
			return
		case <-ticker.C:
			for i := range pixels {
				h = o.H + (rand.Float64() * hl) - (hl / 2)
				s = o.S + (rand.Float64() * sl) - (sl / 2)
				v = o.V + (rand.Float64() * vl) - (vl / 2)

				pixels[i].H = h
				pixels[i].S = s
				pixels[i].V = v
			}
		}
	}
}

// SparkleSkipHSV "sparkles" an HSV slice of pixels,
//
// skip - A map of pixels to skip, based on the map index.
// o - Original pixel value to sparkle.
// hl, sl, vl - Hue, Saturation and Value limits to sparkle.
// d - Sparkle delay, in ms.
func ahSparkleSkipHSV(ctx context.Context, pixels []color.HSV, skip map[int]bool, o color.HSV, hl, sl, vl float64, d int) {
	ticker := time.NewTicker(time.Duration(d) * time.Millisecond)

	var h, s, v float64
	for {
		select {
		case <-ctx.Done():
			log.Printf("Exiting: Sparkle")
			return
		case <-ticker.C:
			for i := range pixels {
				if skip[i] {
					continue
				}

				h = o.H + (rand.Float64() * hl) - (hl / 2)
				s = o.S + (rand.Float64() * sl) - (sl / 2)
				v = o.V + (rand.Float64() * vl) - (vl / 2)

				pixels[i].H = h
				pixels[i].S = s
				pixels[i].V = v
			}
		}
	}
}
