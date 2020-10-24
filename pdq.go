package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

//Global url variable that will be used to keep track of the current page url for comparison
var url string

//Creates a new browser, sets headless to be false, and clears the defualt device to laod the page correctly
func newBrowser() *rod.Browser {
	browser := rod.New().
		ControlURL(
			launcher.New().
				Headless(false).
				MustLaunch(),
		).
		MustConnect().DefaultDevice(devices.Clear, true)
	return browser
}

//Listens for a page load event and gets pageload metrics and memory usage
func onPageLoadEventFired(page *rod.Page, quit chan bool) func() {
	wait := page.EachEvent(func(p *proto.PageLoadEventFired) {
		select {
		case <-quit:
			close(quit)
			return
		default:
			getPageLoadTimingsOverlay(page)
			logMemoryUsage(page)
			quit <- true
			break
		}
	})
	return wait
}

//Listens for a target change event, gets FPS overlay, and waits for page load events
func onTargetInfoChanged(browser *rod.Browser) func() {
	wait := browser.EachEvent(func(t *proto.TargetTargetInfoChanged) {
		quit := make(chan bool)
		if getURL() != t.TargetInfo.URL {
			page, err := browser.PageFromTarget(t.TargetInfo.TargetID)
			explain(err)
			getShowFPSCounterOverlay(page)
			w := onPageLoadEventFired(page, quit)
			go w()
			setURL(t.TargetInfo.URL)
		}
	})
	return wait
}

//Sets a global url variable
func setURL(u string) {
	url = u
}

//Gets the value of the global url variable
func getURL() string {
	u := url
	return u
}

/**
 * Entrypoint function for the program
 * Creates a new browser session
 * Listens for Target Changes and Page Loads to display performance metrics that correlate to the page currently in scope
 */
func main() {
	browser := newBrowser()
	wait := onTargetInfoChanged(browser)
	wait()
}
