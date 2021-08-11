package accounts

import "fmt"

const defaultApiVersion = "2021-02-01"

func userAgent() string {
	return fmt.Sprintf("pandora/accounts/%s", defaultApiVersion)
}
