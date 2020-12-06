package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func checkError(err error) {
	if err == nil {
		return
	}

	debug.PrintStack()
	panic(err)
}

func main() {
	// runServer()
	go runServer()
	doLights()
}

func runServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Holiday Lights, are Go!\n"))
	})
	addr := ":80"
	fmt.Println("Piow Light Show")
	log.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

const (
	// ledLum = 128
	ledLum = 200
	ledNum = 50

	// colorRed = uint32(0x00FF00)
	colorRed   = uint32(0xFF0000)
	colorBlue  = uint32(0x00FF00)
	colorGreen = uint32(0x0000FF)
	colorWhite = uint32(0xFFFFFF)
)

func doLights() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = ledLum
	opt.Channels[0].LedCount = ledNum
	opt.Channels[0].StripeType = ws2811.WS2811StripBRG

	m, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	checkError(m.Init())
	defer m.Fini()

	leds := []uint32{colorRed, colorWhite, colorGreen, colorWhite}
	l := len(leds)
	offset := 0
	for {
		for i := 0; i < ledNum; i++ {
			m.Leds(0)[i] = leds[(i+offset)%l]
		}

		checkError(m.Render())
		offset = (offset + 1) % l
		time.Sleep(250 * time.Millisecond)
	}
}
