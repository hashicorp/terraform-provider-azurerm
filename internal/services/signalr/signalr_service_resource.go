// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/migration"
	signalrValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmSignalRService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmSignalRServiceCreate,
		Read:   resourceArmSignalRServiceRead,
		Update: resourceArmSignalRServiceUpdate,
		Delete: resourceArmSignalRServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceV0ToV1{},
		}),
		SchemaVersion: 1,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSignalRID(id)
			return err
		}),

		Schema: resourceArmSignalRServiceSchema(),
	}
}

func resourceArmSignalRServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))

	id := signalr.NewSignalRID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_signalr_service", id.ID())
	}

	sku := d.Get("sku").([]interface{})
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

	cors := d.Get("cors").([]interface{})
	upstreamSettings := d.Get("upstream_endpoint").(*pluginsdk.Set).List()

	expandedFeatures := make([]signalr.SignalRFeature, 0)
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableConnectivityLogs, strconv.FormatBool(connectivityLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableMessagingLogs, strconv.FormatBool(messagingLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature("EnableLiveTrace", strconv.FormatBool(liveTraceEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsServiceMode, serviceMode))

	// Upstream configurations are only allowed when the SignalR service is in `Serverless` mode
	if len(upstreamSettings) > 0 && !signalRIsInServerlessMode(&expandedFeatures) {
		return fmt.Errorf("Upstream configurations are only allowed when the SignalR Service is in `Serverless` mode")
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

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	resourceLogsData := expandSignalRResourceLogConfig(connectivityLogsEnabled, messagingLogsEnabled, httpLogsEnabled)
	resourceType := signalr.SignalRResource{
		Location: utils.String(location),
		Identity: identity,
		Properties: &signalr.SignalRProperties{
			Cors:                     expandSignalRCors(cors),
			Features:                 &expandedFeatures,
			Upstream:                 expandUpstreamSettings(upstreamSettings),
			LiveTraceConfiguration:   expandSignalRLiveTraceConfig(d.Get("live_trace").([]interface{})),
			ResourceLogConfiguration: resourceLogsData,
			PublicNetworkAccess:      utils.String(publicNetworkAcc),
			DisableAadAuth:           utils.Bool(!d.Get("aad_auth_enabled").(bool)),
			DisableLocalAuth:         utils.Bool(!d.Get("local_auth_enabled").(bool)),
			Tls: &signalr.SignalRTlsSettings{
				ClientCertEnabled: utils.Bool(tlsClientCertEnabled),
			},
			Serverless: &signalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: utils.Int64(int64(d.Get("serverless_connection_timeout_in_seconds").(int))),
			},
		},
		Sku:  expandSignalRServiceSku(sku),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, resourceType); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmSignalRServiceUpdate(d, meta)
}

func resourceArmSignalRServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
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

	keys, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("name", id.SignalRName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if err = d.Set("sku", flattenSignalRServiceSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("hostname", props.HostName)
			d.Set("ip_address", props.ExternalIP)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)

			connectivityLogsEnabled := false
			messagingLogsEnabled := false
			httpLogsEnabled := false
			liveTraceEnabled := false
			serviceMode := "Default"
			for _, feature := range *props.Features {
				if feature.Flag == "EnableLiveTrace" {
					liveTraceEnabled = strings.EqualFold(feature.Value, "True")
				}
				if feature.Flag == signalr.FeatureFlagsServiceMode {
					serviceMode = feature.Value
				}
			}

			d.Set("live_trace_enabled", liveTraceEnabled)
			d.Set("service_mode", serviceMode)

			aadAuthEnabled := true
			if props.DisableAadAuth != nil {
				aadAuthEnabled = !(*props.DisableAadAuth)
			}
			d.Set("aad_auth_enabled", aadAuthEnabled)

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !(*props.DisableLocalAuth)
			}
			d.Set("local_auth_enabled", localAuthEnabled)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = strings.EqualFold(*props.PublicNetworkAccess, "Enabled")
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			tlsClientCertEnabled := false
			if props.Tls != nil && props.Tls.ClientCertEnabled != nil {
				tlsClientCertEnabled = *props.Tls.ClientCertEnabled
			}
			d.Set("tls_client_cert_enabled", tlsClientCertEnabled)

			if props.Serverless != nil && props.Serverless.ConnectionTimeoutInSeconds != nil {
				d.Set("serverless_connection_timeout_in_seconds", int(*props.Serverless.ConnectionTimeoutInSeconds))
			}

			if err := d.Set("cors", flattenSignalRCors(props.Cors)); err != nil {
				return fmt.Errorf("setting `cors`: %+v", err)
			}

			if err := d.Set("upstream_endpoint", flattenUpstreamSettings(props.Upstream)); err != nil {
				return fmt.Errorf("setting `upstream_endpoint`: %+v", err)
			}

			if err := d.Set("live_trace", flattenSignalRLiveTraceConfig(props.LiveTraceConfiguration)); err != nil {
				return fmt.Errorf("setting `live_trace`:%+v", err)
			}

			if props.ResourceLogConfiguration != nil && props.ResourceLogConfiguration.Categories != nil {
				for _, item := range *props.ResourceLogConfiguration.Categories {
					name := ""
					if item.Name != nil {
						name = *item.Name
					}

					var cateEnabled string
					if item.Enabled != nil {
						cateEnabled = *item.Enabled
					}

					switch name {
					case "MessagingLogs":
						messagingLogsEnabled = strings.EqualFold(cateEnabled, "true")
					case "ConnectivityLogs":
						connectivityLogsEnabled = strings.EqualFold(cateEnabled, "true")
					case "HttpRequestLogs":
						httpLogsEnabled = strings.EqualFold(cateEnabled, "true")
					default:
						continue
					}
				}
				d.Set("connectivity_logs_enabled", connectivityLogsEnabled)
				d.Set("messaging_logs_enabled", messagingLogsEnabled)
				d.Set("http_request_logs_enabled", httpLogsEnabled)
			}
			identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_access_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
	}

	return nil
}

func resourceArmSignalRServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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

	currentSku := ""
	if existing.Model != nil && existing.Model.Sku != nil {
		currentSku = existing.Model.Sku.Name
	}

	if d.HasChange("sku") {
		sku := d.Get("sku").([]interface{})
		resourceType.Sku = expandSignalRServiceSku(sku)
		currentSku = resourceType.Sku.Name
	}

	if d.HasChanges("cors", "upstream_endpoint", "serverless_connection_timeout_in_seconds", "identity",
		"public_network_access_enabled", "local_auth_enabled", "aad_auth_enabled", "tls_client_cert_enabled",
		"features", "connectivity_logs_enabled", "messaging_logs_enabled", "http_request_logs_enabled", "service_mode", "live_trace_enabled", "live_trace") {
		resourceType.Properties = &signalr.SignalRProperties{}

		if d.HasChange("cors") {
			corsRaw := d.Get("cors").([]interface{})
			resourceType.Properties.Cors = expandSignalRCors(corsRaw)
		}

		if d.HasChange("upstream_endpoint") {
			featuresRaw := d.Get("upstream_endpoint").(*pluginsdk.Set).List()
			resourceType.Properties.Upstream = expandUpstreamSettings(featuresRaw)
		}

		if d.HasChange("serverless_connection_timeout_in_seconds") {
			resourceType.Properties.Serverless = &signalr.ServerlessSettings{
				ConnectionTimeoutInSeconds: utils.Int64(int64(d.Get("serverless_connection_timeout_in_seconds").(int))),
			}
		}

		if d.HasChange("identity") {
			identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
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
				return fmt.Errorf("SKU Free_F1 does not support disabling public network access")
			}
			resourceType.Properties.PublicNetworkAccess = utils.String(publicNetworkAcc)
		}

		if d.HasChange("local_auth_enabled") {
			resourceType.Properties.DisableLocalAuth = utils.Bool(!d.Get("local_auth_enabled").(bool))
		}

		if d.HasChange("aad_auth_enabled") {
			resourceType.Properties.DisableAadAuth = utils.Bool(!d.Get("aad_auth_enabled").(bool))
		}

		if d.HasChange("tls_client_cert_enabled") {
			tlsClientCertEnabled := d.Get("tls_client_cert_enabled").(bool)
			resourceType.Properties.Tls = &signalr.SignalRTlsSettings{
				ClientCertEnabled: utils.Bool(tlsClientCertEnabled),
			}
			if currentSku == "Free_F1" && tlsClientCertEnabled {
				return fmt.Errorf("SKU Free_F1 does not support enabling tls client cert")
			}
		}

		if d.HasChanges("connectivity_logs_enabled", "messaging_logs_enabled", "http_request_logs_enabled", "live_trace_enabled", "service_mode") {
			features := make([]signalr.SignalRFeature, 0)
			if d.HasChange("connectivity_logs_enabled") || d.HasChange("messaging_logs_enabled") || d.HasChange("http_request_logs_enabled") {
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
			resourceType.Properties.LiveTraceConfiguration = expandSignalRLiveTraceConfig(d.Get("live_trace").([]interface{}))
		}
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		resourceType.Tags = tags.Expand(tagsRaw)
	}

	if err := client.UpdateThenPoll(ctx, *id, resourceType); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
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

func resourceArmSignalRServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// @tombuildsstuff: we can't use DeleteThenPoll here since the API returns a 404 on the Future in time
	future, err := client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		if r := future.Poller.LatestResponse(); r == nil || !response.WasNotFound(r.Response) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}

