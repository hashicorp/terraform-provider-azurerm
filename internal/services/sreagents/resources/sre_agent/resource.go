// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sre_agent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-07-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_sre_agent"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return agents.ValidateAgentID
}

func (Resource) Identity() resourceids.ResourceId {
	return &agents.AgentId{}
}
