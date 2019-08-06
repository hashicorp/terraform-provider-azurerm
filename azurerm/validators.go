package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func validateRFC3339Date(v interface{}, k string) (warnings []string, errors []error) {
	return validate.RFC3339Time(v, k)
}

func validateUUID(v interface{}, k string) (warnings []string, errors []error) {
	return validate.UUID(v, k)
}

func validateIso8601Duration() schema.SchemaValidateFunc {
	return validate.ISO8601Duration
}

func validateAzureVirtualMachineTimeZone() schema.SchemaValidateFunc {
	return validate.VirtualMachineTimeZone()
}

func validateCollation() schema.SchemaValidateFunc {
	return validate.DatabaseCollation
}