func signalRIsInServerlessMode(features *[]signalr.SignalRFeature) bool {
	if features == nil {
		return false
	}

	for _, feature := range *features {
		if feature.Flag == signalr.FeatureFlagsServiceMode {
			return strings.EqualFold(feature.Value, "Serverless")
		}
	}

	return false
}

func signalRFeature(featureFlag signalr.FeatureFlags, value string) signalr.SignalRFeature {
	return signalr.SignalRFeature{
		Flag:  featureFlag,
		Value: value,
	}
}

func expandUpstreamSettings(input []interface{}) *signalr.ServerlessUpstreamSettings {
	upstreamTemplates := make([]signalr.UpstreamTemplate, 0)

	for _, upstreamSetting := range input {
		setting := upstreamSetting.(map[string]interface{})
		authTypeNone := signalr.UpstreamAuthTypeNone
		authTypeManagedIdentity := signalr.UpstreamAuthTypeManagedIdentity
		auth := signalr.UpstreamAuthSettings{
			Type: &authTypeNone,
		}
		upstreamTemplate := signalr.UpstreamTemplate{
			HubPattern:      utils.String(strings.Join(*utils.ExpandStringSlice(setting["hub_pattern"].([]interface{})), ",")),
			EventPattern:    utils.String(strings.Join(*utils.ExpandStringSlice(setting["event_pattern"].([]interface{})), ",")),
			CategoryPattern: utils.String(strings.Join(*utils.ExpandStringSlice(setting["category_pattern"].([]interface{})), ",")),
			UrlTemplate:     setting["url_template"].(string),
			Auth:            &auth,
		}

		if setting["user_assigned_identity_id"].(string) != "" {
			upstreamTemplate.Auth = &signalr.UpstreamAuthSettings{
				Type: &authTypeManagedIdentity,
				ManagedIdentity: &signalr.ManagedIdentitySettings{
					Resource: utils.String(setting["user_assigned_identity_id"].(string)),
				},
			}
		}

		upstreamTemplates = append(upstreamTemplates, upstreamTemplate)
	}

	return &signalr.ServerlessUpstreamSettings{
		Templates: &upstreamTemplates,
	}
}

func flattenUpstreamSettings(upstreamSettings *signalr.ServerlessUpstreamSettings) []interface{} {
	result := make([]interface{}, 0)
	if upstreamSettings == nil || upstreamSettings.Templates == nil {
		return result
	}

	for _, settings := range *upstreamSettings.Templates {
		categoryPattern := make([]interface{}, 0)
		if settings.CategoryPattern != nil {
			categoryPatterns := strings.Split(*settings.CategoryPattern, ",")
			categoryPattern = utils.FlattenStringSlice(&categoryPatterns)
		}

		eventPattern := make([]interface{}, 0)
		if settings.EventPattern != nil {
			eventPatterns := strings.Split(*settings.EventPattern, ",")
			eventPattern = utils.FlattenStringSlice(&eventPatterns)
		}

		hubPattern := make([]interface{}, 0)
		if settings.HubPattern != nil {
			hubPatterns := strings.Split(*settings.HubPattern, ",")
			hubPattern = utils.FlattenStringSlice(&hubPatterns)
		}

		var managedIdentityId string
		if upstreamAuth := settings.Auth; upstreamAuth != nil && upstreamAuth.Type != nil && *upstreamAuth.Type != signalr.UpstreamAuthTypeNone {
			if upstreamAuth.ManagedIdentity != nil && upstreamAuth.ManagedIdentity.Resource != nil {
				managedIdentityId = *upstreamAuth.ManagedIdentity.Resource
			}
		}

		result = append(result, map[string]interface{}{
			"url_template":              settings.UrlTemplate,
			"hub_pattern":               hubPattern,
			"event_pattern":             eventPattern,
			"category_pattern":          categoryPattern,
			"user_assigned_identity_id": managedIdentityId,
		})
	}
	return result
}

func expandSignalRCors(input []interface{}) *signalr.SignalRCorsSettings {
	corsSettings := signalr.SignalRCorsSettings{}

	if len(input) == 0 || input[0] == nil {
		return &corsSettings
	}

	setting := input[0].(map[string]interface{})
	origins := setting["allowed_origins"].(*pluginsdk.Set).List()

	allowedOrigins := make([]string, 0)
	for _, param := range origins {
		allowedOrigins = append(allowedOrigins, param.(string))
	}

	corsSettings.AllowedOrigins = &allowedOrigins

	return &corsSettings
}

