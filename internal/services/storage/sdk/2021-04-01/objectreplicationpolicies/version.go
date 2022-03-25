package objectreplicationpolicies

import "fmt"

const defaultApiVersion = "2021-04-01"

func userAgent() string {
	return fmt.Sprintf("pandora/objectreplicationpolicies/%s", defaultApiVersion)
}
