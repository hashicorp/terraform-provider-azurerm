package collection

import "fmt"

const defaultApiVersion = "2021-01-11"

func userAgent() string {
	return fmt.Sprintf("pandora/collection/%s", defaultApiVersion)
}
