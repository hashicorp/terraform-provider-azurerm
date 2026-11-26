// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_custom_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func expandDataFactoryLinkedServiceIntegrationRuntimeV2(input []any) *datafactory.IntegrationRuntimeReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	return &datafactory.IntegrationRuntimeReference{
		ReferenceName: pointer.To(v["name"].(string)),
		Type:          pointer.To("IntegrationRuntimeReference"),
		Parameters:    v["parameters"].(map[string]any),
	}
}
