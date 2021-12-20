package resourceskus

import "fmt"

const defaultApiVersion = "2021-08-01"

func userAgent() string {
	return fmt.Sprintf("pandora/resourceskus/%s", defaultApiVersion)
}
