package privatelinkresources

import "fmt"

const defaultApiVersion = "2021-04-30"

func userAgent() string {
	return fmt.Sprintf("pandora/privatelinkresources/%s", defaultApiVersion)
}
