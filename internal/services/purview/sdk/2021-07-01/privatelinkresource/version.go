package privatelinkresource

import "fmt"

const defaultApiVersion = "2021-07-01"

func userAgent() string {
	return fmt.Sprintf("pandora/privatelinkresource/%s", defaultApiVersion)
}
