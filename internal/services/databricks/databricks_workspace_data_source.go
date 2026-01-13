// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2025-10-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDatabricksWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDatabricksWorkspaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"compute_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"managed_disk_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"storage_account_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enhanced_security_compliance": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"automatic_cluster_update_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"compliance_security_profile_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"compliance_security_profile_standards": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"enhanced_security_monitoring_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"custom_parameters": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"machine_learning_workspace_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"nat_gateway_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"no_public_ip": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"private_subnet_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"public_ip_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"public_subnet_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"storage_account_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"storage_account_sku_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"virtual_network_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"vnet_address_prefix": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

// This functions is used to Flatten the custom parameters in data source.
// It is similar to flattenWorkspaceCustomParameters but does not return the backend address pool ID.
// It also omits the public and private subnet NSG association IDs since they are not available in API.
func flattenWorkspaceCustomParametersForDataSource(input *workspaces.WorkspaceCustomParameters) []interface{} {
	if input == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if v := input.AmlWorkspaceId; v != nil {
		parameters["machine_learning_workspace_id"] = v.Value
	}

	if v := input.NatGatewayName; v != nil {
		parameters["nat_gateway_name"] = v.Value
	}

	if v := input.EnableNoPublicIP; v != nil {
		parameters["no_public_ip"] = v.Value
	}

	if v := input.CustomPrivateSubnetName; v != nil {
		parameters["private_subnet_name"] = v.Value
	}

	if v := input.PublicIPName; v != nil {
		parameters["public_ip_name"] = v.Value
	}

	if v := input.CustomPublicSubnetName; v != nil {
		parameters["public_subnet_name"] = v.Value
	}

	if v := input.StorageAccountName; v != nil {
		parameters["storage_account_name"] = v.Value
	}

	if v := input.StorageAccountSkuName; v != nil {
		parameters["storage_account_sku_name"] = v.Value
	}

	if v := input.CustomVirtualNetworkId; v != nil {
		parameters["virtual_network_id"] = v.Value
	}

	if v := input.VnetAddressPrefix; v != nil {
		parameters["vnet_address_prefix"] = v.Value
	}

	return []interface{}{parameters}
}

func dataSourceDatabricksWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if sku := resp.Model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}
		d.Set("compute_mode", string(model.Properties.ComputeMode))
		d.Set("workspace_id", model.Properties.WorkspaceId)
		if err := d.Set("storage_account_identity", flattenWorkspaceManagedIdentity(model.Properties.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}
		if err := d.Set("managed_disk_identity", flattenWorkspaceManagedIdentity(model.Properties.ManagedDiskIdentity)); err != nil {
			return fmt.Errorf("setting `managed_disk_identity`: %+v", err)
		}
		d.Set("workspace_url", model.Properties.WorkspaceURL)
		d.Set("location", model.Location)
		if err := d.Set("enhanced_security_compliance", flattenWorkspaceEnhancedSecurity(model.Properties.EnhancedSecurityCompliance)); err != nil {
			return fmt.Errorf("setting `enhanced_security_compliance`: %+v", err)
		}

		if err := d.Set("custom_parameters", flattenWorkspaceCustomParametersForDataSource(model.Properties.Parameters)); err != nil {
			return fmt.Errorf("setting `custom_parameters`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
