// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_workflow

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var logicAppResourceName = "azurerm_logic_app"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppWorkflowCreate,
		Read:   resourceLogicAppWorkflowRead,
		Update: resourceLogicAppWorkflowUpdate,
		Delete: resourceLogicAppWorkflowDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workflows.ParseWorkflowID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: logicAppWorkflowSchema(),
	}
}
