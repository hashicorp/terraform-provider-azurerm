package authorizationruleseventhubs

import "fmt"

const defaultApiVersion = "2017-04-01"

func userAgent() string {
	return fmt.Sprintf("pandora/authorizationruleseventhubs/%s", defaultApiVersion)
}
