package roles

import "fmt"

const defaultApiVersion = "2020-10-05-privatepreview"

func userAgent() string {
	return fmt.Sprintf("pandora/roles/%s", defaultApiVersion)
}
