package main

import (
	"fmt"
	"log"

	"github.com/go-rod/rod"
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
		MustConnect()
	wait := browser.EachEvent(func(e *proto.TargetTargetInfoChanged) {
		page, err := browser.PageFromTarget(e.TargetInfo.TargetID)
		if err != nil {
			log.Fatal("could not get page from target using targetID!")
		}
		{
			_ = proto.OverlaySetShowFPSCounter{
				Show: true,
			}.Call(page)
		}
	}, func(e *proto.PageLoadEventFired) {
		fmt.Println("pageload event fired!")
	})
	wait()
}
