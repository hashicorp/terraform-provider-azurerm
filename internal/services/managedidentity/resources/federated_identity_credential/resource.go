// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package federated_identity_credential

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name federated_identity_credential -service-package-name managedidentity -properties "name" -compare-values "subscription_id:user_assigned_identity_id,resource_group_name:user_assigned_identity_id,user_assigned_identity_name:user_assigned_identity_id" -test-sequential

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &federatedidentitycredentials.FederatedIdentityCredentialId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return federatedidentitycredentials.ValidateFederatedIdentityCredentialID
}

func (r Resource) ResourceType() string {
	return "azurerm_federated_identity_credential"
}
