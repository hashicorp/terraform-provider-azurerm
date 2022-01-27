package profiles

import "fmt"

const defaultApiVersion = "2021-06-01"

func userAgent() string {
	return fmt.Sprintf("pandora/profiles/%s", defaultApiVersion)
}
