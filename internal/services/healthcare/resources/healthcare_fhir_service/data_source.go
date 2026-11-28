// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

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

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Healthcare FHIR Service (Fast Healthcare Interoperability Resources).",
		Read:        dataSourceHealthcareApisFhirServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Healthcare FHIR Service.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FhirServiceName(),
			},

			"workspace_id": {
				Description:  "The ID of the Healthcare Workspace in which the Healthcare FHIR Service exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.LocationComputed(),

			"kind": {
				Description: "The kind of the Healthcare FHIR Service.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"access_policy_object_ids": {
				Description: "The list of the access policies of the service instance.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
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
							Description: "The Azure Active Directory (tenant) that serves as the authentication authority to access the service. Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id}.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"audience": {
							Description: "The intended audience to receive authentication tokens for the service.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
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
				Description: "The list of azure container registry settings used for convert data operation of the service instance.",
				Type:        pluginsdk.TypeList,
				Computed:    true,
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
							Description: "The set of origins to be allowed via CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_headers": {
							Description: "The set of headers to be allowed via CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"allowed_methods": {
							Description: "The methods to be allowed via CORS.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"max_age_in_seconds": {
							Description: "The max age to be allowed via CORS.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"credentials_allowed": {
							Description: "Whether credentials are allowed via CORS.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"configuration_export_storage_account_name": {
				Description: "The name of the Storage Account which the operation configuration information is exported to.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
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
		d.Set("kind", pointer.FromEnum(m.Kind))

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
