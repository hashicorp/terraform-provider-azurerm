package loadtests

import "fmt"

const defaultApiVersion = "2021-12-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/loadtests/%s", defaultApiVersion)
}
