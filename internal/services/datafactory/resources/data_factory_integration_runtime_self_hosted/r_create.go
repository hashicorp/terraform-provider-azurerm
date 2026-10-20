// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_self_hosted

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
)

func expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesLinkedInfo(input []any) *integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization {
	if len(input) == 0 {
		return nil
	}

	rbacConfig := input[0].(map[string]any)
	rbac := rbacConfig["resource_id"].(string)
	return &integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization{
		ResourceId:        rbac,
		AuthorizationType: string(helpers.AuthorizationTypeRBAC),
	}
}
