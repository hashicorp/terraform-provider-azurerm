// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	frontDoorValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontDoor() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCreate,
		Read:   resourceFrontDoorRead,
		Update: resourceFrontDoorUpdate,
		Delete: resourceFrontDoorDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := frontdoors.ParseFrontDoorID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.FrontDoorUpgradeV0ToV1{},
			1: migration.FrontDoorUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: resourceFrontDoorSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(frontDoorCustomizeDiff),
	}
}

func resourceFrontDoorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := frontdoors.NewFrontDoorID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_frontdoor", id.ID())
	}

	var backendCertNameCheck bool
	var backendPoolsSendReceiveTimeoutSeconds int64
	if bps, ok := d.Get("backend_pool_settings").([]interface{}); ok && len(bps) > 0 {
		bpsMap := bps[0].(map[string]interface{})
		if v, ok := bpsMap["enforce_backend_pools_certificate_name_check"].(bool); ok {
			backendCertNameCheck = v
		}
		if v, ok := bpsMap["backend_pools_send_receive_timeout_seconds"].(int); ok {
			backendPoolsSendReceiveTimeoutSeconds = int64(v)
		}
	}

	friendlyName := d.Get("friendly_name").(string)
	routingRules := d.Get("routing_rule").([]interface{})
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})
	backendPools := d.Get("backend_pool").([]interface{})
	frontendEndpoints := d.Get("frontend_endpoint").([]interface{})

	enabledState := expandFrontDoorEnabledState(d.Get("load_balancer_enabled").(bool))
	t := d.Get("tags").(map[string]interface{})

	frontDoorParameters := frontdoors.FrontDoor{
		Location: utils.String("Global"),
		Properties: &frontdoors.FrontDoorProperties{
			FriendlyName:          utils.String(friendlyName),
			RoutingRules:          expandFrontDoorRoutingRule(routingRules, id, nil),
			BackendPools:          expandFrontDoorBackendPools(backendPools, id),
			BackendPoolsSettings:  expandFrontDoorBackendPoolsSettings(backendCertNameCheck, backendPoolsSendReceiveTimeoutSeconds),
			FrontendEndpoints:     expandFrontDoorFrontendEndpoint(frontendEndpoints, id),
			HealthProbeSettings:   expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, id),
			LoadBalancingSettings: expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, id),
			EnabledState:          &enabledState,
		},
		Tags: tags.Expand(t),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, frontDoorParameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))

	d.SetId(id.ID())
	return resourceFrontDoorRead(d, meta)
}

func resourceFrontDoorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := frontdoors.NewFrontDoorID(subscriptionId, resourceGroup, name)

	// remove in 3.0
	// due to a change in the RP, if a Frontdoor exists in a location other than 'Global' it may continue to
	// exist in that location, if this is a brand new Frontdoor it must be created in the 'Global' location
	var location string

	exists, err := client.Get(ctx, id)
	if err != nil || exists.Model == nil {
		return fmt.Errorf("locating %s: %+v", id, err)
	} else {
		location = azure.NormalizeLocation(*exists.Model.Location)
	}

	cfgLocation, hasLocation := d.GetOk("location")
	if hasLocation {
		if location != azure.NormalizeLocation(cfgLocation) {
			return fmt.Errorf("the Front Door %q (Resource Group %q) already exists in %q and cannot be moved to the %q location", name, resourceGroup, location, cfgLocation)
		}
	}

	existingModel := *exists.Model

	if d.HasChange("friendly_name") {
		existingModel.Properties.FriendlyName = utils.String(d.Get("friendly_name").(string))
	}

	routingRules := d.Get("routing_rule").([]interface{})
	if d.HasChange("routing_rule") {
		rulesEngines := make(map[string]*frontdoors.SubResource)
		if existingModel.Properties != nil && existingModel.Properties.RoutingRules != nil {
			for _, rule := range *existingModel.Properties.RoutingRules {
				if rule.Properties != nil && rule.Properties.RulesEngine != nil {
					rulesEngines[*rule.Name] = rule.Properties.RulesEngine
				}
			}
		}
		existingModel.Properties.RoutingRules = expandFrontDoorRoutingRule(routingRules, id, &rulesEngines)
	}

	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	if d.HasChange("backend_pool_load_balancing") {
		existingModel.Properties.LoadBalancingSettings = expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, id)
	}

	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})
	if d.HasChange("backend_pool_health_probe") {
		existingModel.Properties.HealthProbeSettings = expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, id)
	}

	backendPools := d.Get("backend_pool").([]interface{})
	if d.HasChange("backend_pool") {
		existingModel.Properties.BackendPools = expandFrontDoorBackendPools(backendPools, id)
	}

	frontendEndpoints := d.Get("frontend_endpoint").([]interface{})
	if d.HasChange("frontend_endpoint") {
		existingModel.Properties.FrontendEndpoints = expandFrontDoorFrontendEndpoint(frontendEndpoints, id)
	}

	if d.HasChange("backend_pool_settings") {
		var backendCertNameCheck bool
		var backendPoolsSendReceiveTimeoutSeconds int64
		if bps, ok := d.Get("backend_pool_settings").([]interface{}); ok && len(bps) > 0 {
			bpsMap := bps[0].(map[string]interface{})
			if v, ok := bpsMap["enforce_backend_pools_certificate_name_check"].(bool); ok {
				backendCertNameCheck = v
			}
			if v, ok := bpsMap["backend_pools_send_receive_timeout_seconds"].(int); ok {
				backendPoolsSendReceiveTimeoutSeconds = int64(v)
			}
			existingModel.Properties.BackendPoolsSettings = expandFrontDoorBackendPoolsSettings(backendCertNameCheck, backendPoolsSendReceiveTimeoutSeconds)
		}
	}

	if d.HasChange("load_balancer_enabled") {
		enabledState := expandFrontDoorEnabledState(d.Get("load_balancer_enabled").(bool))
		existingModel.Properties.EnabledState = &enabledState
	}

	if d.HasChanges("tags") {
		existingModel.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	// If the explicitResourceOrder is empty and it's not a new resource set the mapping table to the state file and return an error.
	// If the explicitResourceOrder is empty and it is a new resource it will run the CreateOrUpdate as expected
	// If the explicitResourceOrder is NOT empty and it is NOT a new resource it will run the CreateOrUpdate as expected
	explicitResourceOrder := d.Get("explicit_resource_order").([]interface{})
	if len(explicitResourceOrder) == 0 {
		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, existingModel); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))
	}

	d.SetId(id.ID())
	return resourceFrontDoorRead(d, meta)
}

func resourceFrontDoorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontDoorIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Front Door %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.FrontDoorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			explicitResourceOrder := d.Get("explicit_resource_order").([]interface{})
			flattenedBackendPools, err := flattenFrontDoorBackendPools(props.BackendPools, *id, explicitResourceOrder)
			if err != nil {
				return fmt.Errorf("flattening `backend_pool`: %+v", err)
			}
			if err := d.Set("backend_pool", flattenedBackendPools); err != nil {
				return fmt.Errorf("setting `backend_pool`: %+v", err)
			}

			backendPoolSettings := flattenFrontDoorBackendPoolsSettings(props.BackendPoolsSettings)
			out := map[string]interface{}{
				"enforce_backend_pools_certificate_name_check": backendPoolSettings.enforceBackendPoolsCertificateNameCheck,
				"backend_pools_send_receive_timeout_seconds":   backendPoolSettings.backendPoolsSendReceiveTimeoutSeconds,
			}
			d.Set("backend_pool_settings", []interface{}{out})

			d.Set("cname", props.Cname)
			d.Set("header_frontdoor_id", props.FrontdoorId)
			if props.EnabledState != nil {
				d.Set("load_balancer_enabled", *props.EnabledState == frontdoors.FrontDoorEnabledStateEnabled)
			}
			d.Set("friendly_name", props.FriendlyName)

			// Need to call frontEndEndpointClient here to get the frontEndEndpoint information from that client
			// because the information is hidden from the main frontDoorClient "by design"...
			frontEndEndpointsClient := meta.(*clients.Client).Frontdoor.FrontDoorsClient
			frontEndEndpointInfo, err := retrieveFrontEndEndpointInformation(ctx, frontEndEndpointsClient, *id, props.FrontendEndpoints)
			if err != nil {
				return fmt.Errorf("retrieving FrontEnd Endpoint Information: %+v", err)
			}

			// Force the returned flattenFrontEndEndpoints into the order defined in the explicit_resource_order mapping table
			frontDoorFrontendEndpoints, err := flattenFrontEndEndpoints(frontEndEndpointInfo, *id, explicitResourceOrder)
			if err != nil {
				return fmt.Errorf("flattening `frontend_endpoint`: %+v", err)
			}
			if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoints); err != nil {
				return fmt.Errorf("setting `frontend_endpoint`: %+v", err)
			}

			// Force the returned flattenFrontDoorHealthProbeSettingsModel into the order defined in the explicit_resource_order mapping table
			if err := d.Set("backend_pool_health_probe", flattenFrontDoorHealthProbeSettingsModel(props.HealthProbeSettings, *id, explicitResourceOrder)); err != nil {
				return fmt.Errorf("setting `backend_pool_health_probe`: %+v", err)
			}

			// Force the returned flattenFrontDoorLoadBalancingSettingsModel into the order defined in the explicit_resource_order mapping table
			if err := d.Set("backend_pool_load_balancing", flattenFrontDoorLoadBalancingSettingsModel(props.LoadBalancingSettings, *id, explicitResourceOrder)); err != nil {
				return fmt.Errorf("setting `backend_pool_load_balancing`: %+v", err)
			}

			var flattenedRoutingRules *[]interface{}
			// Force the returned flattenedRoutingRules into the order defined in the explicit_resource_order mapping table
			flattenedRoutingRules, err = flattenFrontDoorRoutingRule(props.RoutingRules, d.Get("routing_rule"), *id, explicitResourceOrder)
			if err != nil {
				return fmt.Errorf("flattening `routing_rules`: %+v", err)
			}
			if err := d.Set("routing_rule", flattenedRoutingRules); err != nil {
				return fmt.Errorf("setting `routing_rules`: %+v", err)
			}

			// Populate computed values
			bpHealthProbeSettings := make(map[string]string)
			if props.HealthProbeSettings != nil {
				for _, v := range *props.HealthProbeSettings {
					if v.Name == nil || v.Id == nil {
						continue
					}
					rid, err := parse.HealthProbeIDInsensitively(*v.Id)
					if err != nil {
						continue
					}
					bpHealthProbeSettings[*v.Name] = rid.ID()
				}
			}
			if err := d.Set("backend_pool_health_probes", bpHealthProbeSettings); err != nil {
				return fmt.Errorf("setting `backend_pool_health_probes`: %+v", err)
			}

			bpLBSettings := make(map[string]string)
			if props.LoadBalancingSettings != nil {
				for _, v := range *props.LoadBalancingSettings {
					if v.Name == nil || v.Id == nil {
						continue
					}
					rid, err := parse.LoadBalancingIDInsensitively(*v.Id)
					if err != nil {
						continue
					}
					bpLBSettings[*v.Name] = rid.ID()
				}
			}
			if err := d.Set("backend_pool_load_balancing_settings", bpLBSettings); err != nil {
				return fmt.Errorf("setting `backend_pool_load_balancing_settings`: %+v", err)
			}

			backendPools := make(map[string]string)
			if props.BackendPools != nil {
				for _, v := range *props.BackendPools {
					if v.Name == nil || v.Id == nil {
						continue
					}
					rid, err := parse.BackendPoolIDInsensitively(*v.Id)
					if err != nil {
						continue
					}
					backendPools[*v.Name] = rid.ID()
				}
			}
			if err := d.Set("backend_pools", backendPools); err != nil {
				return fmt.Errorf("setting `backend_pools`: %+v", err)
			}

			frontendEndpoints := make(map[string]string)
			if props.FrontendEndpoints != nil {
				for _, v := range *props.FrontendEndpoints {
					if v.Name == nil || v.Id == nil {
						continue
					}
					rid, err := parse.FrontendEndpointIDInsensitively(*v.Id)
					if err != nil {
						continue
					}
					frontendEndpoints[*v.Name] = rid.ID()
				}
			}
			if err := d.Set("frontend_endpoints", frontendEndpoints); err != nil {
				return fmt.Errorf("setting `frontend_endpoints`: %+v", err)
			}

			routingRules := make(map[string]string)
			if props.RoutingRules != nil {
				for _, v := range *props.RoutingRules {
					if v.Name == nil || v.Id == nil {
						continue
					}
					rid, err := parse.RoutingRuleIDInsensitively(*v.Id)
					if err != nil {
						continue
					}
					routingRules[*v.Name] = rid.ID()
				}
			}
			if err := d.Set("routing_rules", routingRules); err != nil {
				return fmt.Errorf("setting `routing_rules`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceFrontDoorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontDoorIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandFrontDoorBackendPools(input []interface{}, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.BackendPool {
	if len(input) == 0 {
		return &[]frontdoors.BackendPool{}
	}

	output := make([]frontdoors.BackendPool, 0)

	for _, bp := range input {
		backendPool := bp.(map[string]interface{})
		backendPoolName := backendPool["name"].(string)
		backendPoolLoadBalancingName := backendPool["load_balancing_name"].(string)
		backendPoolHealthProbeName := backendPool["health_probe_name"].(string)
		backends := backendPool["backend"].([]interface{})

		backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolName).ID()
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolHealthProbeName).ID()
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolLoadBalancingName).ID()

		result := frontdoors.BackendPool{
			Id:   utils.String(backendPoolId),
			Name: utils.String(backendPoolName),
			Properties: &frontdoors.BackendPoolProperties{
				Backends: expandFrontDoorBackend(backends),
				HealthProbeSettings: &frontdoors.SubResource{
					Id: utils.String(healthProbeId),
				},
				LoadBalancingSettings: &frontdoors.SubResource{
					Id: utils.String(loadBalancingId),
				},
			},
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorBackend(input []interface{}) *[]frontdoors.Backend {
	if len(input) == 0 {
		return &[]frontdoors.Backend{}
	}

	output := make([]frontdoors.Backend, 0)

	for _, be := range input {
		backend := be.(map[string]interface{})
		address := backend["address"].(string)
		hostHeader := backend["host_header"].(string)
		enabled := expandFrontDoorBackendEnabledState(backend["enabled"].(bool))
		httpPort := int64(backend["http_port"].(int))
		httpsPort := int64(backend["https_port"].(int))
		priority := int64(backend["priority"].(int))
		weight := int64(backend["weight"].(int))

		result := frontdoors.Backend{
			Address:           utils.String(address),
			BackendHostHeader: utils.String(hostHeader),
			EnabledState:      &enabled,
			HTTPPort:          utils.Int64(httpPort),
			HTTPSPort:         utils.Int64(httpsPort),
			Priority:          utils.Int64(priority),
			Weight:            utils.Int64(weight),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorBackendEnabledState(isEnabled bool) frontdoors.BackendEnabledState {
	if isEnabled {
		return frontdoors.BackendEnabledStateEnabled
	}
	return frontdoors.BackendEnabledStateDisabled
}

func expandFrontDoorBackendPoolsSettings(enforceCertificateNameCheck bool, backendPoolsSendReceiveTimeoutSeconds int64) *frontdoors.BackendPoolsSettings {
	enforceCheck := frontdoors.EnforceCertificateNameCheckEnabledStateDisabled

	if enforceCertificateNameCheck {
		enforceCheck = frontdoors.EnforceCertificateNameCheckEnabledStateEnabled
	}

	result := frontdoors.BackendPoolsSettings{
		EnforceCertificateNameCheck: &enforceCheck,
		SendRecvTimeoutSeconds:      utils.Int64(backendPoolsSendReceiveTimeoutSeconds),
	}

	return &result
}

func expandFrontDoorFrontendEndpoint(input []interface{}, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.FrontendEndpoint {
	if len(input) == 0 {
		return &[]frontdoors.FrontendEndpoint{}
	}

	output := make([]frontdoors.FrontendEndpoint, 0)

	for _, frontendEndpoints := range input {
		frontendEndpoint := frontendEndpoints.(map[string]interface{})
		hostName := frontendEndpoint["host_name"].(string)
		isSessionAffinityEnabled := frontendEndpoint["session_affinity_enabled"].(bool)
		sessionAffinityTtlSeconds := int64(frontendEndpoint["session_affinity_ttl_seconds"].(int))
		waf := frontendEndpoint["web_application_firewall_policy_link_id"].(string)
		name := frontendEndpoint["name"].(string)
		id := frontdoors.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
		sessionAffinityEnabled := frontdoors.SessionAffinityEnabledStateDisabled

		if isSessionAffinityEnabled {
			sessionAffinityEnabled = frontdoors.SessionAffinityEnabledStateEnabled
		}

		result := frontdoors.FrontendEndpoint{
			Id:   utils.String(id),
			Name: utils.String(name),
			Properties: &frontdoors.FrontendEndpointProperties{
				HostName:                    utils.String(hostName),
				SessionAffinityEnabledState: &sessionAffinityEnabled,
				SessionAffinityTtlSeconds:   utils.Int64(sessionAffinityTtlSeconds),
			},
		}

		if waf != "" {
			result.Properties.WebApplicationFirewallPolicyLink = &frontdoors.FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink{
				Id: utils.String(waf),
			}
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorHealthProbeSettingsModel(input []interface{}, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.HealthProbeSettingsModel {
	if len(input) == 0 {
		return &[]frontdoors.HealthProbeSettingsModel{}
	}

	output := make([]frontdoors.HealthProbeSettingsModel, 0)

	for _, hps := range input {
		v := hps.(map[string]interface{})
		path := v["path"].(string)
		protocol := frontdoors.FrontDoorProtocol(v["protocol"].(string))
		intervalInSeconds := int64(v["interval_in_seconds"].(int))
		name := v["name"].(string)
		enabled := v["enabled"].(bool)

		healthProbeEnabled := frontdoors.HealthProbeEnabledEnabled
		if !enabled {
			healthProbeEnabled = frontdoors.HealthProbeEnabledDisabled
		}
		healthProbeId := parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		probeMethod := frontdoors.FrontDoorHealthProbeMethod(v["probe_method"].(string))

		result := frontdoors.HealthProbeSettingsModel{
			Id:   utils.String(healthProbeId),
			Name: utils.String(name),
			Properties: &frontdoors.HealthProbeSettingsProperties{
				IntervalInSeconds: utils.Int64(intervalInSeconds),
				Path:              utils.String(path),
				Protocol:          &protocol,
				HealthProbeMethod: &probeMethod,
				EnabledState:      &healthProbeEnabled,
			},
		}

		output = append(output, result)
	}

	return &output
}

func expandFrontDoorLoadBalancingSettingsModel(input []interface{}, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.LoadBalancingSettingsModel {
	if len(input) == 0 {
		return &[]frontdoors.LoadBalancingSettingsModel{}
	}

	output := make([]frontdoors.LoadBalancingSettingsModel, 0)

	for _, lbs := range input {
		loadBalanceSetting := lbs.(map[string]interface{})
		name := loadBalanceSetting["name"].(string)
		sampleSize := int64(loadBalanceSetting["sample_size"].(int))
		successfulSamplesRequired := int64(loadBalanceSetting["successful_samples_required"].(int))
		additionalLatencyMilliseconds := int64(loadBalanceSetting["additional_latency_milliseconds"].(int))
		loadBalancingId := parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		result := frontdoors.LoadBalancingSettingsModel{
			Id:   utils.String(loadBalancingId),
			Name: utils.String(name),
			Properties: &frontdoors.LoadBalancingSettingsProperties{
				SampleSize:                    utils.Int64(sampleSize),
				SuccessfulSamplesRequired:     utils.Int64(successfulSamplesRequired),
				AdditionalLatencyMilliseconds: utils.Int64(additionalLatencyMilliseconds),
			},
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorRoutingRule(input []interface{}, frontDoorId frontdoors.FrontDoorId, rulesEngines *map[string]*frontdoors.SubResource) *[]frontdoors.RoutingRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]frontdoors.RoutingRule, 0)

	for _, rr := range input {
		routingRule := rr.(map[string]interface{})
		name := routingRule["name"].(string)
		frontendEndpoints := routingRule["frontend_endpoints"].([]interface{})
		acceptedProtocols := routingRule["accepted_protocols"].([]interface{})
		ptm := routingRule["patterns_to_match"].([]interface{})
		enabled := frontdoors.RoutingRuleEnabledState(expandFrontDoorEnabledState(routingRule["enabled"].(bool)))

		patternsToMatch := make([]string, 0)
		for _, p := range ptm {
			patternsToMatch = append(patternsToMatch, p.(string))
		}

		var routingConfiguration frontdoors.RouteConfiguration
		if rc := routingRule["redirect_configuration"].([]interface{}); len(rc) != 0 {
			routingConfiguration = expandFrontDoorRedirectConfiguration(rc)
		} else if fc := routingRule["forwarding_configuration"].([]interface{}); len(fc) != 0 {
			routingConfiguration = expandFrontDoorForwardingConfiguration(fc, frontDoorId)
		}
		routingRuleId := parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()

		currentRoutingRule := frontdoors.RoutingRule{
			Id:   utils.String(routingRuleId),
			Name: utils.String(name),
			Properties: &frontdoors.RoutingRuleProperties{
				FrontendEndpoints:  expandFrontDoorFrontEndEndpoints(frontendEndpoints, frontDoorId),
				AcceptedProtocols:  expandFrontDoorAcceptedProtocols(acceptedProtocols),
				PatternsToMatch:    &patternsToMatch,
				EnabledState:       &enabled,
				RouteConfiguration: routingConfiguration,
			},
		}

		// Preserve existing rules engine for this routing rule
		// https://github.com/hashicorp/terraform-provider-azurerm/issues/7455#issuecomment-882769364
		if rulesEngines != nil {
			if rulesEngine, ok := (*rulesEngines)[name]; ok {
				currentRoutingRule.Properties.RulesEngine = rulesEngine
			}
		}

		output = append(output, currentRoutingRule)
	}

	return &output
}

func expandFrontDoorAcceptedProtocols(input []interface{}) *[]frontdoors.FrontDoorProtocol {
	if len(input) == 0 {
		return &[]frontdoors.FrontDoorProtocol{}
	}

	output := make([]frontdoors.FrontDoorProtocol, 0)

	for _, ap := range input {
		result := frontdoors.FrontDoorProtocolHTTPS
		if ap.(string) == string(frontdoors.FrontDoorProtocolHTTP) {
			result = frontdoors.FrontDoorProtocolHTTP
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorFrontEndEndpoints(input []interface{}, frontDoorId frontdoors.FrontDoorId) *[]frontdoors.SubResource {
	if len(input) == 0 {
		return &[]frontdoors.SubResource{}
	}

	output := make([]frontdoors.SubResource, 0)

	for _, name := range input {
		frontendEndpointId := parse.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name.(string)).ID()
		result := frontdoors.SubResource{
			Id: utils.String(frontendEndpointId),
		}
		output = append(output, result)
	}

	return &output
}

func expandFrontDoorEnabledState(enabled bool) frontdoors.FrontDoorEnabledState {
	if enabled {
		return frontdoors.FrontDoorEnabledStateEnabled
	}
	return frontdoors.FrontDoorEnabledStateDisabled
}

func expandFrontDoorRedirectConfiguration(input []interface{}) frontdoors.RedirectConfiguration {
	if len(input) == 0 {
		return frontdoors.RedirectConfiguration{}
	}

	v := input[0].(map[string]interface{})
	redirectType := frontdoors.FrontDoorRedirectType(v["redirect_type"].(string))
	redirectProtocol := frontdoors.FrontDoorRedirectProtocol(v["redirect_protocol"].(string))
	customHost := v["custom_host"].(string)
	customPath := v["custom_path"].(string)
	customFragment := v["custom_fragment"].(string)
	customQueryString := v["custom_query_string"].(string)

	redirectConfiguration := frontdoors.RedirectConfiguration{
		CustomHost:       utils.String(customHost),
		RedirectType:     &redirectType,
		RedirectProtocol: &redirectProtocol,
	}
	// The way the API works is if you don't include the attribute in the structure
	// it is treated as Preserve instead of Replace...
	if customHost != "" {
		redirectConfiguration.CustomHost = utils.String(customHost)
	}
	if customPath != "" {
		redirectConfiguration.CustomPath = utils.String(customPath)
	}
	if customFragment != "" {
		redirectConfiguration.CustomFragment = utils.String(customFragment)
	}
	if customQueryString != "" {
		redirectConfiguration.CustomQueryString = utils.String(customQueryString)
	}
	return redirectConfiguration
}

func expandFrontDoorForwardingConfiguration(input []interface{}, frontDoorId frontdoors.FrontDoorId) frontdoors.ForwardingConfiguration {
	if len(input) == 0 {
		return frontdoors.ForwardingConfiguration{}
	}

	v := input[0].(map[string]interface{})
	customForwardingPath := v["custom_forwarding_path"].(string)
	forwardingProtocol := frontdoors.FrontDoorForwardingProtocol(v["forwarding_protocol"].(string))
	backendPoolName := v["backend_pool_name"].(string)
	cacheUseDynamicCompression := v["cache_use_dynamic_compression"].(bool)
	cacheQueryParameterStripDirective := frontdoors.FrontDoorQuery(v["cache_query_parameter_strip_directive"].(string))
	cacheQueryParameters := v["cache_query_parameters"].([]interface{})
	cacheDuration := v["cache_duration"].(string)
	cacheEnabled := v["cache_enabled"].(bool)

	// convert list of cache_query_parameters into an array into a comma-separated list
	queryParametersArray := make([]string, 0)
	for _, p := range cacheQueryParameters {
		queryParametersArray = append(queryParametersArray, p.(string))
	}
	queryParametersString := strings.Join(queryParametersArray, ",")

	backendPoolId := parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, backendPoolName).ID()
	backend := &frontdoors.SubResource{
		Id: utils.String(backendPoolId),
	}

	forwardingConfiguration := frontdoors.ForwardingConfiguration{
		ForwardingProtocol: &forwardingProtocol,
		BackendPool:        backend,
	}
	// Per the portal, if you enable the cache the cache_query_parameter_strip_directive
	// is then a required attribute else the CacheConfiguration type is null
	if cacheEnabled {
		// Set the default value for dynamic compression or use the value defined in the config
		dynamicCompression := frontdoors.DynamicCompressionEnabledEnabled
		if !cacheUseDynamicCompression {
			dynamicCompression = frontdoors.DynamicCompressionEnabledDisabled
		}
		if cacheQueryParameterStripDirective == "" {
			// Set Default Value for strip directive is not in the key slice and cache is enabled
			cacheQueryParameterStripDirective = frontdoors.FrontDoorQueryStripAll
		}
		// set cacheQueryParameters to "" when StripDirective is "StripAll" or "StripNone"
		if cacheQueryParameterStripDirective == "StripAll" || cacheQueryParameterStripDirective == "StripNone" {
			queryParametersString = ""
		}

		// Making sure that duration is empty when cacheDuration is empty
		var duration *string
		if cacheDuration == "" {
			duration = nil
		} else {
			duration = utils.String(cacheDuration)
		}

		forwardingConfiguration.CacheConfiguration = &frontdoors.CacheConfiguration{
			DynamicCompression:           &dynamicCompression,
			QueryParameterStripDirective: &cacheQueryParameterStripDirective,
			QueryParameters:              utils.String(queryParametersString),
			CacheDuration:                duration,
		}
	}

	if customForwardingPath != "" {
		forwardingConfiguration.CustomForwardingPath = utils.String(customForwardingPath)
	}

	return forwardingConfiguration
}

func flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings []interface{}, frontDoorId frontdoors.FrontDoorId) *[]interface{} {
	output := make([]interface{}, 0)
	var backendPoolOrder []string
	var frontedEndpointOrder []string
	var routingRulesOrder []string
	var backendPoolLoadBalancingOrder []string
	var backendPoolHealthProbeOrder []string
	if len(backendPools) > 0 {
		flattenendBackendPools, err := flattenFrontDoorBackendPools(expandFrontDoorBackendPools(backendPools, frontDoorId), frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendBackendPools {
				backendPool := ids.(map[string]interface{})
				backendPoolOrder = append(backendPoolOrder, backendPool["id"].(string))
			}
		}
	}
	if len(frontendEndpoints) > 0 {
		flattenendfrontendEndpoints, err := flattenFrontEndEndpoints(expandFrontDoorFrontendEndpoint(frontendEndpoints, frontDoorId), frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendfrontendEndpoints {
				frontendEndPoint := ids.(map[string]interface{})
				frontedEndpointOrder = append(frontedEndpointOrder, frontendEndPoint["id"].(string))
			}
		}
	}
	if len(routingRules) > 0 {
		var oldBlocks interface{}
		flattenendRoutingRules, err := flattenFrontDoorRoutingRule(expandFrontDoorRoutingRule(routingRules, frontDoorId, nil), oldBlocks, frontDoorId, make([]interface{}, 0))
		if err == nil {
			for _, ids := range *flattenendRoutingRules {
				routingRule := ids.(map[string]interface{})
				routingRulesOrder = append(routingRulesOrder, routingRule["id"].(string))
			}
		}
	}
	if len(loadBalancingSettings) > 0 {
		flattenendLoadBalancingSettings := flattenFrontDoorLoadBalancingSettingsModel(expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, frontDoorId), frontDoorId, make([]interface{}, 0))

		if len(flattenendLoadBalancingSettings) > 0 {
			for _, ids := range flattenendLoadBalancingSettings {
				loadBalancingSetting := ids.(map[string]interface{})
				backendPoolLoadBalancingOrder = append(backendPoolLoadBalancingOrder, loadBalancingSetting["id"].(string))
			}
		}
	}
	if len(healthProbeSettings) > 0 {
		flattenendHealthProbeSettings := flattenFrontDoorHealthProbeSettingsModel(expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, frontDoorId), frontDoorId, make([]interface{}, 0))

		if len(flattenendHealthProbeSettings) > 0 {
			for _, ids := range flattenendHealthProbeSettings {
				healthProbeSetting := ids.(map[string]interface{})
				backendPoolHealthProbeOrder = append(backendPoolHealthProbeOrder, healthProbeSetting["id"].(string))
			}
		}
	}

	output = append(output, map[string]interface{}{
		"backend_pool_ids":                backendPoolOrder,
		"frontend_endpoint_ids":           frontedEndpointOrder,
		"routing_rule_ids":                routingRulesOrder,
		"backend_pool_load_balancing_ids": backendPoolLoadBalancingOrder,
		"backend_pool_health_probe_ids":   backendPoolHealthProbeOrder,
	})

	return &output
}

func combineBackendPools(allPools []frontdoors.BackendPool, orderedIds []interface{}, frontDoorId frontdoors.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, backend := range allPools {
			if strings.EqualFold(v.(string), *backend.Id) {
				orderedBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
				if err == nil {
					output = append(output, orderedBackendPool)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, backend := range allPools {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *backend.Id) {
				found = true
				break
			}
		}

		if !found {
			newBackendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, newBackendPool)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontDoorBackendPools(input *[]frontdoors.BackendPool, frontDoorId frontdoors.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedBackendPools := explicitOrder[0].(map[string]interface{})
		orderedBackendPoolsIds := orderedBackendPools["backend_pool_ids"].([]interface{})
		combinedBackendPools, err := combineBackendPools(*input, orderedBackendPoolsIds, frontDoorId)
		if err == nil {
			output = combinedBackendPools
		} else {
			return nil, err
		}
	} else {
		for _, backend := range *input {
			backendPool, err := flattenSingleFrontDoorBackendPools(&backend, frontDoorId)
			if err == nil {
				output = append(output, backendPool)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorBackendPools(input *frontdoors.BackendPool, frontDoorId frontdoors.FrontDoorId) (map[string]interface{}, error) {
	if input == nil {
		return make(map[string]interface{}), nil
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewBackendPoolID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	backend := make([]interface{}, 0)
	healthProbeName := ""
	loadBalancingName := ""
	if props := input.Properties; props != nil {
		backend = flattenFrontDoorBackend(props.Backends)
		if props.HealthProbeSettings != nil && props.HealthProbeSettings.Id != nil {
			name, err := parse.HealthProbeIDInsensitively(*props.HealthProbeSettings.Id)
			if err != nil {
				return nil, err
			}
			healthProbeName = name.HealthProbeSettingName
		}

		if props.LoadBalancingSettings != nil && props.LoadBalancingSettings.Id != nil {
			name, err := parse.LoadBalancingIDInsensitively(*props.LoadBalancingSettings.Id)
			if err != nil {
				return nil, err
			}
			loadBalancingName = name.LoadBalancingSettingName
		}
	}

	output := map[string]interface{}{
		"backend":             backend,
		"health_probe_name":   healthProbeName,
		"id":                  id,
		"load_balancing_name": loadBalancingName,
		"name":                name,
	}

	return output, nil
}

type flattenedBackendPoolSettings struct {
	enforceBackendPoolsCertificateNameCheck bool
	backendPoolsSendReceiveTimeoutSeconds   int
}

func flattenFrontDoorBackendPoolsSettings(input *frontdoors.BackendPoolsSettings) flattenedBackendPoolSettings {
	if input == nil {
		return flattenedBackendPoolSettings{
			enforceBackendPoolsCertificateNameCheck: true,
			backendPoolsSendReceiveTimeoutSeconds:   60,
		}
	}

	enforceCertificateNameCheck := false
	sendReceiveTimeoutSeconds := 0
	if input.EnforceCertificateNameCheck != nil && *input.EnforceCertificateNameCheck != "" && *input.EnforceCertificateNameCheck == frontdoors.EnforceCertificateNameCheckEnabledStateEnabled {
		enforceCertificateNameCheck = true
	}
	if input.SendRecvTimeoutSeconds != nil {
		sendReceiveTimeoutSeconds = int(*input.SendRecvTimeoutSeconds)
	}

	return flattenedBackendPoolSettings{
		enforceBackendPoolsCertificateNameCheck: enforceCertificateNameCheck,
		backendPoolsSendReceiveTimeoutSeconds:   sendReceiveTimeoutSeconds,
	}
}

func flattenFrontDoorBackend(input *[]frontdoors.Backend) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	output := make([]interface{}, 0)
	for _, v := range *input {
		result := make(map[string]interface{})
		if address := v.Address; address != nil {
			result["address"] = *address
		}
		if backendHostHeader := v.BackendHostHeader; backendHostHeader != nil {
			result["host_header"] = *backendHostHeader
		}
		if v.EnabledState != nil {
			result["enabled"] = *v.EnabledState == frontdoors.BackendEnabledStateEnabled
		}
		if httpPort := v.HTTPPort; httpPort != nil {
			result["http_port"] = int(*httpPort)
		}
		if httpsPort := v.HTTPSPort; httpsPort != nil {
			result["https_port"] = int(*httpsPort)
		}
		if priority := v.Priority; priority != nil {
			result["priority"] = int(*priority)
		}
		if weight := v.Weight; weight != nil {
			result["weight"] = int(*weight)
		}
		output = append(output, result)
	}

	return output
}

func retrieveFrontEndEndpointInformation(ctx context.Context, client *frontdoors.FrontDoorsClient, frontDoorId frontdoors.FrontDoorId, endpoints *[]frontdoors.FrontendEndpoint) (*[]frontdoors.FrontendEndpoint, error) {
	output := make([]frontdoors.FrontendEndpoint, 0)
	if endpoints == nil {
		return &output, nil
	}

	for _, endpoint := range *endpoints {
		if endpoint.Name == nil {
			continue
		}

		name := *endpoint.Name
		endpointID := frontdoors.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name)
		resp, err := client.FrontendEndpointsGet(ctx, endpointID)
		if err != nil {
			return nil, fmt.Errorf("retrieving Custom HTTPS Configuration for Frontend Endpoint %q (FrontDoor %q / Resource Group %q): %+v", name, frontDoorId.FrontDoorName, frontDoorId.ResourceGroupName, err)
		}
		if resp.Model != nil {
			output = append(output, *resp.Model)
		}
	}

	return &output, nil
}

func combineFrontEndEndpoints(allEndpoints []frontdoors.FrontendEndpoint, orderedIds []interface{}, frontDoorId frontdoors.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, frontendEndpoint := range allEndpoints {
			if strings.EqualFold(v.(string), *frontendEndpoint.Id) {
				orderedFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
				if err == nil {
					output = append(output, orderedFrontendEndpoint)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, frontendEndpoint := range allEndpoints {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *frontendEndpoint.Id) {
				found = true
				break
			}
		}

		if !found {
			newFrontendEndpoint, err := flattenSingleFrontEndEndpoints(frontendEndpoint, frontDoorId)
			if err == nil {
				output = append(output, newFrontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontEndEndpoints(input *[]frontdoors.FrontendEndpoint, frontDoorId frontdoors.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	output := make([]interface{}, 0)
	if input == nil {
		return &output, nil
	}

	if len(explicitOrder) > 0 {
		orderedFrontEnd := explicitOrder[0].(map[string]interface{})
		orderedFrontEndIds := orderedFrontEnd["frontend_endpoint_ids"].([]interface{})
		combinedFrontEndEndpoints, err := combineFrontEndEndpoints(*input, orderedFrontEndIds, frontDoorId)
		if err == nil {
			output = combinedFrontEndEndpoints
		} else {
			return nil, err
		}
	} else {
		for _, v := range *input {
			frontendEndpoint, err := flattenSingleFrontEndEndpoints(v, frontDoorId)
			if err == nil {
				output = append(output, frontendEndpoint)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontEndEndpoints(input frontdoors.FrontendEndpoint, frontDoorId frontdoors.FrontDoorId) (map[string]interface{}, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = parse.NewFrontendEndpointID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, *input.Name).ID()
		name = *input.Name
	}
	hostName := ""
	sessionAffinityEnabled := false
	sessionAffinityTlsSeconds := 0
	webApplicationFirewallPolicyLinkId := ""
	if props := input.Properties; props != nil {
		if props.HostName != nil {
			hostName = *props.HostName
		}
		if props.SessionAffinityEnabledState != nil && *props.SessionAffinityEnabledState != "" {
			sessionAffinityEnabled = *props.SessionAffinityEnabledState == frontdoors.SessionAffinityEnabledStateEnabled
		}
		if props.SessionAffinityTtlSeconds != nil {
			sessionAffinityTlsSeconds = int(*props.SessionAffinityTtlSeconds)
		}
		if waf := props.WebApplicationFirewallPolicyLink; waf != nil && waf.Id != nil {
			// rewrite the ID to ensure it's consistent
			parsed, err := parse.WebApplicationFirewallPolicyIDInsensitively(*waf.Id)
			if err != nil {
				return nil, err
			}
			webApplicationFirewallPolicyLinkId = parsed.ID()
		}
		// flattenedHttpsConfig := flattenCustomHttpsConfiguration(props)
		// customHTTPSConfiguration = flattenedHttpsConfig.CustomHTTPSConfiguration
		// customHttpsProvisioningEnabled = flattenedHttpsConfig.CustomHTTPSProvisioningEnabled
	}

	output := map[string]interface{}{
		// "custom_https_configuration":        customHTTPSConfiguration,
		// "custom_https_provisioning_enabled": customHttpsProvisioningEnabled,
		"host_name":                    hostName,
		"id":                           id,
		"name":                         name,
		"session_affinity_enabled":     sessionAffinityEnabled,
		"session_affinity_ttl_seconds": sessionAffinityTlsSeconds,
		"web_application_firewall_policy_link_id": webApplicationFirewallPolicyLinkId,
	}

	return output, nil
}

func combineHealthProbeSettingsModel(allHealthProbeSettings []frontdoors.HealthProbeSettingsModel, orderedIds []interface{}, frontDoorId frontdoors.FrontDoorId) []interface{} {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, healthProbeSetting := range allHealthProbeSettings {
			if strings.EqualFold(v.(string), *healthProbeSetting.Id) {
				orderedHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
				output = append(output, orderedHealthProbeSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, healthProbeSetting := range allHealthProbeSettings {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *healthProbeSetting.Id) {
				found = true
				break
			}
		}

		if !found {
			newHealthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&healthProbeSetting, frontDoorId)
			output = append(output, newHealthProbeSetting)
		}
	}

	return output
}

func flattenFrontDoorHealthProbeSettingsModel(input *[]frontdoors.HealthProbeSettingsModel, frontDoorId frontdoors.FrontDoorId, explicitOrder []interface{}) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	if len(explicitOrder) > 0 {
		orderedHealthProbeSetting := explicitOrder[0].(map[string]interface{})
		orderedHealthProbeSettingIds := orderedHealthProbeSetting["backend_pool_health_probe_ids"].([]interface{})
		output = combineHealthProbeSettingsModel(*input, orderedHealthProbeSettingIds, frontDoorId)
	} else {
		for _, v := range *input {
			healthProbeSetting := flattenSingleFrontDoorHealthProbeSettingsModel(&v, frontDoorId)
			output = append(output, healthProbeSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorHealthProbeSettingsModel(input *frontdoors.HealthProbeSettingsModel, frontDoorId frontdoors.FrontDoorId) map[string]interface{} {
	if input == nil {
		return make(map[string]interface{})
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewHealthProbeID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	enabled := false
	intervalInSeconds := 0
	path := ""
	probeMethod := ""
	protocol := ""

	if properties := input.Properties; properties != nil {
		if properties.IntervalInSeconds != nil {
			intervalInSeconds = int(*properties.IntervalInSeconds)
		}
		if properties.Path != nil {
			path = *properties.Path
		}
		if healthProbeMethod := properties.HealthProbeMethod; healthProbeMethod != nil && *healthProbeMethod != "" {
			// I have to upper this as the frontdoor.GET and frontdoor.HEAD types are uppercased
			// but Azure stores them in the resource as sentence cased (e.g. "Get" and "Head")
			probeMethod = strings.ToUpper(string(*healthProbeMethod))
		}
		if properties.EnabledState != nil && *properties.EnabledState != "" {
			enabled = *properties.EnabledState == frontdoors.HealthProbeEnabledEnabled
		}
		if properties.Protocol != nil {
			protocol = string(*properties.Protocol)
		}
	}

	output := map[string]interface{}{
		"enabled":             enabled,
		"id":                  id,
		"name":                name,
		"protocol":            protocol,
		"interval_in_seconds": intervalInSeconds,
		"path":                path,
		"probe_method":        probeMethod,
	}

	return output
}

func combineLoadBalancingSettingsModel(allLoadBalancingSettings []frontdoors.LoadBalancingSettingsModel, orderedIds []interface{}, frontDoorId frontdoors.FrontDoorId) []interface{} {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, loadBalancingSetting := range allLoadBalancingSettings {
			if strings.EqualFold(v.(string), *loadBalancingSetting.Id) {
				orderedLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalancingSetting, frontDoorId)
				output = append(output, orderedLoadBalanceSetting)
				break
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, loadBalanceSetting := range allLoadBalancingSettings {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *loadBalanceSetting.Id) {
				found = true
				break
			}
		}

		if !found {
			newLoadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&loadBalanceSetting, frontDoorId)
			output = append(output, newLoadBalanceSetting)
		}
	}

	return output
}

func flattenFrontDoorLoadBalancingSettingsModel(input *[]frontdoors.LoadBalancingSettingsModel, frontDoorId frontdoors.FrontDoorId, explicitOrder []interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedLoadBalancingSettings := explicitOrder[0].(map[string]interface{})
		orderedLoadBalancingIds := orderedLoadBalancingSettings["backend_pool_load_balancing_ids"].([]interface{})
		output = combineLoadBalancingSettingsModel(*input, orderedLoadBalancingIds, frontDoorId)
	} else {
		for _, v := range *input {
			loadBalanceSetting := flattenSingleFrontDoorLoadBalancingSettingsModel(&v, frontDoorId)
			output = append(output, loadBalanceSetting)
		}
	}

	return output
}

func flattenSingleFrontDoorLoadBalancingSettingsModel(input *frontdoors.LoadBalancingSettingsModel, frontDoorId frontdoors.FrontDoorId) map[string]interface{} {
	if input == nil {
		return make(map[string]interface{})
	}

	id := ""
	name := ""
	if input.Name != nil {
		name = *input.Name
		// rewrite the ID to ensure it's consistent
		id = parse.NewLoadBalancingID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, name).ID()
	}

	additionalLatencyMilliseconds := 0
	sampleSize := 0
	successfulSamplesRequired := 0
	if properties := input.Properties; properties != nil {
		if properties.AdditionalLatencyMilliseconds != nil {
			additionalLatencyMilliseconds = int(*properties.AdditionalLatencyMilliseconds)
		}
		if properties.SampleSize != nil {
			sampleSize = int(*properties.SampleSize)
		}
		if properties.SuccessfulSamplesRequired != nil {
			successfulSamplesRequired = int(*properties.SuccessfulSamplesRequired)
		}
	}

	output := map[string]interface{}{
		"additional_latency_milliseconds": additionalLatencyMilliseconds,
		"id":                              id,
		"name":                            name,
		"sample_size":                     sampleSize,
		"successful_samples_required":     successfulSamplesRequired,
	}

	return output
}

func combineRoutingRules(allRoutingRules []frontdoors.RoutingRule, oldBlocks interface{}, orderedIds []interface{}, frontDoorId frontdoors.FrontDoorId) ([]interface{}, error) {
	output := make([]interface{}, 0)
	found := false

	// first find all of the ones in the ordered mapping list and add them in the correct order
	for _, v := range orderedIds {
		for _, routingRule := range allRoutingRules {
			if strings.EqualFold(v.(string), *routingRule.Id) {
				orderedRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
				if err == nil {
					output = append(output, orderedRoutingRule)
					break
				} else {
					return nil, err
				}
			}
		}
	}

	// Now check to see if items were added to the resource via the portal and add them to the state file
	for _, routingRule := range allRoutingRules {
		found = false
		for _, orderedId := range orderedIds {
			if strings.EqualFold(orderedId.(string), *routingRule.Id) {
				found = true
				break
			}
		}

		if !found {
			newRoutingRule, err := flattenSingleFrontDoorRoutingRule(routingRule, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, newRoutingRule)
			} else {
				return nil, err
			}
		}
	}

	return output, nil
}

func flattenFrontDoorRoutingRule(input *[]frontdoors.RoutingRule, oldBlocks interface{}, frontDoorId frontdoors.FrontDoorId, explicitOrder []interface{}) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	if len(explicitOrder) > 0 {
		orderedRule := explicitOrder[0].(map[string]interface{})
		orderedRountingRuleIds := orderedRule["routing_rule_ids"].([]interface{})
		combinedRoutingRules, err := combineRoutingRules(*input, oldBlocks, orderedRountingRuleIds, frontDoorId)
		if err != nil {
			return nil, err
		}
		output = combinedRoutingRules
	} else {
		for _, v := range *input {
			routingRule, err := flattenSingleFrontDoorRoutingRule(v, oldBlocks, frontDoorId)
			if err == nil {
				output = append(output, routingRule)
			} else {
				return nil, err
			}
		}
	}

	return &output, nil
}

func flattenSingleFrontDoorRoutingRule(input frontdoors.RoutingRule, oldBlocks interface{}, frontDoorId frontdoors.FrontDoorId) (map[string]interface{}, error) {
	id := ""
	name := ""
	if input.Name != nil {
		// rewrite the ID to ensure it's consistent
		id = parse.NewRoutingRuleID(frontDoorId.SubscriptionId, frontDoorId.ResourceGroupName, frontDoorId.FrontDoorName, *input.Name).ID()
		name = *input.Name
	}

	acceptedProtocols := make([]string, 0)
	enabled := false
	forwardingConfiguration := make([]interface{}, 0)
	frontEndEndpoints := make([]string, 0)
	patternsToMatch := make([]string, 0)
	redirectConfiguration := make([]interface{}, 0)

	if props := input.Properties; props != nil {
		acceptedProtocols = flattenFrontDoorAcceptedProtocol(props.AcceptedProtocols)
		if props.EnabledState != nil {
			enabled = *props.EnabledState == frontdoors.RoutingRuleEnabledStateEnabled
		}
		forwardConfiguration, err := flattenRoutingRuleForwardingConfiguration(props.RouteConfiguration, oldBlocks)
		if err != nil {
			return nil, fmt.Errorf("flattening `forward_configuration`: %+v", err)
		}

		forwardingConfiguration = *forwardConfiguration
		frontendEndpoints, err := flattenFrontDoorFrontendEndpointsSubResources(props.FrontendEndpoints)
		if err != nil {
			return nil, fmt.Errorf("flattening `frontend_endpoints`: %+v", err)
		}

		frontEndEndpoints = *frontendEndpoints
		if props.PatternsToMatch != nil {
			patternsToMatch = *props.PatternsToMatch
		}
		redirectConfiguration = flattenRoutingRuleRedirectConfiguration(props.RouteConfiguration)
	}

	output := map[string]interface{}{
		"accepted_protocols":       acceptedProtocols,
		"enabled":                  enabled,
		"forwarding_configuration": forwardingConfiguration,
		"frontend_endpoints":       frontEndEndpoints,
		"id":                       id,
		"name":                     name,
		"patterns_to_match":        patternsToMatch,
		"redirect_configuration":   redirectConfiguration,
	}

	return output, nil
}

func flattenRoutingRuleForwardingConfiguration(config frontdoors.RouteConfiguration, oldConfig interface{}) (*[]interface{}, error) {
	v, ok := config.(frontdoors.ForwardingConfiguration)
	if !ok {
		return &[]interface{}{}, nil
	}

	name := ""
	if v.BackendPool != nil && v.BackendPool.Id != nil {
		backendPoolId, err := parse.BackendPoolIDInsensitively(*v.BackendPool.Id)
		if err != nil {
			return nil, err
		}
		name = backendPoolId.Name
	}
	customForwardingPath := ""
	if v.CustomForwardingPath != nil {
		customForwardingPath = *v.CustomForwardingPath
	}

	cacheEnabled := false
	cacheQueryParameterStripDirective := string(frontdoors.FrontDoorQueryStripAll)
	cacheUseDynamicCompression := false

	var cacheQueryParameters []interface{}
	var cacheQueryParametersArray []string
	var cacheDuration *string

	if cacheConfiguration := v.CacheConfiguration; cacheConfiguration != nil {
		cacheEnabled = true
		if stripDirective := cacheConfiguration.QueryParameterStripDirective; stripDirective != nil && *stripDirective != "" {
			cacheQueryParameterStripDirective = string(*stripDirective)
		}
		if dynamicCompression := cacheConfiguration.DynamicCompression; dynamicCompression != nil && *dynamicCompression != "" {
			cacheUseDynamicCompression = string(*dynamicCompression) == string(frontdoors.DynamicCompressionEnabledEnabled)
		}
		if queryParameters := cacheConfiguration.QueryParameters; queryParameters != nil {
			cacheQueryParametersArray = strings.Split(*queryParameters, ",")
		}
		if duration := cacheConfiguration.CacheDuration; duration != nil {
			cacheDuration = duration
		}
	} else {
		// if the cache is disabled, use the default values or revert to what they were in the previous plan
		old, ok := oldConfig.([]interface{})
		if ok {
			for _, oldValue := range old {
				oldVal, ok := oldValue.(map[string]interface{})
				if ok {
					thisName := oldVal["name"].(string)
					if name == thisName {
						oldConfigs := oldVal["forwarding_configuration"].([]interface{})
						if len(oldConfigs) > 0 {
							ofc := oldConfigs[0].(map[string]interface{})
							cacheQueryParameterStripDirective = ofc["cache_query_parameter_strip_directive"].(string)
							cacheUseDynamicCompression = ofc["cache_use_dynamic_compression"].(bool)
							cacheDuration = utils.String(ofc["cache_duration"].(string))

							cacheQueryParameters = ofc["cache_query_parameters"].([]interface{})
							for _, p := range cacheQueryParameters {
								cacheQueryParametersArray = append(cacheQueryParametersArray, p.(string))
							}
						}
					}
				}
			}
		}
	}

	forwardingProtocol := ""
	if v.ForwardingProtocol != nil {
		forwardingProtocol = string(*v.ForwardingProtocol)
	}

	return &[]interface{}{
		map[string]interface{}{
			"backend_pool_name":                     name,
			"custom_forwarding_path":                customForwardingPath,
			"forwarding_protocol":                   forwardingProtocol,
			"cache_enabled":                         cacheEnabled,
			"cache_query_parameter_strip_directive": cacheQueryParameterStripDirective,
			"cache_use_dynamic_compression":         cacheUseDynamicCompression,
			"cache_query_parameters":                cacheQueryParametersArray,
			"cache_duration":                        cacheDuration,
		},
	}, nil
}

func flattenRoutingRuleRedirectConfiguration(config frontdoors.RouteConfiguration) []interface{} {
	v, ok := config.(frontdoors.RedirectConfiguration)
	if !ok {
		return []interface{}{}
	}

	customFragment := ""
	if v.CustomFragment != nil {
		customFragment = *v.CustomFragment
	}
	customHost := ""
	if v.CustomHost != nil {
		customHost = *v.CustomHost
	}
	customQueryString := ""
	if v.CustomQueryString != nil {
		customQueryString = *v.CustomQueryString
	}
	customPath := ""
	if v.CustomPath != nil {
		customPath = *v.CustomPath
	}

	redirectProtocol := ""
	if v.RedirectProtocol != nil {
		redirectProtocol = string(*v.RedirectProtocol)
	}

	redirectType := ""
	if v.RedirectType != nil {
		redirectType = string(*v.RedirectType)
	}

	return []interface{}{
		map[string]interface{}{
			"custom_host":         customHost,
			"custom_fragment":     customFragment,
			"custom_query_string": customQueryString,
			"custom_path":         customPath,
			"redirect_protocol":   redirectProtocol,
			"redirect_type":       redirectType,
		},
	}
}

func flattenFrontDoorAcceptedProtocol(input *[]frontdoors.FrontDoorProtocol) []string {
	if input == nil {
		return make([]string, 0)
	}

	output := make([]string, 0)

	for _, p := range *input {
		output = append(output, string(p))
	}

	return output
}

func flattenFrontDoorFrontendEndpointsSubResources(input *[]frontdoors.SubResource) (*[]string, error) {
	output := make([]string, 0)
	if input == nil {
		return &output, nil
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		id, err := parse.FrontendEndpointIDInsensitively(*v.Id)
		if err != nil {
			return nil, err
		}
		output = append(output, id.Name)
	}

	return &output, nil
}

func resourceFrontDoorSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: frontDoorValidate.FrontDoorName,
		},

		"cname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"header_frontdoor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"load_balancer_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"routing_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 500,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"accepted_protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(frontdoors.FrontDoorProtocolHTTP),
								string(frontdoors.FrontDoorProtocolHTTPS),
							}, false),
						},
					},
					"patterns_to_match": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 25,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"frontend_endpoints": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 500,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"redirect_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_fragment": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_host": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_query_string": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"redirect_protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(frontdoors.FrontDoorRedirectProtocolHTTPOnly),
										string(frontdoors.FrontDoorRedirectProtocolHTTPSOnly),
										string(frontdoors.FrontDoorRedirectProtocolMatchRequest),
									}, false),
								},
								"redirect_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(frontdoors.FrontDoorRedirectTypeFound),
										string(frontdoors.FrontDoorRedirectTypeMoved),
										string(frontdoors.FrontDoorRedirectTypePermanentRedirect),
										string(frontdoors.FrontDoorRedirectTypeTemporaryRedirect),
									}, false),
								},
							},
						},
					},
					"forwarding_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"backend_pool_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
								},
								"cache_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"cache_use_dynamic_compression": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"cache_query_parameter_strip_directive": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(frontdoors.FrontDoorQueryStripAll),
									ValidateFunc: validation.StringInSlice([]string{
										string(frontdoors.FrontDoorQueryStripAll),
										string(frontdoors.FrontDoorQueryStripNone),
										string(frontdoors.FrontDoorQueryStripOnly),
										string(frontdoors.FrontDoorQueryStripAllExcept),
									}, false),
								},
								"cache_query_parameters": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"cache_duration": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.ISO8601DurationBetween("PT1S", "P365D"),
								},
								"custom_forwarding_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"forwarding_protocol": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(frontdoors.FrontDoorForwardingProtocolHTTPSOnly),
									ValidateFunc: validation.StringInSlice([]string{
										string(frontdoors.FrontDoorForwardingProtocolHTTPOnly),
										string(frontdoors.FrontDoorForwardingProtocolHTTPSOnly),
										string(frontdoors.FrontDoorForwardingProtocolMatchRequest),
									}, false),
								},
							},
						},
					},
				},
			},
		},

		"backend_pool_load_balancing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 5000,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
					},
					"sample_size": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  4,
					},
					"successful_samples_required": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  2,
					},
					"additional_latency_milliseconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},
				},
			},
		},

		"backend_pool_health_probe": {
			Type:     pluginsdk.TypeList,
			MaxItems: 5000,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "/",
					},
					"protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(frontdoors.FrontDoorProtocolHTTP),
						ValidateFunc: validation.StringInSlice([]string{
							string(frontdoors.FrontDoorProtocolHTTP),
							string(frontdoors.FrontDoorProtocolHTTPS),
						}, false),
					},
					"probe_method": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(frontdoors.FrontDoorHealthProbeMethodGET),
						ValidateFunc: validation.StringInSlice([]string{
							string(frontdoors.FrontDoorHealthProbeMethodGET),
							string(frontdoors.FrontDoorHealthProbeMethodHEAD),
						}, false),
					},
					"interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  120,
					},
				},
			},
		},

		"backend_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend": {
						Type:     pluginsdk.TypeList,
						MaxItems: 500,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
								"address": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"http_port": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"https_port": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"weight": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      50,
									ValidateFunc: validation.IntBetween(1, 1000),
								},
								"priority": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      1,
									ValidateFunc: validation.IntBetween(1, 5),
								},
								"host_header": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
					},
					"health_probe_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"load_balancing_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"backend_pool_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enforce_backend_pools_certificate_name_check": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"backend_pools_send_receive_timeout_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      60,
						ValidateFunc: validation.IntBetween(0, 240),
					},
				},
			},
		},

		"frontend_endpoint": {
			Type:     pluginsdk.TypeList,
			MaxItems: 500,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: frontDoorValidate.BackendPoolRoutingRuleName,
					},
					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"session_affinity_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"session_affinity_ttl_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},
					"web_application_firewall_policy_link_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: frontDoorValidate.WebApplicationFirewallPolicyID,
					},
				},
			},
		},

		// Computed values
		"explicit_resource_order": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend_pool_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"frontend_endpoint_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"routing_rule_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"backend_pool_load_balancing_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"backend_pool_health_probe_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"backend_pool_health_probes": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pool_load_balancing_settings": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pools": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"frontend_endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"routing_rules": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
