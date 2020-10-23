package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var url string

/**
 * Helper function that easily explains errors by logging them in a readable format
 * Takes in the runtime caller information and displays it in the format below:
 * [Jan-02-06 3:04pm] Error Warning: example.go main() line:9 Error invalid argument
 */
func explain(err error) {
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		function := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
		fmt.Println("[" + time.Now().Format("Jan-02-06 3:04pm") + "] Error Warning:" + file + " " + function + "() line:" + strconv.Itoa(line) + " " + err.Error())
	}
}

func setURL(u string) {
	url = u
}

func getURL() string {
	u := url
	return u
}
