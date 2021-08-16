package capacities

import "fmt"

const defaultApiVersion = "2021-01-01"

func userAgent() string {
	return fmt.Sprintf("pandora/capacities/%s", defaultApiVersion)
}
