package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/flags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

var requireResourcesToBeImported = flags.RequireResourcesToBeImported

func tagsSchema() *schema.Schema {
	return tags.Schema()
}

func tagsForceNewSchema() *schema.Schema {
	return tags.ForceNewSchema()
}

func tagsForDataSourceSchema() *schema.Schema {
	return tags.DataSourceSchema()
}

func tagValueToString(v interface{}) (string, error) {
	return tags.TagValueToString(v)
}

func validateAzureRMTags(v interface{}, k string) (warnings []string, errors []error) {
	return tags.Validate(v, k)
}

func expandTags(tagsMap map[string]interface{}) map[string]*string {
	return tags.Expand(tagsMap)
}

func filterTags(tagsMap map[string]*string, tagNames ...string) map[string]*string {
	return tags.Filter(tagsMap, tagNames...)
}

func flattenAndSetTags(d *schema.ResourceData, tagMap map[string]*string) {
	// we intentionally ignore the error here, since this method doesn't expose it
	_ = tags.FlattenAndSet(d, tagMap)
}

// migrated

// nolint: deadcode unused
func parseAzureResourceID(id string) (*azure.ResourceID, error) {
	return azure.ParseAzureResourceID(id)
}

// nolint: deadcode unused
func ignoreCaseDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}

// nolint: deadcode unused
func azureRMLockByName(name string, resourceType string) {
	locks.ByName(name, resourceType)
}

// nolint: deadcode unused
func azureRMLockMultipleByName(names *[]string, resourceType string) {
	locks.MultipleByName(names, resourceType)
}

// nolint: deadcode unused
func azureRMUnlockByName(name string, resourceType string) {
	locks.UnlockByName(name, resourceType)
}

// nolint: deadcode unused
func azureRMUnlockMultipleByName(names *[]string, resourceType string) {
	locks.UnlockMultipleByName(names, resourceType)
}

func validateRFC3339Date(v interface{}, k string) (warnings []string, errors []error) {
	return validate.RFC3339Time(v, k)
}

// nolint: deadcode unused
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
