package servers

import "fmt"

const defaultApiVersion = "2017-08-01"

func userAgent() string {
	return fmt.Sprintf("pandora/servers/%s", defaultApiVersion)
}
