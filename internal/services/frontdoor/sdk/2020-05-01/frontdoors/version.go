package frontdoors

import "fmt"

const defaultApiVersion = "2020-05-01"

func userAgent() string {
	return fmt.Sprintf("pandora/frontdoors/%s", defaultApiVersion)
}
