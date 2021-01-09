package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jdtobe/Zero/color"
	"github.com/jdtobe/Zero/env"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	envShowName        = "SHOW_NAME"
	envShowNameDefault = ""
)

var (
	ledNum = env.GetDefaultInt("LED_NUM", 50)
	ledLum = env.GetDefaultInt("LED_LUM", 255/2) // half-brightness
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

	lightShow(leds, env.GetDefault("SHOW_NAME", originalShowName))
}

func lightShow(leds *ws2811.WS2811, name string) {
	ctx := context.Background()

	pixels := make([]color.HSV, ledNum)
	go showHSV(ctx, leds, pixels, 24)

	aFn := animationByName(env.GetDefault(envShowName, envShowNameDefault))
	var arr error
	for arr == nil {
		arr = aFn(ctx, pixels)
	}
}

// showHSV displays a color.HSV slice to the LED String at the provided FPS.
func showHSV(ctx context.Context, m *ws2811.WS2811, pixels []color.HSV, fps int) {
	ticker := time.NewTicker(time.Second / time.Duration(fps))

	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting: Show")
			return
		case <-ticker.C:
			for i := 0; i < len(pixels); i++ {
				m.Leds(0)[i] = pixels[i].UInt32()
			}
			if err := m.Render(); err != nil {
				log.Fatal("Show Error: ", err)
			}
		}
	}
}
