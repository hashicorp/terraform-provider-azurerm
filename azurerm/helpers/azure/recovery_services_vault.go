package azure

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"
)

func ValidateRecoveryServicesVaultName(v interface{}, k string) (warnings []string, errors []error) {
	regexpValidator := validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,49}$"),
		"Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens.",
	)
	return regexpValidator(v, k)
}

// This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
func HandleAzureSdkForGoBug2824(id string) string {
	return strings.Replace(id, "/Subscriptions/", "/subscriptions/", 1)
}
