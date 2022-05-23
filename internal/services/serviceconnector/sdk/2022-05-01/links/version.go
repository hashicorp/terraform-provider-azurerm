package links

import "fmt"

const defaultApiVersion = "2022-05-01"

func userAgent() string {
	return fmt.Sprintf("pandora/links/%s", defaultApiVersion)
}
