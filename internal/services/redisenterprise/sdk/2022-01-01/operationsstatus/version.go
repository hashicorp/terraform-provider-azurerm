package operationsstatus

import "fmt"

const defaultApiVersion = "2022-01-01"

func userAgent() string {
	return fmt.Sprintf("pandora/operationsstatus/%s", defaultApiVersion)
}
