package eventhubsclusters

import "fmt"

const defaultApiVersion = "2018-01-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/eventhubsclusters/%s", defaultApiVersion)
}
