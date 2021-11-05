package notificationhub

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/migration"

	"github.com/hashicorp/terraform-provider-azurerm/internal/location"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var notificationHubNamespaceResourceName = "azurerm_notification_hub_namespace"

func resourceNotificationHubNamespace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNotificationHubNamespaceCreateUpdate,
		Read:   resourceNotificationHubNamespaceRead,
		Update: resourceNotificationHubNamespaceCreateUpdate,
		Delete: resourceNotificationHubNamespaceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NamespaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NotificationHubNamespaceResourceV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(notificationhubs.Basic),
					string(notificationhubs.Free),
					string(notificationhubs.Standard),
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"namespace_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(notificationhubs.Messaging),
					string(notificationhubs.NotificationHub),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"tags": tags.Schema(),

			"servicebus_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNotificationHubNamespaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_notification_hub_namespace", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	parameters := notificationhubs.NamespaceCreateOrUpdateParameters{
		Location: utils.String(location),
		Sku: &notificationhubs.Sku{
			Name: notificationhubs.SkuName(d.Get("sku_name").(string)),
		},
		NamespaceProperties: &notificationhubs.NamespaceProperties{
			Region:        utils.String(location),
			NamespaceType: notificationhubs.NamespaceType(d.Get("namespace_type").(string)),
			Enabled:       utils.Bool(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be created..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   notificationHubNamespaceStateRefreshFunc(ctx, client, id),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %ss to finish replicating: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNotificationHubNamespaceRead(d, meta)
}

func resourceNotificationHubNamespaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	skuName := ""
	if resp.Sku != nil {
		skuName = string(resp.Sku.Name)
	}
	d.Set("sku_name", skuName)

	if props := resp.NamespaceProperties; props != nil {
		d.Set("enabled", props.Enabled)
		d.Set("namespace_type", props.NamespaceType)
		d.Set("servicebus_endpoint", props.ServiceBusEndpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNotificationHubNamespaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Notification Hub Namespace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	// the future returned from the Delete method is broken 50% of the time - let's poll ourselves for now
	// Related Bug: https://github.com/Azure/azure-sdk-for-go/issues/2254
	log.Printf("[DEBUG] Waiting for Notification Hub Namespace %q (Resource Group %q) to be deleted", id.Name, id.ResourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: notificationHubNamespaceDeleteStateRefreshFunc(ctx, client, *id),
		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Notification Hub %q (Resource Group %q) to be deleted: %s", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func notificationHubNamespaceStateRefreshFunc(ctx context.Context, client *notificationhubs.NamespacesClient, id parse.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "404", nil
			}

			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func notificationHubNamespaceDeleteStateRefreshFunc(ctx context.Context, client *notificationhubs.NamespacesClient, id parse.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		// Note: this exists as the Delete API only seems to work some of the time
		// in this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attepmpt
		// Upstream Bug: https://github.com/Azure/azure-sdk-for-go/issues/2254
		if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
			log.Printf("re-issuing deletion request for %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
