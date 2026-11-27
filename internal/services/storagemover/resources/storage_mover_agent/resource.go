// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_agent

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/agents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_mover_agent -service-package-name storagemover -properties "name" -compare-values "subscription_id:storage_mover_id,resource_group_name:storage_mover_id,storage_mover_name:storage_mover_id"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &agents.AgentId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_storage_mover_agent"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return agents.ValidateAgentID
}
