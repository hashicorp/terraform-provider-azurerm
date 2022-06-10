package eventhubsclustersnamespace

import "fmt"

const defaultApiVersion = "2021-11-01"

func userAgent() string {
	return fmt.Sprintf("pandora/eventhubsclustersnamespace/%s", defaultApiVersion)
}
