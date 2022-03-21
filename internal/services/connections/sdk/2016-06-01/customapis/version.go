package customapis

import "fmt"

const defaultApiVersion = "2016-06-01"

func userAgent() string {
	return fmt.Sprintf("pandora/customapis/%s", defaultApiVersion)
}
