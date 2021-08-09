package managedidentity

import "fmt"

const defaultApiVersion = "2018-11-30"

func userAgent() string {
	return fmt.Sprintf("pandora/managedidentity/%s", defaultApiVersion)
}
