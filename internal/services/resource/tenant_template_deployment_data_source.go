// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceTenantTemplateDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceTenantTemplateDeploymentRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateDeploymentName,
			},

			// Computed
			"output_content": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				// NOTE:  outputs can be strings, ints, objects etc - whilst using a nested object was considered
				// parsing the JSON using `jsondecode` allows the users to interact with/map objects as required
			},
		},
	}
}

func dataSourceTenantTemplateDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewTenantTemplateDeploymentID(d.Get("name").(string))

	resp, err := client.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("template %s in tenant was not found", id.DeploymentName)
		}

		return fmt.Errorf("retrieving Tenant Template Deployment %q: %+v", id.DeploymentName, err)
	}

	d.SetId(id.ID())

	if props := resp.Properties; props != nil {
		flattenedOutputs, err := flattenTemplateDeploymentBody(props.Outputs)
		if err != nil {
			return fmt.Errorf("flattening `output_content`: %+v", err)
		}
		return d.Set("output_content", flattenedOutputs)
	}

	return nil
}
