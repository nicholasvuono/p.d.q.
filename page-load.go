package main

import (
	"fmt"
	"strconv"

	"github.com/go-rod/rod"
	"github.com/tidwall/gjson"
)

//Gets First Contentful Paint information for a given page and returns the value as a string
func getFirstContentfulPaint(page *rod.Page) string {
	firstContentfulPaint, err := page.Eval("JSON.stringify(performance.getEntriesByName('first-contentful-paint'))")
	explain(err)
	fcp := gjson.Get(firstContentfulPaint.Value.String(), "0.startTime").Uint()
	return strconv.FormatUint(fcp, 10)
}

//Gets Time to Interactive information for a given page and returns the value as a string
func getTimeToInteractive(page *rod.Page) string {
	metrics, err := page.Eval("JSON.stringify(performance.toJSON())")
	explain(err)
	loadEventEnd := gjson.Get(metrics.Value.String(), "timing.loadEventEnd").Uint()
	navigationStart := gjson.Get(metrics.Value.String(), "timing.navigationStart").Uint()
	timeToInteractive := loadEventEnd - navigationStart
	return strconv.FormatUint(timeToInteractive, 10)
}

//Gets Total Blocking Time information for a given page and returns the value as a string
func getTotalBlockingTime(firstContentfulPaint string, timeToInteractive string) string {
	fcp, err := strconv.Atoi(firstContentfulPaint)
	explain(err)
	tti, err := strconv.Atoi(timeToInteractive)
	explain(err)
	return strconv.Itoa(tti - fcp)
}

////Gets Time to First Byte information for a given page and returns the value as a string
func getTimeToFirstByte(page *rod.Page) string {
	metrics, err := page.Eval("JSON.stringify(performance.toJSON())")
	explain(err)
	requestStart := gjson.Get(metrics.Value.String(), "timing.requestStart").Uint()
	responseStart := gjson.Get(metrics.Value.String(), "timing.responseStart").Uint()
	timeToFirstByte := responseStart - requestStart
	return strconv.FormatUint(timeToFirstByte, 10)
}

/**
 * Calls functions to gather four different important performance page laod time metrics
 * Formats the metrics into an easily readable multiline string and returns it
 */
func getPageLoadTimings(page *rod.Page) string {
	fcp := getFirstContentfulPaint(page)
	tti := getTimeToInteractive(page)
	tbt := getTotalBlockingTime(fcp, tti)
	ttfb := getTimeToFirstByte(page)
	pageLoadTimings := fmt.Sprintf(`
	First Contentful Paint (FCP):		%vms
	Time to Interactive (TTI):			%vms
	Total Blocking Time (TBT):		%vms
	Time to First Byte (TTFB):		%vms
	`, fcp, tti, tbt, ttfb)
	return pageLoadTimings
}

/**
 * Calls function to gather page laod timing metrics
 * Creates a div to display these metrics as an overlay on the page
 */
func getPageLoadTimingsOverlay(page *rod.Page) {
	pageLoadTimings := getPageLoadTimings(page)
	inlineJS := `(pageLoadTimings) => {let existing = document.getElementById('PageLoadTimings');if(existing === null) {let body = document.querySelector('body');let div = document.createElement('div');div.setAttribute('id', 'PageLoadTimings');div.setAttribute('style', 'font-family:Segoe UI;color:#00FFFF;position:fixed;background:rgba(0,0,0,.8);z-index:99000000000;left:0px;bottom:0px;width:300px;height:100px;pointer-events:none;font-size:small;white-space:pre-wrap;');div.textContent = pageLoadTimings;body.appendChild(div);} else {existing.textContent = pageLoadTimings;}}`
	page.MustEval(inlineJS, pageLoadTimings)
}
