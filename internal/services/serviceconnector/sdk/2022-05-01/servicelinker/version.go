package servicelinker

import "fmt"

const defaultApiVersion = "2022-05-01"

func userAgent() string {
	return fmt.Sprintf("pandora/servicelinker/%s", defaultApiVersion)
}
