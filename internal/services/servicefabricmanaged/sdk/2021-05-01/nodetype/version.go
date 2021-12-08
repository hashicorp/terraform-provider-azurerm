package nodetype

import "fmt"

const defaultApiVersion = "2021-05-01"

func userAgent() string {
	return fmt.Sprintf("pandora/nodetype/%s", defaultApiVersion)
}
