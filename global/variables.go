package global

import (
	"os"
	"strconv"
)

// Exported global variable (capitalized)
var KSFetchAvailable bool
var StratForgeFetchAvailable bool

// init() runs automatically when the package is imported
func init() {
	val := os.Getenv("KSFetchAvailable")
	a, err := strconv.ParseBool(val)
	if err != nil {
		KSFetchAvailable = false // default
	} else {
		KSFetchAvailable = a
	}

	StratForgeFetchAvailableEnv := os.Getenv("StratForgeFetchAvailable")
	b, err := strconv.ParseBool(StratForgeFetchAvailableEnv)
	if err != nil {
		StratForgeFetchAvailable = true // default
	} else {
		StratForgeFetchAvailable = b
	}
}
