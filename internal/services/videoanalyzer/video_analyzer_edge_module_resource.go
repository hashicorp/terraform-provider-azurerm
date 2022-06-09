package videoanalyzer

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/sdk/2021-05-01-preview/videoanalyzer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/videoanalyzer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := videoanalyzer.ParseEdgeModuleID(id)
			return err
		}),

		DeprecationMessage: `Video Analyzer (Preview) is now Deprecated and will be Retired on 2022-11-30 - as such the 'azurerm_video_analyzer_edge_module' resource is deprecated and will be removed in v4.0 of the AzureRM Provider`,

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
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id := videoanalyzer.NewEdgeModuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("video_analyzer_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.EdgeModulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_video_analyzer_edge_module", id.ID())
		}
	}

	if _, err := client.EdgeModulesCreateOrUpdate(ctx, id, videoanalyzer.EdgeModuleEntity{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVideoAnalyzerEdgeModuleRead(d, meta)
}

func resourceVideoAnalyzerEdgeModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := videoanalyzer.ParseEdgeModuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.EdgeModulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.EdgeModuleName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("video_analyzer_name", id.AccountName)

	return nil
}

func resourceVideoAnalyzerEdgeModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).VideoAnalyzer.VideoAnalyzersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := videoanalyzer.ParseEdgeModuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.EdgeModulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
