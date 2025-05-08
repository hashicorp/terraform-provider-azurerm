// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	mgValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceManagementGroupTemplateDeployment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagementGroupTemplateDeploymentRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateDeploymentName,
			},

			"management_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: mgValidate.ManagementGroupID,
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

func dataSourceManagementGroupTemplateDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LegacyDeploymentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId := d.Get("management_group_id").(string)
	deploymentName := d.Get("name").(string)

	mgmtGroupId, err := mgParse.ManagementGroupID(managementGroupId)
	if err != nil {
		return err
	}

	resp, err := client.GetAtManagementGroupScope(ctx, mgmtGroupId.Name, deploymentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("deployment %s in Management Group %s was not found", deploymentName, managementGroupId)
		}

		return fmt.Errorf("retrieving Management Group Template Deployment %s in management group %s: %+v", deploymentName, managementGroupId, err)
	}

	templateId, err := parse.ManagementGroupTemplateDeploymentID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(templateId.ID())

	if props := resp.Properties; props != nil {
		flattenedOutputs, err := flattenTemplateDeploymentBody(props.Outputs)
		if err != nil {
			return fmt.Errorf("flattening `output_content`: %+v", err)
		}
		return d.Set("output_content", flattenedOutputs)
	}

	return nil
}
