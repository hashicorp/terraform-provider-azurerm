package healthcare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				ValidateFunc: validate.WorkspaceID,
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

			"identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"acr_login_servers": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"cors_configuration": {
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

			"export_storage_account_name": {
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

	workspaceId, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing workspace id error: %+v", err)
	}
	id := parse.NewFhirServiceID(subscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)

	d.Set("workspace_id", workspaceId.ID())

	if resp.Location != nil {
		d.Set("location", location.NormalizeNilable(resp.Location))
	}
	if err := d.Set("identity", flattenFhirManagedIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	if err := d.Set("kind", resp.Kind); err != nil {
		return fmt.Errorf("setting `kind`: %+v", err)
	}

	if props := resp.FhirServiceProperties; props != nil {
		d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
		d.Set("authentication", flattenFhirAuthentication(props.AuthenticationConfiguration))
		d.Set("cors_configuration", flattenFhirCorsConfiguration(props.CorsConfiguration))
		d.Set("acr_login_servers", flattenFhirAcrLoginServer(props.AcrConfiguration))
		if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
			d.Set("export_storage_account_name", props.ExportConfiguration.StorageAccountName)
		}
	}
	if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
		return err
	}

	return nil
}
