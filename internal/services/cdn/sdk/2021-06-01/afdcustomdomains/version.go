package afdcustomdomains

import "fmt"

const defaultApiVersion = "2021-06-01"

func userAgent() string {
	return fmt.Sprintf("pandora/afdcustomdomains/%s", defaultApiVersion)
}
