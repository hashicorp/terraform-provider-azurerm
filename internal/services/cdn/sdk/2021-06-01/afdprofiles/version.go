package afdprofiles

import "fmt"

const defaultApiVersion = "2021-06-01"

func userAgent() string {
	return fmt.Sprintf("pandora/afdprofiles/%s", defaultApiVersion)
}
