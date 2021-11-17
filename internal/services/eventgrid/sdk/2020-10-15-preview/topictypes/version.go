package topictypes

import "fmt"

const defaultApiVersion = "2020-10-15-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/topictypes/%s", defaultApiVersion)
}
