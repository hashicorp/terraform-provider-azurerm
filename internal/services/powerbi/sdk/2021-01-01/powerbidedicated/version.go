package powerbidedicated

import "fmt"

const defaultApiVersion = "2021-01-01"

func userAgent() string {
	return fmt.Sprintf("pandora/powerbidedicated/%s", defaultApiVersion)
}
