// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionruleassociations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule_association"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionruleassociations.ValidateScopedDataCollectionRuleAssociationID
}
