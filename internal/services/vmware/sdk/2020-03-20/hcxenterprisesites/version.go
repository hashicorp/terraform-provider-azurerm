package hcxenterprisesites

import "fmt"

const defaultApiVersion = "2020-03-20"

func userAgent() string {
	return fmt.Sprintf("pandora/hcxenterprisesites/%s", defaultApiVersion)
}
