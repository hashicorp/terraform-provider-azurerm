package eventhubs

import "fmt"

const defaultApiVersion = "2017-04-01"

func userAgent() string {
	return fmt.Sprintf("pandora/eventhubs/%s", defaultApiVersion)
}
