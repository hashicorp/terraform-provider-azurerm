// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubDPS() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSCreate,
		Read:   resourceIotHubDPSRead,
		Update: resourceIotHubDPSUpdate,
		Delete: resourceIotHubDPSDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseProvisioningServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},

			"resource_group_name": commonschema.ResourceGroupName(), // azure.SchemaResourceGroupNameDiffSuppress(),

			"location": commonschema.Location(),

			"sku": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.IotHubSkuS1),
							}, false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 200),
						},
					},
				},
			},

			"linked_hub": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"connection_string": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							// Azure returns the key as ****. We'll suppress that here.
							DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
								secretKeyRegex := regexp.MustCompile("(SharedAccessKey)=[^;]+")
								maskedNew := secretKeyRegex.ReplaceAllString(new, "$1=****")
								return (new == d.Get(k).(string)) && (maskedNew == old)
							},
							Sensitive: true,
						},
						"location": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							StateFunc:    azure.NormalizeLocation,
						},
						"apply_allocation_policy": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"allocation_weight": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"hostname": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ip_filter_rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"ip_mask": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.CIDR,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.IPFilterActionTypeAccept),
								string(devices.IPFilterActionTypeReject),
							}, false),
						},
						"target": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(iotdpsresource.IPFilterTargetTypeAll),
								string(iotdpsresource.IPFilterTargetTypeServiceApi),
								string(iotdpsresource.IPFilterTargetTypeDeviceApi),
							}, false),
						},
					},
				},
			},

			"allocation_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(iotdpsresource.AllocationPolicyHashed),
				ValidateFunc: validation.StringInSlice([]string{
					string(iotdpsresource.AllocationPolicyHashed),
					string(iotdpsresource.AllocationPolicyGeoLatency),
					string(iotdpsresource.AllocationPolicyStatic),
				}, false),
			},

			"data_residency_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"device_provisioning_host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"id_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_operations_host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIotHubDPSCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing IoT Device Provisioning Service %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_iothub_dps", id.ID())
	}

	publicNetworkAccess := iotdpsresource.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = iotdpsresource.PublicNetworkAccessDisabled
	}

	allocationPolicy := iotdpsresource.AllocationPolicy(d.Get("allocation_policy").(string))
	iotdps := iotdpsresource.ProvisioningServiceDescription{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Name:     utils.String(id.ProvisioningServiceName),
		Sku:      expandIoTHubDPSSku(d),
		Properties: iotdpsresource.IotDpsPropertiesDescription{
			IotHubs:             expandIoTHubDPSIoTHubs(d.Get("linked_hub").([]interface{})),
			AllocationPolicy:    &allocationPolicy,
			EnableDataResidency: utils.Bool(d.Get("data_residency_enabled").(bool)),
			IPFilterRules:       expandDpsIPFilterRules(d),
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, iotdps); err != nil {
		return fmt.Errorf("creating IoT Device Provisioning Service %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSRead(d, meta)
}

func resourceIotHubDPSRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseProvisioningServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ProvisioningServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		sku := flattenIoTHubDPSSku(model.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		props := model.Properties
		if err := d.Set("linked_hub", flattenIoTHubDPSLinkedHub(props.IotHubs)); err != nil {
			return fmt.Errorf("setting `linked_hub`: %+v", err)
		}

		ipFilterRules := flattenDpsIPFilterRules(props.IPFilterRules)
		if err := d.Set("ip_filter_rule", ipFilterRules); err != nil {
			return fmt.Errorf("setting `ip_filter_rule` in IoTHub DPS %q: %+v", id.ProvisioningServiceName, err)
		}

		d.Set("service_operations_host_name", props.ServiceOperationsHostName)
		d.Set("device_provisioning_host_name", props.DeviceProvisioningHostName)
		d.Set("id_scope", props.IdScope)

		allocationPolicy := string(iotdpsresource.AllocationPolicyHashed)
		if props.AllocationPolicy != nil {
			allocationPolicy = string(*props.AllocationPolicy)
		}
		d.Set("allocation_policy", allocationPolicy)

		enableDataResidency := false
		if props.EnableDataResidency != nil {
			enableDataResidency = *props.EnableDataResidency
		}
		d.Set("data_residency_enabled", enableDataResidency)

		publicNetworkAccess := true
		if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != "" {
			publicNetworkAccess = strings.EqualFold("Enabled", string(*props.PublicNetworkAccess))
		}
		d.Set("public_network_access_enabled", publicNetworkAccess)

		d.Set("tags", flattenTags(model.Tags))
	}

	return nil
}

func resourceIotHubDPSUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	iotdps := resp.Model
	if iotdps == nil {
		return fmt.Errorf("retrieving model of %s: %+v", id, err)
	}

	if d.HasChanges("allocation_policy") {
		allocationPolicy := iotdpsresource.AllocationPolicy(d.Get("allocation_policy").(string))
		iotdps.Properties.AllocationPolicy = &allocationPolicy
	}

	if d.HasChanges("public_network_access_enabled") {
		publicNetworkAccess := iotdpsresource.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = iotdpsresource.PublicNetworkAccessDisabled
		}
		iotdps.Properties.PublicNetworkAccess = &publicNetworkAccess
	}

	if d.HasChanges("ip_filter_rule") {
		iotdps.Properties.IPFilterRules = expandDpsIPFilterRules(d)
	}

	if d.HasChanges("linked_hub") {
		iotdps.Properties.IotHubs = expandIoTHubDPSIoTHubs(d.Get("linked_hub").([]interface{}))
	}

	if d.HasChanges("sku") {
		iotdps.Sku = expandIoTHubDPSSku(d)
	}

	if d.HasChanges("tags") {
		iotdps.Tags = expandTags(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, *iotdps); err != nil {
		return fmt.Errorf("updating IoT Device Provisioning Service %s: %+v", id, err)
	}

	return resourceIotHubDPSRead(d, meta)
}

func resourceIotHubDPSDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseProvisioningServiceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return waitForIotHubDPSToBeDeleted(ctx, client, *id, d)
}

func waitForIotHubDPSToBeDeleted(ctx context.Context, client *iotdpsresource.IotDpsResourceClient, id commonids.ProvisioningServiceId, d *pluginsdk.ResourceData) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IoT Device Provisioning Service %q to be deleted", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubdpsStateStatusCodeRefreshFunc(ctx, client, id),
		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for IoT Device Provisioning Service %q to be deleted: %+v", id, err)
	}

	return nil
}

func iothubdpsStateStatusCodeRefreshFunc(ctx context.Context, client *iotdpsresource.IotDpsResourceClient, id commonids.ProvisioningServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		log.Printf("Retrieving IoT Device Provisioning Service %q returned Status %q", id, statusCode)

		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, statusCode, nil
			}
			return nil, "", fmt.Errorf("polling for the status of the IoT Device Provisioning Service %q: %+v", id, err)
		}

		return res, statusCode, nil
	}
}

func expandIoTHubDPSSku(d *pluginsdk.ResourceData) iotdpsresource.IotDpsSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})

	skuName := iotdpsresource.IotDpsSku(skuMap["name"].(string))
	return iotdpsresource.IotDpsSkuInfo{
		Name:     &skuName,
		Capacity: utils.Int64(int64(skuMap["capacity"].(int))),
	}
}

func expandIoTHubDPSIoTHubs(input []interface{}) *[]iotdpsresource.IotHubDefinitionDescription {
	linkedHubs := make([]iotdpsresource.IotHubDefinitionDescription, 0)

	for _, attr := range input {
		linkedHubConfig := attr.(map[string]interface{})
		linkedHub := iotdpsresource.IotHubDefinitionDescription{
			ConnectionString:      linkedHubConfig["connection_string"].(string),
			AllocationWeight:      utils.Int64(int64(linkedHubConfig["allocation_weight"].(int))),
			ApplyAllocationPolicy: utils.Bool(linkedHubConfig["apply_allocation_policy"].(bool)),
			Location:              azure.NormalizeLocation(linkedHubConfig["location"].(string)),
		}

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return &linkedHubs
}

func flattenIoTHubDPSSku(input iotdpsresource.IotDpsSkuInfo) []interface{} {
	output := make(map[string]interface{})

	name := ""
	if input.Name != nil {
		name = string(*input.Name)
	}
	output["name"] = name

	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []interface{}{output}
}

func flattenIoTHubDPSLinkedHub(input *[]iotdpsresource.IotHubDefinitionDescription) []interface{} {
	linkedHubs := make([]interface{}, 0)
	if input == nil {
		return linkedHubs
	}

	for _, attr := range *input {
		linkedHub := make(map[string]interface{})

		if attr.Name != nil {
			linkedHub["hostname"] = *attr.Name
		}
		if attr.ApplyAllocationPolicy != nil {
			linkedHub["apply_allocation_policy"] = *attr.ApplyAllocationPolicy
		}
		if attr.AllocationWeight != nil {
			linkedHub["allocation_weight"] = *attr.AllocationWeight
		}

		linkedHub["connection_string"] = attr.ConnectionString
		linkedHub["location"] = azure.NormalizeLocation(attr.Location)

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return linkedHubs
}

func expandDpsIPFilterRules(d *pluginsdk.ResourceData) *[]iotdpsresource.IPFilterRule {
	ipFilterRuleList := d.Get("ip_filter_rule").([]interface{})
	if len(ipFilterRuleList) == 0 {
		return nil
	}

	rules := make([]iotdpsresource.IPFilterRule, 0)

	for _, r := range ipFilterRuleList {
		rawRule := r.(map[string]interface{})
		ipFilterTargetType := iotdpsresource.IPFilterTargetType(rawRule["target"].(string))
		rule := &iotdpsresource.IPFilterRule{
			FilterName: rawRule["name"].(string),
			Action:     iotdpsresource.IPFilterActionType(rawRule["action"].(string)),
			IPMask:     rawRule["ip_mask"].(string),
			Target:     &ipFilterTargetType,
		}

		rules = append(rules, *rule)
	}
	return &rules
}

func flattenDpsIPFilterRules(in *[]iotdpsresource.IPFilterRule) []interface{} {
	rules := make([]interface{}, 0)
	if in == nil {
		return rules
	}

	for _, r := range *in {
		rawRule := make(map[string]interface{})

		rawRule["name"] = r.FilterName
		rawRule["action"] = string(r.Action)
		rawRule["ip_mask"] = r.IPMask

		if r.Target != nil && *r.Target != "" {
			rawRule["target"] = string(*r.Target)
		}

		rules = append(rules, rawRule)
	}
	return rules
}
