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
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceHealthcareDicomService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceHealthcareApisDicomServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DicomServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

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
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},

			"private_endpoint": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_partitions_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
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

						"allow_credentials": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"encryption_key_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_system_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"storage_account_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing workspace id error: %+v", err)
	}

	id := dicomservices.NewDicomServiceID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
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

		if props := m.Properties; props != nil {
			d.Set("authentication", flattenDicomAuthentication(props.AuthenticationConfiguration))
			d.Set("private_endpoint", flattenDicomServicePrivateEndpoint(props.PrivateEndpointConnections))
			d.Set("service_url", props.ServiceURL)

			d.Set("data_partitions_enabled", pointer.From(props.EnableDataPartitions))

			d.Set("cors", flattenDicomServiceCorsConfiguration(props.CorsConfiguration))

			if props.Encryption != nil && props.Encryption.CustomerManagedKeyEncryption != nil {
				d.Set("encryption_key_url", pointer.From(props.Encryption.CustomerManagedKeyEncryption.KeyEncryptionKeyURL))
			}

			storage, err := flattenStorageConfiguration(props.StorageConfiguration)
			if err != nil {
				return fmt.Errorf("flattening `storage`: %+v", err)
			}
			d.Set("storage", storage)
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
