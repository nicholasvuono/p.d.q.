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
	launcher := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(launcher).MustConnect().DefaultDevice(devices.Clear)
	return browser
}

//Listens for a page load event and gets pageload metrics and memory usage
func onPageLoadEventFired(page *rod.Page, quit chan bool) func() {
	wait := page.EachEvent(
		func(p *proto.PageLoadEventFired) {
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
	wait := browser.EachEvent(
		func(t *proto.TargetTargetInfoChanged) {
			quit := make(chan bool)
			if url != t.TargetInfo.URL {
				page, err := browser.PageFromTarget(t.TargetInfo.TargetID)
				explain(err)
				getShowFPSCounterOverlay(page)
				w := onPageLoadEventFired(page, quit)
				go w()
				url = t.TargetInfo.URL
			}
		})
	return wait
}