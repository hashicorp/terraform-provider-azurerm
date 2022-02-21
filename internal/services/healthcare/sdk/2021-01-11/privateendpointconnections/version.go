package privateendpointconnections

import "fmt"

const defaultApiVersion = "2021-01-11"

func userAgent() string {
	return fmt.Sprintf("pandora/privateendpointconnections/%s", defaultApiVersion)
}
