package main

import (
	"context"
	"fmt"
	"log"

	colorful "github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

var (
	ledNum         = getEnvDefaultInt("LED_NUM", 50)
	ledLum         = getEnvDefaultInt("LED_LUM", 255/2) // half-brightness
	drawDuration   = getEnvDefaultInt("DRAW_DURATION", 200)
	fadeDuration   = getEnvDefaultInt("FADE_DURATION", 5)
	fadeMultiplier = getEnvDefaultFloat64("FADE_MULTIPLIER", 0.995)

	colorRed    = hsv{H: 0, S: 1.0, V: 1.0}
	colorBlue   = hsv{H: 215, S: 1.0, V: 1.0}
	colorGreen  = hsv{H: 125, S: 1.0, V: 0.75}
	colorPurple = hsv{H: 285, S: 1.0, V: 1.0}
	colorPink   = hsv{H: 305, S: 1.0, V: 1.0}
	colorWhite  = hsv{H: 0, S: 0.0, V: 1.0}
)

func main() {
	fmt.Println("Pi0w Light Show")

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = ledLum
	opt.Channels[0].LedCount = ledNum
	opt.Channels[0].StripeType = ws2811.WS2811StripRGB

	leds, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if err := leds.Init(); err != nil {
		log.Fatal("Error: ", err)
	}
	defer leds.Fini()

	lightShow(leds, getEnvDefault("SHOW_NAME", originalShowName))
}

type hsv struct {
	H, S, V float64
}

func (c *hsv) RgbUInt32() uint32 {
	r, g, b := colorful.Hsv(c.H, c.S, c.V).RGB255()
	return (uint32(r)<<8+uint32(g))<<8 + uint32(b)
}

type color interface {
	RgbUInt32() uint32
}

type animationFunc func(context.Context, []hsv) error

func lightShow(leds *ws2811.WS2811, name string) {
	var animationFn animationFunc
	switch name {
	case xmasTreeShowName:
		animationFn = xmasTreeShow
	case originalShowName:
		fallthrough
	default:
		animationFn = originalShow
	}

	ctx := context.Background()

	// Show
	pixels := make([]hsv, ledNum)
	go ahShowHSV(ctx, leds, pixels)

	var arr error
	for arr == nil {
		arr = animationFn(ctx, pixels)
	}
}
