package nameavailability

import "fmt"

const defaultApiVersion = "2022-05-13"

func userAgent() string {
	return fmt.Sprintf("pandora/nameavailability/%s", defaultApiVersion)
}
