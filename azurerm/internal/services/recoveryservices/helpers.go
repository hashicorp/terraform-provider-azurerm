package recoveryservices

import (
	"strings"
)

// This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
func handleAzureSdkForGoBug2824(id string) string {
	return strings.Replace(id, "/Subscriptions/", "/subscriptions/", 1)
}
