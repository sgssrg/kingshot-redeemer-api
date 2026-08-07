package global

import (
	"os"
	"strconv"
)

// Exported global variable (capitalized)
var KSFetchAvailable bool

// init() runs automatically when the package is imported
func init() {
	val := os.Getenv("KSFetchAvailable")
	b, err := strconv.ParseBool(val)
	if err != nil {
		KSFetchAvailable = false // default
	} else {
		KSFetchAvailable = b
	}
}
