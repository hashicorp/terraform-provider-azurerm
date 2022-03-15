package healthcare

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisDicomService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisDicomServiceCreateUpdate,
		Read:   resourceHealthcareApisDicomServiceRead,
		Update: resourceHealthcareApisDicomServiceCreateUpdate,
		Delete: resourceHealthcareApisDicomServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DicomServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DicomServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"location": commonschema.Location(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"audience": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},

			"private_endpoint_connection": {
				Type:     pluginsdk.TypeSet,
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisDicomServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Dicom Service creation.")

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare workspace error: %+v", err)
	}

	dicomServiceId := parse.NewDicomServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, dicomServiceId.ResourceGroup, dicomServiceId.WorkspaceName, dicomServiceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presense of existing %s: %+v", dicomServiceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcareapis_dicom_service", dicomServiceId.ID())
		}
	}

	t := d.Get("tags").(map[string]interface{})

	publicNetworkAccess := healthcareapis.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = healthcareapis.PublicNetworkAccessDisabled
	}

	parameters := healthcareapis.DicomService{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(t),
		DicomServiceProperties: &healthcareapis.DicomServiceProperties{
			PublicNetworkAccess: publicNetworkAccess,
		},
	}

	future, err := client.CreateOrUpdate(ctx, dicomServiceId.ResourceGroup, dicomServiceId.WorkspaceName, dicomServiceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", dicomServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update %s: %+v", dicomServiceId, err)
	}

	d.SetId(dicomServiceId.ID())
	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("workspace_id", workspaceId.ID())

	if resp.Location != nil {
		d.Set("location", location.Normalize(*resp.Location))
	}

	if props := resp.DicomServiceProperties; props != nil {
		d.Set("authentication_configuration", flattenDicomAuthentication(props.AuthenticationConfiguration))
		d.Set("private_endpoint_connection", flattenDicomServicePrivateEndpoint(props.PrivateEndpointConnections))
		d.Set("service_url", props.ServiceURL)
	}

	if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
		return err
	}

	return nil
}

func resourceHealthcareApisDicomServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, id.WorkspaceName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func flattenDicomAuthentication(input *healthcareapis.DicomServiceAuthenticationConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	authBlock := make(map[string]interface{})
	if input.Authority != nil {
		authBlock["authority"] = *input.Authority
	}

	audience := make([]interface{}, 0)
	if input.Audiences != nil {
		for _, data := range *input.Audiences {
			audience = append(audience, data)
		}
	}
	authBlock["audience"] = pluginsdk.NewSet(pluginsdk.HashString, audience)

	return []interface{}{authBlock}
}

func flattenDicomServicePrivateEndpoint(input *[]healthcareapis.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]interface{}{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}

		if endpoint.ID != nil {
			result["id"] = *endpoint.ID
		}
	}
	return results
}
