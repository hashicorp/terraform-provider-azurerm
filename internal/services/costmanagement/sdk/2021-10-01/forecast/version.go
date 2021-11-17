package forecast

import "fmt"

const defaultApiVersion = "2021-10-01"

func userAgent() string {
	return fmt.Sprintf("pandora/forecast/%s", defaultApiVersion)
}
