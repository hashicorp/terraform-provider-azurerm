package disasterrecoveryconfigs

import "fmt"

const defaultApiVersion = "2017-04-01"

func userAgent() string {
	return fmt.Sprintf("pandora/disasterrecoveryconfigs/%s", defaultApiVersion)
}
