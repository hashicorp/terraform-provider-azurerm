package services

import "fmt"

const defaultApiVersion = "2021-05-01"

func userAgent() string {
	return fmt.Sprintf("pandora/services/%s", defaultApiVersion)
}
