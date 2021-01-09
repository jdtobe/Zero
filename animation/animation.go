package animation

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jdtobe/Zero/color"
	"github.com/jdtobe/Zero/env"
)

//////////////////// DELETE THESE ///////////////////////////
var (
	ledNum         = env.GetDefaultInt("LED_NUM", 50)
	ledLum         = env.GetDefaultInt("LED_LUM", 255/2) // half-brightness
	drawDuration   = env.GetDefaultInt("DRAW_DURATION", 200)
	fadeDuration   = env.GetDefaultInt("FADE_DURATION", 5)
	fadeMultiplier = env.GetDefaultFloat64("FADE_MULTIPLIER", 0.995)
)

//////////////////// DELETE THESE ///////////////////////////

// XmasTree is an animation that resembles a glittery green tree with occasional bulbs that show up and fade out.
func XmasTree(ctx context.Context, pixels []color.HSV) error {
	fmt.Println("Starting X-Mas Tree Show...")

	skip := map[int]bool{}

	// Draw "Bulbs"
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break
			case <-ticker.C:
				// Pick a random bulb.
				i := rand.Intn(ledNum)
				if skip[i] {
					continue
				}
				go func(ctx context.Context, c color.HSV) {
					skip[i] = true
					pixels[i] = c
					fadingOut := true

					// Bulb Fade Routine
					ticker := time.NewTicker(time.Duration(fadeDuration) * time.Millisecond)

					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							var v float64
							if fadingOut {
								v = pixels[i].V * fadeMultiplier
								if v < 0.1 {
									fadingOut = false
									pixels[i].H = color.Green.H
									pixels[i].S = color.Green.S
								}
								pixels[i].V = v
								continue
							}

							v = pixels[i].V * (2 - fadeMultiplier)
							if v > color.Green.V {
								pixels[i].V = color.Green.V
								skip[i] = false
								ticker.Stop()
								return
							}
							pixels[i].V = v
						}
					}
				}(ctx, []color.HSV{color.Red, color.Blue, color.Purple, color.Pink, color.White}[rand.Intn(5)])
			}
		}
	}()

	// Sparkle
	o := color.Green
	o.S -= 0.25
	go ahSparkleSkipHSV(ctx, pixels, skip, o, 80, 0.5, 0.25, 5000)

	select {
	case <-ctx.Done():
		log.Println("Animation Done!")
	}

	return nil
}

// Original is the first animation written.
func Original(ctx context.Context, pixels []color.HSV) error {
	fmt.Println("Starting Original Show...")

	ledChain := []color.HSV{color.Green, color.Red, color.Green, color.White}

	// Setup Faders
	go ahFadeHSV(ctx, pixels, fadeDuration, fadeMultiplier)

	// Draw Loop
	go func() {
		offset := 0
		ledChainLen := len(ledChain)
		for {
			select {
			case <-ctx.Done():
				log.Println("Exiting: Draw Loop")
				return
			default:
			}

			// Set LED Colors:
			const split = 10
			for j := 0; j < ledNum/split; j++ {
				for k := j; k < ledNum; k += ledNum / split {
					pixels[k] = ledChain[(k+offset)%len(ledChain)]
				}
				time.Sleep(time.Duration(drawDuration) * time.Millisecond)
			}
			offset = (offset + 1) % ledChainLen
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Animation Done!")
	}

	return nil
}
