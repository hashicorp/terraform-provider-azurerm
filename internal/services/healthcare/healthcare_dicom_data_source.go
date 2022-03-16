package healthcare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceHealthcareApisDicomService() *pluginsdk.Resource {
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

			"location": commonschema.LocationComputed(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"authentication_configuration": {
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

			"service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		d.Set("location", azure.NormalizeLocation(*resp.Location))
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
