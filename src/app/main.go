package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cast"
)

var (
	ledNum         = 50
	ledLum         = 255 / 2 // half-brightness
	drawDuration   = 200
	fadeDuration   = 10
	fadeMultiplier = 0.995

	colorRed   = hsv{H: 0, S: 1.0, V: 1.0}
	colorBlue  = hsv{H: 235, S: 1.0, V: 1.0}
	colorGreen = hsv{H: 135, S: 1.0, V: 1.0}
	colorWhite = hsv{H: 0, S: 0.0, V: 1.0}
)

func main() {
	fmt.Println("Pi0w Light Show")

	if ln := cast.ToInt(os.Getenv("LED_NUM")); ln != 0 {
		ledNum = ln
	}

	if ll := cast.ToInt(os.Getenv("LED_LUM")); ll != 0 {
		ledLum = ll
	}

	if dd := cast.ToInt(os.Getenv("DRAW_DURATION")); dd != 0 {
		drawDuration = dd
	}

	if fd := cast.ToInt(os.Getenv("FADE_DURATION")); fd != 0 {
		fadeDuration = fd
	}

	if fm := cast.ToFloat64(os.Getenv("FADE_MULTIPLIER")); fm != 0 {
		fadeMultiplier = fm
	}

	lightShow(ledNum, ledLum)
}

type hsv struct {
	H, S, V float64
}

type animationFunc func(context.Context, *ws2811.WS2811) error

func lightShow(ledNum, ledLum int) {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = ledLum
	opt.Channels[0].LedCount = ledNum
	opt.Channels[0].StripeType = ws2811.WS2811StripRGB

	m, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if err := m.Init(); err != nil {
		log.Fatal("Error: ", err)
	}
	defer m.Fini()

	var animationFn animationFunc
	animationFn = originalShow

	var arr error
	for arr == nil {
		arr = animationFn(context.Background(), m)
	}
}

func originalShow(ctx context.Context, m *ws2811.WS2811) error {
	ledChain := []hsv{colorGreen, colorRed, colorGreen, colorWhite}

	ledNum := len(m.Leds(0))
	leds := make([]hsv, ledNum)

	// Setup Faders
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Exiting: Fader")
				return
			default:
			}

			for i := range leds {
				leds[i].V = leds[i].V * fadeMultiplier
			}

			time.Sleep(time.Duration(fadeDuration) * time.Millisecond)
		}
	}()

	// Run Draw Loop
	go func() {
		offset := 0
		ledChainLen := len(ledChain)
		for {
			select {
			case <-ctx.Done():
				log.Println("Exiting; Draw Loop")
				return
			default:
			}

			// Set LED Colors:
			const split = 10
			for j := 0; j < ledNum/split; j++ {
				for k := j; k < ledNum; k += ledNum / split {
					leds[k] = ledChain[(k+offset)%len(ledChain)]
				}
				time.Sleep(time.Duration(drawDuration) * time.Millisecond)
			}
			offset = (offset + 1) % ledChainLen
		}
	}()

	// Show
	go func() {
		ticker := time.NewTicker(time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				log.Println("Exiting: Show")
				return
			case <-ticker.C:
				for i := 0; i < ledNum; i++ {
					r, g, b := colorful.Hsv(leds[i].H, leds[i].S, leds[i].V).RGB255()
					m.Leds(0)[i] = (uint32(r)<<8+uint32(g))<<8 + uint32(b)
				}
				if err := m.Render(); err != nil {
					log.Fatal("Show Error: ", err)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Animation Done!")
	}

	return nil
}
