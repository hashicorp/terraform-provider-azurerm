package privatelinkresources

import "fmt"

const defaultApiVersion = "2022-01-01"

func userAgent() string {
	return fmt.Sprintf("pandora/privatelinkresources/%s", defaultApiVersion)
}
