package deploymentinfo

import "fmt"

const defaultApiVersion = "2020-07-01"

func userAgent() string {
	return fmt.Sprintf("pandora/deploymentinfo/%s", defaultApiVersion)
}
