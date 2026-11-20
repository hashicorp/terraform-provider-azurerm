// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSignalRServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := location.Normalize(d.Get("location").(string))

	id := signalr.NewSignalRID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_signalr_service", id.ID())
		}
	}

	sku := d.Get("sku").([]any)
	connectivityLogsEnabled := false
	if v, ok := d.GetOk("connectivity_logs_enabled"); ok {
		connectivityLogsEnabled = v.(bool)
	}
	messagingLogsEnabled := false
	if v, ok := d.GetOk("messaging_logs_enabled"); ok {
		messagingLogsEnabled = v.(bool)
	}

	httpLogsEnabled := false
	if v, ok := d.GetOk("http_request_logs_enabled"); ok {
		httpLogsEnabled = v.(bool)
	}
	liveTraceEnabled := false
	if v, ok := d.GetOk("live_trace_enabled"); ok {
		liveTraceEnabled = v.(bool)
	}
	serviceMode := "Default"
	if v, ok := d.GetOk("service_mode"); ok {
		serviceMode = v.(string)
	}

	cors := d.Get("cors").([]any)
	upstreamSettings := d.Get("upstream_endpoint").(*pluginsdk.Set).List()

	expandedFeatures := make([]signalr.SignalRFeature, 0)
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableConnectivityLogs, strconv.FormatBool(connectivityLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableMessagingLogs, strconv.FormatBool(messagingLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature("EnableLiveTrace", strconv.FormatBool(liveTraceEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsServiceMode, serviceMode))

	// Upstream configurations are only allowed when the SignalR service is in `Serverless` mode
	if len(upstreamSettings) > 0 && !signalRIsInServerlessMode(&expandedFeatures) {
		return errors.New("upstream configurations are only allowed when the SignalR Service is in `Serverless` mode")
	}

	publicNetworkAcc := "Enabled"
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAcc = "Disabled"
	}

	tlsClientCertEnabled := d.Get("tls_client_cert_enabled").(bool)

	if expandSignalRServiceSku(sku).Name == "Free_F1" {
		if publicNetworkAcc == "Disabled" {
			return fmt.Errorf("SKU Free_F1 does not support disabling public network access")
		}
		if tlsClientCertEnabled {
			return fmt.Errorf("SKU Free_F1 does not support enabling tls client cert")
		}
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	resourceLogsData := expandSignalRResourceLogConfig(connectivityLogsEnabled, messagingLogsEnabled, httpLogsEnabled)
	resourceType := signalr.SignalRResource{
		Location: location,
		Identity: identity,
		Properties: &signalr.SignalRProperties{
			Cors:                     expandSignalRCors(cors),
			Features:                 &expandedFeatures,
			Upstream:                 expandUpstreamSettings(upstreamSettings),
			LiveTraceConfiguration:   expandSignalRLiveTraceConfig(d.Get("live_trace").([]any)),
			ResourceLogConfiguration: resourceLogsData,
			PublicNetworkAccess:      pointer.To(publicNetworkAcc),
			DisableAadAuth:           pointer.To(!d.Get("aad_auth_enabled").(bool)),
			DisableLocalAuth:         pointer.To(!d.Get("local_auth_enabled").(bool)),
			Tls: &signalr.SignalRTlsSettings{
				ClientCertEnabled: pointer.To(tlsClientCertEnabled),
			},
			Serverless: &signalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: pointer.To(int64(d.Get("serverless_connection_timeout_in_seconds").(int))),
			},
		},
		Sku:  expandSignalRServiceSku(sku),
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, resourceType, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceArmSignalRServiceRead(d, meta)
}
