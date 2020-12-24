package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const xmasTreeShowName = "xmas-tree"

func xmasTreeShow(ctx context.Context, pixels []hsv) error {
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
				go func(ctx context.Context, color hsv) {
					skip[i] = true
					pixels[i] = color
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
									pixels[i].H = colorGreen.H
									pixels[i].S = colorGreen.S
								}
								pixels[i].V = v
								continue
							}

							v = pixels[i].V * (2 - fadeMultiplier)
							if v > colorGreen.V {
								pixels[i].V = colorGreen.V
								skip[i] = false
								ticker.Stop()
								return
							}
							pixels[i].V = v
						}
					}
				}(ctx, []hsv{colorRed, colorBlue, colorPurple, colorPink, colorWhite}[rand.Intn(5)])
			}
		}
	}()

	// Sparkle
	o := colorGreen
	o.S -= 0.25
	go ahSparkleSkipHSV(ctx, pixels, skip, o, 80, 0.5, 0.25, 5000)

	select {
	case <-ctx.Done():
		log.Println("Animation Done!")
	}

	return nil
}

const originalShowName = "original"

func originalShow(ctx context.Context, pixels []hsv) error {
	fmt.Println("Starting Original Show...")

	ledChain := []hsv{colorGreen, colorRed, colorGreen, colorWhite}

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
