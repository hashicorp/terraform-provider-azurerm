// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSignalRServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).SignalRClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resourceType := signalr.SignalRResource{}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model != nil {
		resourceType.Location = existing.Model.Location
	}

	currentSku := ""
	if existing.Model != nil && existing.Model.Sku != nil {
		currentSku = existing.Model.Sku.Name
	}

	if d.HasChange("sku") {
		sku := d.Get("sku").([]any)
		resourceType.Sku = expandSignalRServiceSku(sku)
		currentSku = resourceType.Sku.Name
	}

	if d.HasChangesExcept("sku", "tags") {
		resourceType.Properties = &signalr.SignalRProperties{}

		if d.HasChange("cors") {
			corsRaw := d.Get("cors").([]any)
			resourceType.Properties.Cors = expandSignalRCors(corsRaw)
		}

		if d.HasChange("upstream_endpoint") {
			featuresRaw := d.Get("upstream_endpoint").(*pluginsdk.Set).List()
			resourceType.Properties.Upstream = expandUpstreamSettings(featuresRaw)
		}

		if d.HasChange("serverless_connection_timeout_in_seconds") {
			resourceType.Properties.Serverless = &signalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: pointer.To(int64(d.Get("serverless_connection_timeout_in_seconds").(int))),
			}
		}

		if d.HasChange("identity") {
			identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			resourceType.Identity = identity
		}

		if d.HasChange("public_network_access_enabled") {
			publicNetworkAcc := "Enabled"
			if !d.Get("public_network_access_enabled").(bool) {
				publicNetworkAcc = "Disabled"
			}
			if currentSku == "Free_F1" && publicNetworkAcc == "Disabled" {
				return errors.New("SKU Free_F1 does not support disabling public network access")
			}
			resourceType.Properties.PublicNetworkAccess = pointer.To(publicNetworkAcc)
		}

		if d.HasChange("local_auth_enabled") {
			resourceType.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
		}

		if d.HasChange("aad_auth_enabled") {
			resourceType.Properties.DisableAadAuth = pointer.To(!d.Get("aad_auth_enabled").(bool))
		}

		if d.HasChange("tls_client_cert_enabled") {
			tlsClientCertEnabled := d.Get("tls_client_cert_enabled").(bool)
			resourceType.Properties.Tls = &signalr.SignalRTlsSettings{
				ClientCertEnabled: pointer.To(tlsClientCertEnabled),
			}
			if currentSku == "Free_F1" && tlsClientCertEnabled {
				return errors.New("SKU Free_F1 does not support enabling tls client cert")
			}
		}

		// lintignore:R019 // deliberate subset: only the attributes mapped to SignalR feature flags; the other properties are handled by the surrounding branches
		if d.HasChanges(
			"connectivity_logs_enabled",
			"messaging_logs_enabled",
			"http_request_logs_enabled",
			"live_trace_enabled",
			"service_mode",
		) {
			features := make([]signalr.SignalRFeature, 0)
			if d.HasChanges("connectivity_logs_enabled", "messaging_logs_enabled", "http_request_logs_enabled") {
				connectivityLogsNew := d.Get("connectivity_logs_enabled")
				features = append(features, signalRFeature(signalr.FeatureFlagsEnableConnectivityLogs, strconv.FormatBool(connectivityLogsNew.(bool))))

				messagingLogsNew := d.Get("messaging_logs_enabled")
				features = append(features, signalRFeature(signalr.FeatureFlagsEnableMessagingLogs, strconv.FormatBool(messagingLogsNew.(bool))))

				httpLogsNew := d.Get("http_request_logs_enabled")

				resourceType.Properties.ResourceLogConfiguration = expandSignalRResourceLogConfig(connectivityLogsNew.(bool), messagingLogsNew.(bool), httpLogsNew.(bool))
			}

			if d.HasChange("live_trace_enabled") {
				liveTraceEnabled := false
				if v, ok := d.GetOk("live_trace_enabled"); ok {
					liveTraceEnabled = v.(bool)
				}
				features = append(features, signalRFeature("EnableLiveTrace", strconv.FormatBool(liveTraceEnabled)))
			}

			if d.HasChange("service_mode") {
				serviceMode := "Default"
				if v, ok := d.GetOk("service_mode"); ok {
					serviceMode = v.(string)
				}
				features = append(features, signalRFeature(signalr.FeatureFlagsServiceMode, serviceMode))
			}

			resourceType.Properties.Features = &features
		}

		if d.HasChange("live_trace") {
			resourceType.Properties.LiveTraceConfiguration = expandSignalRLiveTraceConfig(d.Get("live_trace").([]any))
		}
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]any)
		resourceType.Tags = tags.Expand(tagsRaw)
	}

	if err := client.UpdateThenPoll(ctx, *id, resourceType); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			string(signalr.ProvisioningStateUpdating),
			string(signalr.ProvisioningStateCreating),
			string(signalr.ProvisioningStateMoving),
			string(signalr.ProvisioningStateRunning),
		},
		Target:                    []string{string(signalr.ProvisioningStateSucceeded)},
		Refresh:                   signalrServiceProvisioningStateRefreshFunc(ctx, client, *id),
		Timeout:                   time.Until(deadline),
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 20,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	return resourceArmSignalRServiceRead(d, meta)
}
