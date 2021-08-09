package videoanalyzer

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/videoanalyzer/mgmt/2021-05-01-preview/videoanalyzer"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVideoAnalyzerEdgeModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVideoAnalyzerEdgeModuleCreateUpdate,
		Read:   resourceVideoAnalyzerEdgeModuleRead,
		Delete: resourceVideoAnalyzerEdgeModuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EdgeModulesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,30}[a-zA-Z0-9])$`),
					"Video Analyzer Edge Module name must be 1 - 32 characters long, begin and end with a letter or number and may contain only letters, numbers or underscore.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"video_analyzer_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VideoAnalyzerName(),
			},
		},
	}
}

func resourceVideoAnalyzerEdgeModuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.EdgeModulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	resourceId := parse.NewEdgeModulesID(subscriptionId, d.Get("resource_group_name").(string), d.Get("video_analyzer_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.VideoAnalyzerName, resourceId.EdgeModuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_video_analyzer_edge_module", resourceId.ID())
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.VideoAnalyzerName, resourceId.EdgeModuleName, videoanalyzer.EdgeModuleEntity{}); err != nil {
		return fmt.Errorf("creating %s: %+v", resourceId, err)
	}

	d.SetId(resourceId.ID())
	return resourceVideoAnalyzerEdgeModuleRead(d, meta)
}

func resourceVideoAnalyzerEdgeModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.EdgeModulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EdgeModulesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VideoAnalyzerName, id.EdgeModuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Video Analyzer Edge Module %q was not found in Resource Group %q - removing from state", id.EdgeModuleName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Video Analyzer Edge Module %q (Resource Group %q): %+v", id.EdgeModuleName, id.ResourceGroup, err)
	}

	d.Set("name", id.EdgeModuleName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("video_analyzer_name", id.VideoAnalyzerName)

	return nil
}

func resourceVideoAnalyzerEdgeModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.EdgeModulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EdgeModulesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.VideoAnalyzerName, id.EdgeModuleName)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM delete request for Video Analyzer Edge Module '%s': %+v", id.EdgeModuleName, err)
	}

	return nil
}
