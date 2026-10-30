// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubID(d.Id())
	if err != nil {
		return err
	}

	hub, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(hub.Response.Response) {
			log.Printf("[DEBUG] %s was not found!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.Name); err == nil {
		keyList := keysResp.Response()
		if err := d.Set("shared_access_policy", flattenIoTHubSharedAccessPolicy(keyList.Value)); err != nil {
			return fmt.Errorf("setting `shared_access_policy` in IoTHub %q: %+v", id.Name, err)
		}
	}

	if properties := hub.Properties; properties != nil {
		for k, v := range properties.EventHubEndpoints {
			if v == nil {
				continue
			}

			switch k {
			case "events":
				d.Set("event_hub_events_endpoint", v.Endpoint)

				if *v.Endpoint != "" {
					uri, err := url.Parse(*v.Endpoint)
					if err == nil {
						d.Set("event_hub_events_namespace", strings.Split(uri.Hostname(), ".")[0])
					}
				}

				d.Set("event_hub_events_path", v.Path)
				d.Set("event_hub_partition_count", v.PartitionCount)
				d.Set("event_hub_retention_in_days", v.RetentionTimeInDays)
			case "operationsMonitoringEvents":
				d.Set("event_hub_operations_endpoint", v.Endpoint)
				d.Set("event_hub_operations_path", v.Path)
			}
		}

		d.Set("hostname", properties.HostName)

		if err := d.Set("endpoint", flattenIoTHubEndpoint(properties.Routing)); err != nil {
			return fmt.Errorf("setting `endpoint` in IoTHub %q: %+v", id.Name, err)
		}

		if err := d.Set("route", flattenIoTHubRoute(properties.Routing)); err != nil {
			return fmt.Errorf("setting `route` in IoTHub %q: %+v", id.Name, err)
		}

		if err := d.Set("enrichment", flattenIoTHubEnrichment(properties.Routing)); err != nil {
			return fmt.Errorf("setting `enrichment` in IoTHub %q: %+v", id.Name, err)
		}

		if err := d.Set("fallback_route", flattenIoTHubFallbackRoute(properties.Routing)); err != nil {
			return fmt.Errorf("setting `fallbackRoute` in IoTHub %q: %+v", id.Name, err)
		}

		if err := d.Set("network_rule_set", flattenNetworkRuleSetProperties(properties.NetworkRuleSets)); err != nil {
			return fmt.Errorf("setting `network_rule_set` in IoTHub %q: %+v", id.Name, err)
		}

		if err := d.Set("file_upload", flattenIoTHubFileUpload(properties.StorageEndpoints, properties.MessagingEndpoints, properties.EnableFileUploadNotifications)); err != nil {
			return fmt.Errorf("setting `file_upload` in IoTHub %q: %+v", id.Name, err)
		}

		if enabled := properties.PublicNetworkAccess; enabled != "" {
			d.Set("public_network_access_enabled", enabled == devices.PublicNetworkAccessEnabled)
		}

		if err := d.Set("cloud_to_device", flattenIoTHubCloudToDevice(properties.CloudToDevice)); err != nil {
			return fmt.Errorf("setting `cloudToDevice` in IoTHub %q: %+v", id.Name, err)
		}

		d.Set("min_tls_version", properties.MinTLSVersion)

		localAuthenticationEnabled := true
		if properties.DisableLocalAuth != nil {
			localAuthenticationEnabled = !*properties.DisableLocalAuth
		}
		d.Set("local_authentication_enabled", localAuthenticationEnabled)
	}

	identity, err := flattenIotHubIdentity(hub.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if loc := hub.Location; loc != nil {
		d.Set("location", location.Normalize(*loc))
	}
	if err := d.Set("sku", flattenIoTHubSku(hub.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}
	d.Set("type", hub.Type)
	return tags.FlattenAndSet(d, hub.Tags)
}

//nolint:unparam
func suppressIfTypeIsNot(t string) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		path := strings.Split(k, ".")
		path[len(path)-1] = "type"
		return d.Get(strings.Join(path, ".")).(string) != t
	}
}

func suppressIfTypeIs(t string) pluginsdk.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *pluginsdk.ResourceData) bool {
		path := strings.Split(k, ".")
		path[len(path)-1] = "type"
		return d.Get(strings.Join(path, ".")).(string) == t
	}
}

func IothubConnectionStringSuppress(k, old, new string, d *pluginsdk.ResourceData) bool {
	secretKeyRegex := regexp.MustCompile("(SharedAccessKey|AccountKey)=[^;]+")
	sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

	// Azure will always mask the Access Keys and will include the port number in the GET response
	// 5671 is the default port for Azure Service Bus connections
	maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
	maskedNew = secretKeyRegex.ReplaceAllString(maskedNew, "$1=****")

	oldMap := connectionStringToMap(old)
	maskedNewMap := connectionStringToMap(maskedNew)
	if oldMap == nil || maskedNewMap == nil {
		return false
	}
	if len(oldMap) != len(maskedNewMap) {
		return false
	}
	for k, v := range oldMap {
		newV, ok := maskedNewMap[k]
		if !ok {
			return false
		}
		if newV != v {
			return false
		}
	}

	return true
}

func connectionStringToMap(connectionStr string) map[string]string {
	m := make(map[string]string)
	split := strings.SplitSeq(connectionStr, ";")
	for v := range split {
		// The connection string might contain `=`
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			continue
		}
		m[kv[0]] = kv[1]
	}
	return m
}
