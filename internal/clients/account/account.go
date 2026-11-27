// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package account

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
)

type ResourceManagerAccount struct {
	Environment environments.Environment

	ClientId       string
	ObjectId       string
	SubscriptionId string
	TenantId       string

	AuthenticatedAsAServicePrincipal bool
	RegisteredResourceProviders      resourceproviders.ResourceProviders
}

type Provider interface {
	ClientAccount() *ResourceManagerAccount
	ClientStopContext() context.Context
}

func From(meta any) *ResourceManagerAccount {
	return meta.(Provider).ClientAccount()
}

func StopContext(meta any) context.Context {
	return meta.(Provider).ClientStopContext()
}
