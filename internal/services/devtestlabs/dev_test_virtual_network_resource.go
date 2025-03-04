// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmDevTestVirtualNetwork() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmDevTestVirtualNetworkCreate,
		Read:   resourceArmDevTestVirtualNetworkRead,
		Update: resourceArmDevTestVirtualNetworkUpdate,
		Delete: resourceArmDevTestVirtualNetworkDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworks.ParseVirtualNetworkID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestVirtualNetworkUpgradeV0ToV1{},
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
				ValidateFunc: ValidateDevTestVirtualNetworkName(),
			},

			"lab_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"subnet": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// NOTE: O+C We generate the name and resource ID for this block so this should remain
				Computed: true,
				// whilst the API accepts multiple, in practice only one is usable
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"use_in_virtual_machine_creation": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(virtualnetworks.UsagePermissionTypeAllow),
							ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
						},

						"use_public_ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(virtualnetworks.UsagePermissionTypeAllow),
							ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
						},

						"shared_public_ip_address": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_ports": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"backend_port": {
													Type:     pluginsdk.TypeInt,
													Optional: true,
												},

												"transport_protocol": {
													Type:     pluginsdk.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(virtualnetworks.TransportProtocolTcp),
														string(virtualnetworks.TransportProtocolUdp),
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"unique_identifier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDevTestVirtualNetworkCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Virtual Network creation")

	id := virtualnetworks.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, virtualnetworks.GetOperationOptions{})
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_dev_test_virtual_network", id.ID())
	}

	description := d.Get("description").(string)
	subnetsRaw := d.Get("subnet").([]interface{})
	subnets := expandDevTestVirtualNetworkSubnets(subnetsRaw, subscriptionId, id.ResourceGroupName, id.VirtualNetworkName)

	parameters := virtualnetworks.VirtualNetwork{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &virtualnetworks.VirtualNetworkProperties{
			Description:     pointer.To(description),
			SubnetOverrides: subnets,
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkUpdate(d, meta)
}

func resourceArmDevTestVirtualNetworkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworks.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, virtualnetworks.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			flattenedSubnets := flattenDevTestVirtualNetworkSubnets(props.SubnetOverrides)
			if err := d.Set("subnet", flattenedSubnets); err != nil {
				return fmt.Errorf("setting `subnet`: %+v", err)
			}

			// Computed fields
			d.Set("unique_identifier", props.UniqueIdentifier)
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func resourceArmDevTestVirtualNetworkUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Virtual Network creation")
	id, err := virtualnetworks.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("subnet") {
		subnets := expandDevTestVirtualNetworkSubnets(d.Get("subnet").([]interface{}), subscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
		payload.Properties.SubnetOverrides = subnets
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	err = client.CreateOrUpdateThenPoll(ctx, *id, *payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkRead(d, meta)
}

func resourceArmDevTestVirtualNetworkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworks.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, virtualnetworks.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}

func ValidateDevTestVirtualNetworkName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Virtual Network Name can only include alphanumeric characters, underscores, hyphens.")
}

func expandDevTestVirtualNetworkSubnets(input []interface{}, subscriptionId, resourceGroupName, virtualNetworkName string) *[]virtualnetworks.SubnetOverride {
	results := make([]virtualnetworks.SubnetOverride, 0)
	// default found from the Portal
	name := fmt.Sprintf("%sSubnet", virtualNetworkName)
	subnetId := commonids.NewSubnetID(subscriptionId, resourceGroupName, virtualNetworkName, name)
	allow := virtualnetworks.UsagePermissionTypeAllow

	if len(input) == 0 {
		result := virtualnetworks.SubnetOverride{
			ResourceId:                   pointer.To(subnetId.ID()),
			LabSubnetName:                pointer.To(name),
			UsePublicIPAddressPermission: &allow,
			UseInVMCreationPermission:    &allow,
		}
		results = append(results, result)
		return &results
	}

	for _, val := range input {
		v := val.(map[string]interface{})

		subnet := virtualnetworks.SubnetOverride{
			ResourceId:                         pointer.To(subnetId.ID()),
			LabSubnetName:                      pointer.To(name),
			UsePublicIPAddressPermission:       pointer.To(virtualnetworks.UsagePermissionType(v["use_public_ip_address"].(string))),
			UseInVMCreationPermission:          pointer.To(virtualnetworks.UsagePermissionType(v["use_in_virtual_machine_creation"].(string))),
			SharedPublicIPAddressConfiguration: expandDevTestVirtualNetworkSubnetIpAddressConfiguration(v["shared_public_ip_address"].([]interface{})),
		}
		results = append(results, subnet)
	}

	return &results
}

func expandDevTestVirtualNetworkSubnetIpAddressConfiguration(input []interface{}) *virtualnetworks.SubnetSharedPublicIPAddressConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &virtualnetworks.SubnetSharedPublicIPAddressConfiguration{
		AllowedPorts: expandDevTestVirtualNetworkSubnetAllowedPorts(v["allowed_ports"].([]interface{})),
	}
}

func expandDevTestVirtualNetworkSubnetAllowedPorts(input []interface{}) *[]virtualnetworks.Port {
	results := make([]virtualnetworks.Port, 0)

	for _, val := range input {
		v := val.(map[string]interface{})

		allowedPort := virtualnetworks.Port{
			BackendPort:       pointer.To(int64(v["backend_port"].(int))),
			TransportProtocol: pointer.To(virtualnetworks.TransportProtocol(v["transport_protocol"].(string))),
		}
		results = append(results, allowedPort)
	}

	return &results
}

func flattenDevTestVirtualNetworkSubnets(input *[]virtualnetworks.SubnetOverride) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]interface{})
		if v.LabSubnetName != nil {
			output["name"] = *v.LabSubnetName
		}
		output["use_public_ip_address"] = v.UsePublicIPAddressPermission
		output["use_in_virtual_machine_creation"] = v.UseInVMCreationPermission
		output["shared_public_ip_address"] = flattenDevTestVirtualNetworkSubnetIpAddressConfiguration(v.SharedPublicIPAddressConfiguration)

		outputs = append(outputs, output)
	}

	return outputs
}

func flattenDevTestVirtualNetworkSubnetIpAddressConfiguration(input *virtualnetworks.SubnetSharedPublicIPAddressConfiguration) []interface{} {
	outputs := make([]interface{}, 0)

	if input == nil {
		return outputs
	}

	output := make(map[string]interface{})
	if input.AllowedPorts != nil {
		output["allowed_ports"] = flattenDevTestVirtualNetworkSubnetAllowedPorts(input.AllowedPorts)
	}
	outputs = append(outputs, output)
	return outputs
}

func flattenDevTestVirtualNetworkSubnetAllowedPorts(input *[]virtualnetworks.Port) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]interface{})
		output["backend_port"] = pointer.From(v.BackendPort)
		output["transport_protocol"] = pointer.From(v.TransportProtocol)
		outputs = append(outputs, output)
	}

	return outputs
}
