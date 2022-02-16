package healthcare

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	dicomService "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"log"
	"time"
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
			_, err := dicomService.ParseDicomServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				//todo check the validation func
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": azure.SchemaLocation(),

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							//todo: must follow https://login.microsoft.com/tenantid
						},
						"audience": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"smart_proxy_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
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

	workspace, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare workspace error: %+v", err)
	}

	dicomServiceId := dicomService.NewDicomServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, dicomServiceId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presense of existing %s: %+v", dicomServiceId, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_healthcareapis_dicom_service", dicomServiceId.ID())
		}
	}

	//parameters := dicomService.DicomService{
	//	Name: utils.String(dicomServiceId.DicomServiceName),
	//	Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
	//}
	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceHealthcareApisDicomServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
