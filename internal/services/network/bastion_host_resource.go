// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var skuWeight = map[string]int8{
	"Basic":    1,
	"Standard": 2,
}

func resourceBastionHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBastionHostCreateUpdate,
		Read:   resourceBastionHostRead,
		Update: resourceBastionHostCreateUpdate,
		Delete: resourceBastionHostDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := bastionhosts.ParseBastionHostID(id)
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
				ValidateFunc: validate.BastionHostName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"copy_paste_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"file_copy_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ip_configuration": {
				Type:     pluginsdk.TypeList,
				ForceNew: true,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.BastionIPConfigName,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.BastionSubnetName,
						},
						"public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.PublicIpAddressID,
						},
					},
				},
			},

			"ip_connect_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"scale_units": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 50),
				Default:      2,
			},

			"shareable_link_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(bastionhosts.BastionHostSkuNameBasic),
					string(bastionhosts.BastionHostSkuNameStandard),
				}, false),
				Default: string(bastionhosts.BastionHostSkuNameBasic),
			},

			"tunneling_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"dns_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku", func(ctx context.Context, old, new, meta interface{}) bool {
				// downgrade the SKU is not supported, recreate the resource
				if old.(string) != "" && new.(string) != "" {
					return skuWeight[old.(string)] > skuWeight[new.(string)]
				}
				return false
			}),
		),
	}
}

func resourceBastionHostCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHosts
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for Azure Bastion Host creation.")

	id := bastionhosts.NewBastionHostID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	scaleUnits := d.Get("scale_units").(int)
	sku := d.Get("sku").(string)
	fileCopyEnabled := d.Get("file_copy_enabled").(bool)
	ipConnectEnabled := d.Get("ip_connect_enabled").(bool)
	shareableLinkEnabled := d.Get("shareable_link_enabled").(bool)
	tunnelingEnabled := d.Get("tunneling_enabled").(bool)

	if scaleUnits > 2 && sku == string(bastionhosts.BastionHostSkuNameBasic) {
		return fmt.Errorf("`scale_units` only can be changed when `sku` is `Standard`. `scale_units` is always `2` when `sku` is `Basic`")
	}

	if fileCopyEnabled && sku == string(bastionhosts.BastionHostSkuNameBasic) {
		return fmt.Errorf("`file_copy_enabled` is only supported when `sku` is `Standard`")
	}

	if ipConnectEnabled && sku == string(bastionhosts.BastionHostSkuNameBasic) {
		return fmt.Errorf("`ip_connect_enabled` is only supported when `sku` is `Standard`")
	}

	if shareableLinkEnabled && sku == string(bastionhosts.BastionHostSkuNameBasic) {
		return fmt.Errorf("`shareable_link_enabled` is only supported when `sku` is `Standard`")
	}

	if tunnelingEnabled && sku == string(bastionhosts.BastionHostSkuNameBasic) {
		return fmt.Errorf("`tunneling_enabled` is only supported when `sku` is `Standard`")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_bastion_host", id.ID())
		}
	}

	parameters := bastionhosts.BastionHost{
		Location: &location,
		Properties: &bastionhosts.BastionHostPropertiesFormat{
			DisableCopyPaste:    utils.Bool(!d.Get("copy_paste_enabled").(bool)),
			EnableFileCopy:      utils.Bool(fileCopyEnabled),
			EnableIPConnect:     utils.Bool(ipConnectEnabled),
			EnableShareableLink: utils.Bool(shareableLinkEnabled),
			EnableTunneling:     utils.Bool(tunnelingEnabled),
			IPConfigurations:    expandBastionHostIPConfiguration(d.Get("ip_configuration").([]interface{})),
			ScaleUnits:          utils.Int64(int64(d.Get("scale_units").(int))),
		},
		Sku: &bastionhosts.Sku{
			Name: pointer.To(bastionhosts.BastionHostSkuName(sku)),
		},
		Tags: tags.Expand(t),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBastionHostRead(d, meta)
}

func resourceBastionHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHosts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bastionhosts.ParseBastionHostID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.BastionHostName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(*sku.Name))
		}

		if props := model.Properties; props != nil {
			d.Set("dns_name", props.DnsName)
			d.Set("scale_units", props.ScaleUnits)
			d.Set("file_copy_enabled", props.EnableFileCopy)
			d.Set("ip_connect_enabled", props.EnableIPConnect)
			d.Set("shareable_link_enabled", props.EnableShareableLink)
			d.Set("tunneling_enabled", props.EnableTunneling)

			copyPasteEnabled := true
			if props.DisableCopyPaste != nil {
				copyPasteEnabled = !*props.DisableCopyPaste
			}
			d.Set("copy_paste_enabled", copyPasteEnabled)

			if ipConfigs := props.IPConfigurations; ipConfigs != nil {
				if err := d.Set("ip_configuration", flattenBastionHostIPConfiguration(ipConfigs)); err != nil {
					return fmt.Errorf("flattening `ip_configuration`: %+v", err)
				}
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceBastionHostDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHosts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bastionhosts.ParseBastionHostID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandBastionHostIPConfiguration(input []interface{}) (ipConfigs *[]bastionhosts.BastionHostIPConfiguration) {
	if len(input) == 0 {
		return nil
	}

	property := input[0].(map[string]interface{})
	ipConfName := property["name"].(string)
	subID := property["subnet_id"].(string)
	pipID := property["public_ip_address_id"].(string)

	return &[]bastionhosts.BastionHostIPConfiguration{
		{
			Name: &ipConfName,
			Properties: &bastionhosts.BastionHostIPConfigurationPropertiesFormat{
				Subnet: bastionhosts.SubResource{
					Id: &subID,
				},
				PublicIPAddress: bastionhosts.SubResource{
					Id: &pipID,
				},
			},
		},
	}
}

func flattenBastionHostIPConfiguration(ipConfigs *[]bastionhosts.BastionHostIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		ipConfig := make(map[string]interface{})

		if config.Name != nil {
			ipConfig["name"] = *config.Name
		}

		if props := config.Properties; props != nil {
			subnetId := ""
			if subnet := props.Subnet; subnet.Id != nil {
				subnetId = *subnet.Id
			}
			ipConfig["subnet_id"] = subnetId

			publicIpId := ""
			if pip := props.PublicIPAddress; pip.Id != nil {
				publicIpId = *pip.Id
			}
			ipConfig["public_ip_address_id"] = publicIpId
		}

		result = append(result, ipConfig)
	}
	return result
}
