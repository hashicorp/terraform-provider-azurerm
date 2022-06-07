package confidentialledger

import "fmt"

const defaultApiVersion = "2022-05-13"

func userAgent() string {
	return fmt.Sprintf("pandora/confidentialledger/%s", defaultApiVersion)
}
