//https://docs.microsoft.com/en-us/rest/api/media/privateendpointconnections/createorupdate

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaPrivateEndpointConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaPrivateEndpointConnectionCreate,
		Read:   resourceMediaPrivateEndpointConnectionRead,
		//Update: resourceMediaPrivateEndpointConnectionUpdate,
		Delete: resourceMediaPrivateEndpointConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			//Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PrivateEndpointConnectionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Private Endpoint Connection name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateMediaServicesAccountName,
			},

			"private_link_connection_state": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"actions_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMediaPrivateEndpointConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.PrivateEndpointConnectionsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewPrivateEndpointConnectionID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) || !utils.ResponseWasBadRequest(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private Endpoint Connection %q (Media Service account %q) (ResourceGroup %q): %s", resourceID.Name, resourceID.MediaserviceName, resourceID.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_media_private_endpoint_connection", *existing.ID)
		}
	}

	parameters := media.PrivateEndpointConnection{}

	if _, err := client.CreateOrUpdate(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters); err != nil {
		return fmt.Errorf("Error creating Private Endpoint Connection %q in Media Services Account %q (Resource Group %q): %+v", resourceID.Name, resourceID.MediaserviceName, resourceID.ResourceGroup, err)
	}

	d.SetId(resourceID.ID())

	return resourceMediaPrivateEndpointConnectionRead(d, meta)
}

func resourceMediaPrivateEndpointConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Endpoint %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Private Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)
	d.Set("provisioning_state", resp.ProvisioningState)

	if connectionState := resp.PrivateLinkServiceConnectionState; connectionState != nil {
		d.Set("private_link_connection_state", flattenPrivateLinkConnectionState(connectionState))
	}

	return nil
}

func resourceMediaPrivateEndpointConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name); err != nil {
		return fmt.Errorf("Error deleting Private Endpoint Connection %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}

func flattenPrivateLinkConnectionState(connectionState *media.PrivateLinkServiceConnectionState) []interface{} {
	results := make([]interface{}, 0)

	if connectionState == nil {
		return results
	}

	results = append(results, map[string]interface{}{
		"actions_required": connectionState.ActionsRequired,
		"description":      connectionState.Description,
		"status":           connectionState.Status,
	})

	return results
}
