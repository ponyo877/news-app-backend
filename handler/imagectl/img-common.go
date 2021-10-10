package imagectl

import (
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
	    fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}
