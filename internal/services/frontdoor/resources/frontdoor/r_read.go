// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorRead(d *pluginsdk.ResourceData, meta any) error {
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
			explicitResourceOrder := d.Get("explicit_resource_order").([]any)
			flattenedBackendPools, err := flattenFrontDoorBackendPools(props.BackendPools, *id, explicitResourceOrder)
			if err != nil {
				return fmt.Errorf("flattening `backend_pool`: %+v", err)
			}
			if err := d.Set("backend_pool", flattenedBackendPools); err != nil {
				return fmt.Errorf("setting `backend_pool`: %+v", err)
			}

			backendPoolSettings := flattenFrontDoorBackendPoolsSettings(props.BackendPoolsSettings)
			out := map[string]any{
				"enforce_backend_pools_certificate_name_check": backendPoolSettings.enforceBackendPoolsCertificateNameCheck,
				"backend_pools_send_receive_timeout_seconds":   backendPoolSettings.backendPoolsSendReceiveTimeoutSeconds,
			}
			d.Set("backend_pool_settings", []any{out})

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

			var flattenedRoutingRules *[]any
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
					rid, err := frontdoors.ParseFrontendEndpointIDInsensitively(*v.Id)
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

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
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
