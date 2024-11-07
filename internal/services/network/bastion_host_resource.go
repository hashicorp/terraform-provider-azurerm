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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-01-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var skuWeight = map[string]int8{
	"Developer": 1,
	"Basic":     2,
	"Standard":  3,
	"Premium":   4,
}

func resourceBastionHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBastionHostCreate,
		Read:   resourceBastionHostRead,
		Update: resourceBastionHostUpdate,
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
				Optional: true,
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
							ValidateFunc: commonids.ValidatePublicIPAddressID,
						},
					},
				},
			},

			"ip_connect_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kerberos_enabled": {
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(bastionhosts.PossibleValuesForBastionHostSkuName(), false),
				Default:      string(bastionhosts.BastionHostSkuNameBasic),
			},

			"tunneling_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"session_recording_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
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

func resourceBastionHostCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for Azure Bastion Host creation.")

	id := bastionhosts.NewBastionHostID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	scaleUnits := d.Get("scale_units").(int)
	sku := bastionhosts.BastionHostSkuName(d.Get("sku").(string))
	fileCopyEnabled := d.Get("file_copy_enabled").(bool)
	ipConnectEnabled := d.Get("ip_connect_enabled").(bool)
	kerberosEnabled := d.Get("kerberos_enabled").(bool)
	shareableLinkEnabled := d.Get("shareable_link_enabled").(bool)
	tunnelingEnabled := d.Get("tunneling_enabled").(bool)
	sessionRecordingEnabled := d.Get("session_recording_enabled").(bool)

	if scaleUnits > 2 && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`scale_units` only can be changed when `sku` is `Standard` or `Premium`. `scale_units` is always `2` when `sku` is `Basic`")
	}

	if fileCopyEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`file_copy_enabled` is only supported when `sku` is `Standard` or `Premium`")
	}

	if ipConnectEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`ip_connect_enabled` is only supported when `sku` is `Standard` or `Premium`")
	}

	if kerberosEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`kerberos_enabled` is only supported when `sku` is `Standard` or `Premium`")
	}

	if shareableLinkEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`shareable_link_enabled` is only supported when `sku` is `Standard` or `Premium`")
	}

	if tunnelingEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
		return fmt.Errorf("`tunneling_enabled` is only supported when `sku` is `Standard` or `Premium`")
	}

	if sessionRecordingEnabled && sku != bastionhosts.BastionHostSkuNamePremium {
		return fmt.Errorf("`session_recording_enabled` is only supported when `sku` is `Premium`")
	}

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_bastion_host", id.ID())
	}

	parameters := bastionhosts.BastionHost{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &bastionhosts.BastionHostPropertiesFormat{
			IPConfigurations: expandBastionHostIPConfiguration(d.Get("ip_configuration").([]interface{})),
			ScaleUnits:       pointer.To(int64(d.Get("scale_units").(int))),
		},
		Sku: &bastionhosts.Sku{
			Name: pointer.To(sku),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v := !d.Get("copy_paste_enabled").(bool); v {
		parameters.Properties.DisableCopyPaste = pointer.To(v)
	}

	if fileCopyEnabled {
		parameters.Properties.EnableFileCopy = pointer.To(fileCopyEnabled)
	}

	if ipConnectEnabled {
		parameters.Properties.EnableIPConnect = pointer.To(ipConnectEnabled)
	}

	if kerberosEnabled {
		parameters.Properties.EnableKerberos = pointer.To(kerberosEnabled)
	}

	if shareableLinkEnabled {
		parameters.Properties.EnableShareableLink = pointer.To(shareableLinkEnabled)
	}

	if tunnelingEnabled {
		parameters.Properties.EnableTunneling = pointer.To(tunnelingEnabled)
	}

	if sessionRecordingEnabled {
		parameters.Properties.EnableSessionRecording = pointer.To(sessionRecordingEnabled)
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		if sku != bastionhosts.BastionHostSkuNameDeveloper {
			return fmt.Errorf("`virtual_network_id` is only supported when `sku` is `Developer`")
		}

		parameters.Properties.VirtualNetwork = &bastionhosts.SubResource{
			Id: pointer.To(v.(string)),
		}
	} else if sku == bastionhosts.BastionHostSkuNameDeveloper {
		return fmt.Errorf("`virtual_network_id` is required when `sku` is `Developer`")
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBastionHostRead(d, meta)
}

func resourceBastionHostUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bastionhosts.ParseBastionHostID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
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

	sku := bastionhosts.BastionHostSkuName(d.Get("sku").(string))

	if d.HasChange("sku") {
		payload.Sku = &bastionhosts.Sku{
			Name: pointer.To(sku),
		}
	}

	if d.HasChange("copy_paste_enabled") {
		payload.Properties.DisableCopyPaste = pointer.To(!d.Get("copy_paste_enabled").(bool))
	}

	if d.HasChange("file_copy_enabled") {
		fileCopyEnabled := d.Get("file_copy_enabled").(bool)
		if fileCopyEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`file_copy_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableFileCopy = pointer.To(fileCopyEnabled)
	}

	if d.HasChange("ip_connect_enabled") {
		ipConnectEnabled := d.Get("ip_connect_enabled").(bool)
		if ipConnectEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`ip_connect_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableIPConnect = pointer.To(ipConnectEnabled)
	}

	if d.HasChange("scale_units") {
		scaleUnits := d.Get("scale_units").(int)
		if scaleUnits > 2 && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`scale_units` only can be changed when `sku` is `Standard` or `Premium`. `scale_units` is always `2` when `sku` is `Basic`")
		}
		payload.Properties.ScaleUnits = pointer.To(int64(scaleUnits))
	}

	if d.HasChange("shareable_link_enabled") {
		shareableLinkEnabled := d.Get("shareable_link_enabled").(bool)
		if shareableLinkEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`shareable_link_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableShareableLink = pointer.To(shareableLinkEnabled)
	}

	if d.HasChange("tunneling_enabled") {
		tunnelingEnabled := d.Get("tunneling_enabled").(bool)
		if tunnelingEnabled && (sku != bastionhosts.BastionHostSkuNameStandard && sku != bastionhosts.BastionHostSkuNamePremium) {
			return fmt.Errorf("`tunneling_enabled` is only supported when `sku` is `Standard` or `Premium`")
		}
		payload.Properties.EnableTunneling = pointer.To(tunnelingEnabled)
	}

	if d.HasChange("session_recording_enabled") {
		sessionRecordingEnabled := d.Get("session_recording_enabled").(bool)
		if sessionRecordingEnabled && sku != bastionhosts.BastionHostSkuNamePremium {
			return fmt.Errorf("`session_recording_enabled` is only supported when `sku` is `Premium`")
		}
		payload.Properties.EnableSessionRecording = pointer.To(sessionRecordingEnabled)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))

	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceBastionHostRead(d, meta)
}

func resourceBastionHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
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
		d.Set("location", location.NormalizeNilable(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(*sku.Name))
		}

		if props := model.Properties; props != nil {
			d.Set("dns_name", props.DnsName)
			d.Set("scale_units", props.ScaleUnits)
			d.Set("file_copy_enabled", props.EnableFileCopy)
			d.Set("ip_connect_enabled", props.EnableIPConnect)
			d.Set("kerberos_enabled", props.EnableKerberos)
			d.Set("shareable_link_enabled", props.EnableShareableLink)
			d.Set("tunneling_enabled", props.EnableTunneling)
			d.Set("session_recording_enabled", props.EnableSessionRecording)

			virtualNetworkId := ""
			if vnet := props.VirtualNetwork; vnet != nil {
				vnetId, err := commonids.ParseVirtualNetworkID(pointer.From(vnet.Id))
				if err != nil {
					return err
				}
				virtualNetworkId = vnetId.ID()
			}
			d.Set("virtual_network_id", virtualNetworkId)

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
	client := meta.(*clients.Client).Network.BastionHostsClient
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
