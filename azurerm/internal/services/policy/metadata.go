package policy

import (
	"encoding/json"
	"reflect"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func metadataDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	var oldPolicyAssignmentsMetadata map[string]interface{}
	errOld := json.Unmarshal([]byte(old), &oldPolicyAssignmentsMetadata)
	if errOld != nil {
		return false
	}

	var newPolicyAssignmentsMetadata map[string]interface{}
	if new != "" {
		errNew := json.Unmarshal([]byte(new), &newPolicyAssignmentsMetadata)
		if errNew != nil {
			return false
		}
	}

	// Ignore the following keys if they're found in the metadata JSON
	ignoreKeys := [5]string{"assignedBy", "createdBy", "createdOn", "updatedBy", "updatedOn"}
	for _, key := range ignoreKeys {
		delete(oldPolicyAssignmentsMetadata, key)
		delete(newPolicyAssignmentsMetadata, key)
	}

	return reflect.DeepEqual(oldPolicyAssignmentsMetadata, newPolicyAssignmentsMetadata)
}

func metadataSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Optional:         true,
		Computed:         true,
		ValidateFunc:     validation.StringIsJSON,
		DiffSuppressFunc: metadataDiffSuppressFunc,
	}
}
