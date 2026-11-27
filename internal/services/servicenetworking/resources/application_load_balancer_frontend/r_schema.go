// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (f Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The name which should be used for this Application Gateway for Containers Frontend.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_load_balancer_id": {
			Description:  "The ID of the Application Gateway for Containers.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: frontendsinterface.ValidateTrafficControllerID,
		},

		"tags": commonschema.Tags(),
	}
}

func (f Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fully_qualified_domain_name": {
			Description: "The Fully Qualified Domain Name of the DNS record associated to an Application Gateway for Containers Frontend.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
