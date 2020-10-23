package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

/**
 * Entrypoint function for the program
 * Creates a new browser session
 * Listens for Target Changes and Page Loads to display performance metrics that correlate to the page currently in scope
 */
func main() {

	browser := rod.New().
		ControlURL(
			launcher.New().
				Headless(false).
				MustLaunch(),
		).
		MustConnect().DefaultDevice(devices.Clear, true)

	wait := browser.EachEvent(func(t *proto.TargetTargetInfoChanged) {
		quit := make(chan bool)
		if getURL() != t.TargetInfo.URL {
			page, err := browser.PageFromTarget(t.TargetInfo.TargetID)
			explain(err)
			getShowFPSCounterOverlay(page)
			wait2 := page.EachEvent(func(p *proto.PageLoadEventFired) {
				select {
				case <-quit:
					close(quit)
					return
				default:
					getPageLoadTimingsOverlay(page)
					quit <- true
					break
				}
			})
			go wait2()
			setURL(t.TargetInfo.URL)
		}
	})
	wait()
}
