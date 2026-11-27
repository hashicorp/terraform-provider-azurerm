// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubSubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webpubsub.NewWebPubSubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	liveTraceConfig := d.Get("live_trace").([]any)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(webPubSubResourceType, id.ID())
			}
		}
	}

	publicNetworkAcc := "Enabled"
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAcc = "Disabled"
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := webpubsub.WebPubSubResource{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: identity,
		Properties: &webpubsub.WebPubSubProperties{
			LiveTraceConfiguration: expandLiveTraceConfig(liveTraceConfig),
			PublicNetworkAccess:    pointer.To(publicNetworkAcc),
			DisableAadAuth:         pointer.To(!d.Get("aad_auth_enabled").(bool)),
			DisableLocalAuth:       pointer.To(!d.Get("local_auth_enabled").(bool)),
			Tls: &webpubsub.WebPubSubTlsSettings{
				ClientCertEnabled: pointer.To(d.Get("tls_client_cert_enabled").(bool)),
			},
		},
		Sku: &webpubsub.ResourceSku{
			Name:     d.Get("sku").(string),
			Capacity: pointer.To(int64(d.Get("capacity").(int))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceWebPubSubRead(d, meta)
}

func expandLiveTraceConfig(input []any) *webpubsub.LiveTraceConfiguration {
	resourceCategories := make([]webpubsub.LiveTraceCategory, 0)
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	enabled := "false"
	if v["enabled"].(bool) {
		enabled = "true"
	}

	messageLogEnabled := "false"
	if v["messaging_logs_enabled"].(bool) {
		messageLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("MessagingLogs"),
		Enabled: pointer.To(messageLogEnabled),
	})

	connectivityLogEnabled := "false"
	if v["connectivity_logs_enabled"].(bool) {
		connectivityLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("ConnectivityLogs"),
		Enabled: pointer.To(connectivityLogEnabled),
	})

	httpLogEnabled := "false"
	if v["http_request_logs_enabled"].(bool) {
		httpLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, webpubsub.LiveTraceCategory{
		Name:    pointer.To("HttpRequestLogs"),
		Enabled: pointer.To(httpLogEnabled),
	})

	return &webpubsub.LiveTraceConfiguration{
		Enabled:    &enabled,
		Categories: &resourceCategories,
	}
}
