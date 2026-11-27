// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_elastic_application_performance_monitoring

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

func (s Resource) DeprecationMessage() string {
	return "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_elastic_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information."
}

var (
	_ sdk.ResourceWithUpdate                      = Resource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = Resource{}
)

func (s Resource) ResourceType() string {
	return "azurerm_spring_cloud_elastic_application_performance_monitoring"
}

func (s Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}
