package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	colorful "github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cast"
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

var opt = ws2811.DefaultOptions

func runServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Holiday Lights, are Go!\n"))
	})
	r.HandleFunc("/settings/{setting}", func(w http.ResponseWriter, r *http.Request) {
		v := r.FormValue("value")
		vars := mux.Vars(r)
		switch vars["setting"] {
		case "ledLum", "brightness":
			ledLum = cast.ToInt(v)
			opt.Channels[0].Brightness = ledLum
		case "drawDuration":
			drawDuration = cast.ToInt(v)
		case "fadeDuration":
			fadeDuration = cast.ToInt(v)
		case "fadeMultiplier":
			fadeMultiplier = cast.ToFloat64(v)
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown Setting: %q\n", vars["setting"])
		}
	})
	addr := ":80"
	fmt.Println("Pi0w Light Show")
	log.Println("Listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

const ledNum = 100

var (
	ledLum         = 200
	drawDuration   = 200
	fadeDuration   = 10
	fadeMultiplier = 0.995

	colorRed   = hsv{H: 0, S: 1.0, V: 1.0}
	colorBlue  = hsv{H: 235, S: 1.0, V: 1.0}
	colorGreen = hsv{H: 135, S: 1.0, V: 1.0}
	colorWhite = hsv{H: 0, S: 0.0, V: 1.0}
)

type hsv struct {
	H, S, V float64
}

func doLights() {
	opt.Channels[0].Brightness = ledLum
	opt.Channels[0].LedCount = ledNum
	// opt.Channels[0].StripeType = ws2811.WS2811StripBGR // red, green, red
	// opt.Channels[0].StripeType = ws2811.WS2811StripBRG // red, blue, red
	opt.Channels[0].StripeType = ws2811.WS2811StripRGB

	m, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	checkError(m.Init())
	defer m.Fini()

	ledChain := []hsv{colorGreen, colorRed, colorGreen, colorWhite}

	// muLeds := sync.RWMutex{}
	leds := [ledNum]hsv{}

	done := make(chan struct{})

	// Setup Faders
	go func() {
		for {
			select {
			case <-done:
				log.Printf("Exiting: Fader")
				return
			default:
			}

			// muLeds.Lock()
			for i := range leds {
				leds[i].V = leds[i].V * fadeMultiplier
			}
			// muLeds.Unlock()

			time.Sleep(time.Duration(fadeDuration) * time.Millisecond)
			// angle := float64(0.0)
			// speed := 1.0 - math.Sin(angle)
			// time.Sleep(time.Duration(int64(speed * 10 * float64(time.Millisecond))))
			// angle = angle + math.Pi/50
			// if twopi := (2 * math.Pi); angle > twopi {
			// 	angle = angle - twopi
			// }
		}
	}()

	// Run Draw Loop
	go func() {
		offset := 0
		ledChainLen := len(ledChain)
		for {
			select {
			case <-done:
				log.Println("Exiting; Draw Loop")
				return
			default:
			}

			// Set LED Colors:
			const split = 10
			for j := 0; j < ledNum/split; j++ {
				// muLeds.Lock()
				for k := j; k < ledNum; k += ledNum / split {
					leds[k] = ledChain[(k+offset)%len(ledChain)]
				}
				// muLeds.Unlock()
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
			case <-done:
				log.Println("Exiting: Show")
				return
			case <-ticker.C:
				// muLeds.RLock()
				for i := 0; i < ledNum; i++ {
					r, g, b := colorful.Hsv(leds[i].H, leds[i].S, leds[i].V).RGB255()
					m.Leds(0)[i] = (uint32(r)<<8+uint32(g))<<8 + uint32(b)
				}
				// muLeds.RUnlock()
				checkError(m.Render())
			}
		}
	}()

	select {
	case <-done:
		log.Println("System Exiting!")
	}

	// leds := []uint32{colorRed, colorWhite, colorGreen, colorWhite}
	// l := len(leds)
	// offset := 0
	// for {
	// 	for i := 0; i < ledNum; i++ {
	// 		m.Leds(0)[i] = leds[(i+offset)%l]
	// 	}

	// 	checkError(m.Render())
	// 	offset = (offset + 1) % l
	// 	time.Sleep(250 * time.Millisecond)
	// }
}