func flattenSignalRCors(input *signalr.SignalRCorsSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	allowedOrigins := make([]interface{}, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	return append(results, result)
}

func expandSignalRServiceSku(input []interface{}) *signalr.ResourceSku {
	v := input[0].(map[string]interface{})
	return &signalr.ResourceSku{
		Name:     v["name"].(string),
		Capacity: utils.Int64(int64(v["capacity"].(int))),
	}
}

func flattenSignalRServiceSku(input *signalr.ResourceSku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	capacity := 0
	if input.Capacity != nil {
		capacity = int(*input.Capacity)
	}

	return []interface{}{
		map[string]interface{}{
			"capacity": capacity,
			"name":     input.Name,
		},
	}
}

func expandSignalRLiveTraceConfig(input []interface{}) *signalr.LiveTraceConfiguration {
	resourceCategories := make([]signalr.LiveTraceCategory, 0)
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	enabled := "false"
	if v["enabled"].(bool) {
		enabled = "true"
	}

	messageLogEnabled := "false"
	if v["messaging_logs_enabled"].(bool) {
		messageLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    utils.String("MessagingLogs"),
		Enabled: utils.String(messageLogEnabled),
	})

	connectivityLogEnabled := "false"
	if v["connectivity_logs_enabled"].(bool) {
		connectivityLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    utils.String("ConnectivityLogs"),
		Enabled: utils.String(connectivityLogEnabled),
	})

	httpLogEnabled := "false"
	if v["http_request_logs_enabled"].(bool) {
		httpLogEnabled = "true"
	}
	resourceCategories = append(resourceCategories, signalr.LiveTraceCategory{
		Name:    utils.String("HttpRequestLogs"),
		Enabled: utils.String(httpLogEnabled),
	})

	return &signalr.LiveTraceConfiguration{
		Enabled:    &enabled,
		Categories: &resourceCategories,
	}
}

func flattenSignalRLiveTraceConfig(input *signalr.LiveTraceConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	var enabled bool
	if input.Enabled != nil {
		enabled = strings.EqualFold(*input.Enabled, "true")
	}

	var (
		messagingLogEnabled    bool
		connectivityLogEnabled bool
		httpLogsEnabled        bool
	)

	if input.Categories != nil {
		for _, item := range *input.Categories {
			name := ""
			if item.Name != nil {
				name = *item.Name
			}

			var cateEnabled string
			if item.Enabled != nil {
				cateEnabled = *item.Enabled
			}

			switch name {
			case "MessagingLogs":
				messagingLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "ConnectivityLogs":
				connectivityLogEnabled = strings.EqualFold(cateEnabled, "true")
			case "HttpRequestLogs":
				httpLogsEnabled = strings.EqualFold(cateEnabled, "true")
			default:
				continue
			}
		}
	}
	return []interface{}{map[string]interface{}{
		"enabled":                   enabled,
		"messaging_logs_enabled":    messagingLogEnabled,
		"connectivity_logs_enabled": connectivityLogEnabled,
		"http_request_logs_enabled": httpLogsEnabled,
	}}
}

func expandSignalRResourceLogConfig(connectivityLogEnabled bool, messagingLogEnabled bool, httpLogEnabled bool) *signalr.ResourceLogConfiguration {
	resourceLogCategories := make([]signalr.ResourceLogCategory, 0)

	messagingLog := "false"
	if messagingLogEnabled {
		messagingLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    utils.String("MessagingLogs"),
		Enabled: utils.String(messagingLog),
	})

	connectivityLog := "false"
	if connectivityLogEnabled {
		connectivityLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    utils.String("ConnectivityLogs"),
		Enabled: utils.String(connectivityLog),
	})

	httpLog := "false"
	if httpLogEnabled {
		httpLog = "true"
	}
	resourceLogCategories = append(resourceLogCategories, signalr.ResourceLogCategory{
		Name:    utils.String("HttpRequestLogs"),
		Enabled: utils.String(httpLog),
	})

	return &signalr.ResourceLogConfiguration{
		Categories: &resourceLogCategories,
	}
}

func resourceArmSignalRServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Free_F1",
							"Standard_S1",
							"Premium_P1",
						}, false),
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}),
					},
				},
			},
		},

		"connectivity_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"messaging_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"http_request_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"live_trace_enabled": {
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Default:    false,
			Deprecated: "`live_trace_enabled` has been deprecated in favor of `live_trace` and will be removed in 4.0.",
		},

		"live_trace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"connectivity_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},

					"messaging_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},

					"http_request_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"aad_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tls_client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"serverless_connection_timeout_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  30,
		},

		"service_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Default",
			ValidateFunc: validation.StringInSlice([]string{
				"Serverless",
				"Classic",
				"Default",
			}, false),
		},

		"upstream_endpoint": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"event_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"hub_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"url_template": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: signalrValidate.UrlTemplate,
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"server_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": commonschema.Tags(),
	}
}

func signalrServiceProvisioningStateRefreshFunc(ctx context.Context, client *signalr.SignalRClient, id signalr.SignalRId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for the provisioning state of %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if model.Properties != nil && model.Properties.ProvisioningState != nil {
				return res, string(*model.Properties.ProvisioningState), nil
			}
		}

		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
