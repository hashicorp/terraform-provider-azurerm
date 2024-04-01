// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/app-service"
}

func (r Registration) WebsiteCategories() []string {
	return nil
}

func (r Registration) Name() string {
	return "AppService"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AppServiceEnvironmentV3DataSource{},
		AppServiceSourceControlTokenDataSource{},
		LinuxFunctionAppDataSource{},
		LinuxWebAppDataSource{},
		ServicePlanDataSource{},
		StaticWebAppDataSource{},
		WindowsFunctionAppDataSource{},
		WindowsWebAppDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AppServiceEnvironmentV3Resource{},
		AppServiceSourceControlTokenResource{},
		FunctionAppActiveSlotResource{},
		FunctionAppFunctionResource{},
		FunctionAppHybridConnectionResource{},
		LinuxFunctionAppResource{},
		LinuxFunctionAppSlotResource{},
		LinuxWebAppResource{},
		LinuxWebAppSlotResource{},
		ServicePlanResource{},
		SourceControlResource{},
		SourceControlSlotResource{},
		StaticWebAppResource{},
		StaticWebAppCustomDomainResource{},
		StaticWebAppFunctionAppRegistrationResource{},
		WebAppActiveSlotResource{},
		WebAppHybridConnectionResource{},
		WindowsFunctionAppResource{},
		WindowsFunctionAppSlotResource{},
		WindowsWebAppResource{},
		WindowsWebAppSlotResource{},
	}
}
