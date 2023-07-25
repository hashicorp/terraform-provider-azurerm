// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package notificationhub

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			_, err := namespaces.ParseNamespaceID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(namespaces.SkuNameBasic),
					string(namespaces.SkuNameFree),
					string(namespaces.SkuNameStandard),
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
					string(namespaces.NamespaceTypeMessaging),
					string(namespaces.NamespaceTypeNotificationHub),
				}, false),
			},

			"tags": commonschema.Tags(),

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

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_notification_hub_namespace", id.ID())
		}
	}

	namespaceType := namespaces.NamespaceType(d.Get("namespace_type").(string))
	location := location.Normalize(d.Get("location").(string))
	parameters := namespaces.NamespaceCreateOrUpdateParameters{
		Location: location,
		Sku: &namespaces.Sku{
			Name: namespaces.SkuName(d.Get("sku_name").(string)),
		},
		Properties: &namespaces.NamespaceProperties{
			Region:        utils.String(location),
			NamespaceType: &namespaceType,
			Enabled:       utils.Bool(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
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

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		skuName := ""
		if v := model.Sku; v != nil {
			skuName = string(v.Name)
		}
		d.Set("sku_name", skuName)
		if props := model.Properties; props != nil {
			d.Set("enabled", props.Enabled)
			d.Set("servicebus_endpoint", props.ServiceBusEndpoint)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceNotificationHubNamespaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(future.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		if !response.WasNotFound(future.HttpResponse) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}

func notificationHubNamespaceStateRefreshFunc(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		statusCode := "dropped connection"
		if resp.HttpResponse != nil {
			statusCode = strconv.Itoa(resp.HttpResponse.StatusCode)
		}

		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, statusCode, nil
			}

			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return resp, statusCode, nil
	}
}
