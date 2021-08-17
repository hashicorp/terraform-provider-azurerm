package privateendpointconnection

import "fmt"

const defaultApiVersion = "2020-12-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/privateendpointconnection/%s", defaultApiVersion)
}
