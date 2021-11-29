package communicationservice

import "fmt"

const defaultApiVersion = "2020-08-20"

func userAgent() string {
	return fmt.Sprintf("pandora/communicationservice/%s", defaultApiVersion)
}
