package endpoints

import "fmt"

const defaultApiVersion = "2018-08-01"

func userAgent() string {
	return fmt.Sprintf("pandora/endpoints/%s", defaultApiVersion)
}
