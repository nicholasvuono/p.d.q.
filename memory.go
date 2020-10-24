package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/go-rod/rod"
)

//Gets the currently used memory percentage
func getUsedMemoryPercentage(page *rod.Page) string {
	totalJSHeapSize, err := page.Eval("JSON.stringify(window.performance.memory.totalJSHeapSize)")
	explain(err)
	usedJSHeapSize, err := page.Eval("JSON.stringify(window.performance.memory.usedJSHeapSize)")
	explain(err)
	total, err := strconv.Atoi(totalJSHeapSize.Value.String())
	explain(err)
	used, err := strconv.Atoi(usedJSHeapSize.Value.String())
	explain(err)
	memory := strconv.Itoa(int(math.Ceil((float64(used) / float64(total)) * 100)))
	return memory
}

//Gets the current number of DOM nodes
func getCurrentNumberOfNodes(page *rod.Page) string {
	nodes, err := page.Eval("JSON.stringify(document.getElementsByTagName('*').length)")
	explain(err)
	return nodes.Value.String()
}

//Gets the current page url
func getPageURL(page *rod.Page) string {
	url, err := page.Eval("JSON.stringify(document.URL)")
	explain(err)
	return url.Value.String()
}

//Prints memory metrics into the terminal
func logMemoryUsage(page *rod.Page) {
	url := getPageURL(page)
	nodes := getCurrentNumberOfNodes(page)
	percentage := getUsedMemoryPercentage(page)
	fmt.Println(fmt.Sprintf(`
	Page URL: %v

		- DOM Nodes: %v
		- Memory Percentage: %v 
	`, url, nodes, percentage))
}
