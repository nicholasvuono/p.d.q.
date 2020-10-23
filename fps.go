package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func getShowFPSCounterOverlay(page *rod.Page) {
	_ = proto.OverlaySetShowFPSCounter{
		Show: true,
	}.Call(page)
}
