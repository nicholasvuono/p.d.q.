package main

/**
 * Entrypoint function for the program
 * Creates a new browser session
 * Listens for Target Changes and Page Loads to display performance metrics that correlate to the page currently in scope
 */
 func main() {
	browser := newBrowser()
	browser.MustPage("https://wosp.io")
	wait := onTargetInfoChanged(browser)
	wait()
}