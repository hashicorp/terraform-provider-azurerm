// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceHealthcareApisFhirService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHealthcareApisFhirServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FhirServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.LocationComputed(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_policy_object_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"audience": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"smart_proxy_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"container_registry_login_server_url": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_headers": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"max_age_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"credentials_allowed": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"configuration_export_storage_account_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing workspace id error: %+v", err)
	}
	id := fhirservices.NewFhirServiceID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.WorkspaceName)
	d.Set("workspace_id", workspaceId.ID())

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))
		d.Set("kind", string(pointer.From(m.Kind)))

		if props := m.Properties; props != nil {
			d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
			d.Set("authentication", flattenFhirAuthentication(props.AuthenticationConfiguration))
			d.Set("cors", flattenFhirCorsConfiguration(props.CorsConfiguration))
			d.Set("container_registry_login_server_url", flattenFhirAcrLoginServer(props.AcrConfiguration))
			if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
				d.Set("configuration_export_storage_account_name", props.ExportConfiguration.StorageAccountName)
			}
		}

		i, err := identity.FlattenLegacySystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, m.Tags)
	}

	return nil
}
