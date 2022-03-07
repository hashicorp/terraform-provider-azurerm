package connectiongateways

import "fmt"

const defaultApiVersion = "2016-06-01"

func userAgent() string {
	return fmt.Sprintf("pandora/connectiongateways/%s", defaultApiVersion)
}
