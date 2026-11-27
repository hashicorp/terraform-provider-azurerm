// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	healthcareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Healthcare DICOM Service.",
		Read:        dataSourceHealthcareApisDicomServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Healthcare DICOM Service.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DicomServiceName(),
			},

			"workspace_id": {
				Description:  "The ID of the Healthcare Workspace in which the Healthcare DICOM Service exists.",
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
							Description: "The Azure Active Directory (tenant) that serves as the authentication authority to access the service. Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id}.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"audience": {
							Description: "The intended audience to receive authentication tokens for the service.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
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
				Description: "The URL of the Healthcare DICOM Services.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"data_partitions_enabled": {
				Description: "If data partitions are enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Description: "A list of allowed origins for CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_headers": {
							Description: "A list of allowed headers for CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_methods": {
							Description: "A list of allowed methods for CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"max_age_in_seconds": {
							Description: "The maximum age in seconds for the CORS configuration.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"allow_credentials": {
							Description: "Whether credentials are allowed in CORS.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"encryption_key_url": {
				Description: "The URL of the key to use for encryption as part of the customer-managed key encryption settings.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_system_name": {
							Description: "The filesystem name of connected Storage Account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"storage_account_id": {
							Description: "The resource ID of connected Storage Account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := healthcareClient.From(meta).HealthcareWorkspaceDicomServiceClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
