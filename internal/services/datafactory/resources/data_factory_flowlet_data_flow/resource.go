// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_flowlet_data_flow

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Read:   resourceDataFactoryFlowletDataFlowRead,
		Update: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Delete: resourceDataFactoryFlowletDataFlowDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dataflows.ParseDataflowID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: dataFactoryFlowletDataFlowSchema(),
	}
}
