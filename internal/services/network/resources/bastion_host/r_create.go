// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bastion_host

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/bastionhosts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBastionHostCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &bastionhosts.BastionHostPropertiesFormat{
			IPConfigurations: expandBastionHostIPConfiguration(d.Get("ip_configuration").([]any)),
			ScaleUnits:       pointer.To(int64(d.Get("scale_units").(int))),
		},
		Sku: &bastionhosts.Sku{
			Name: pointer.To(sku),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
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

	if sku == bastionhosts.BastionHostSkuNamePremium {
		_, ok := d.GetOk("ip_configuration.0.public_ip_address_id")
		parameters.Properties.EnablePrivateOnlyBastion = pointer.To(!ok)
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = pointer.To(zones)
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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceBastionHostRead(d, meta)
}

func expandBastionHostIPConfiguration(input []any) (ipConfigs *[]bastionhosts.BastionHostIPConfiguration) {
	if len(input) == 0 {
		return nil
	}

	property := input[0].(map[string]any)

	ipConfig := bastionhosts.BastionHostIPConfiguration{
		Name: pointer.To(property["name"].(string)),
		Properties: &bastionhosts.BastionHostIPConfigurationPropertiesFormat{
			Subnet: bastionhosts.SubResource{
				Id: pointer.To(property["subnet_id"].(string)),
			},
		},
	}

	if pipID, ok := property["public_ip_address_id"].(string); ok && pipID != "" {
		ipConfig.Properties.PublicIPAddress = &bastionhosts.SubResource{
			Id: &pipID,
		}
	}

	return &[]bastionhosts.BastionHostIPConfiguration{ipConfig}
}
