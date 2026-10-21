// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParsePublicIPAddressID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
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

	if d.HasChange("allocation_method") {
		payload.Properties.PublicIPAllocationMethod = pointer.ToEnum[publicipaddresses.IPAllocationMethod](d.Get("allocation_method").(string))
	}

	if d.HasChange("ddos_protection_mode") {
		if payload.Properties.DdosSettings == nil {
			payload.Properties.DdosSettings = &publicipaddresses.DdosSettings{}
		}
		payload.Properties.DdosSettings.ProtectionMode = pointer.ToEnum[publicipaddresses.DdosSettingsProtectionMode](d.Get("ddos_protection_mode").(string))
	}

	if d.HasChange("ddos_protection_plan_id") {
		if !strings.EqualFold(string(*payload.Properties.DdosSettings.ProtectionMode), "enabled") {
			return fmt.Errorf("ddos protection plan id can only be set when ddos protection is enabled")
		}
		if payload.Properties.DdosSettings == nil {
			payload.Properties.DdosSettings = &publicipaddresses.DdosSettings{}
		}
		payload.Properties.DdosSettings.DdosProtectionPlan = &publicipaddresses.SubResource{
			Id: pointer.To(d.Get("ddos_protection_plan_id").(string)),
		}
	}

	if d.HasChange("idle_timeout_in_minutes") {
		payload.Properties.IdleTimeoutInMinutes = pointer.To(int64(d.Get("idle_timeout_in_minutes").(int)))
	}

	if d.HasChanges("domain_name_label", "domain_name_label_scope", "reverse_fqdn") {
		dnl, dnlOk := d.GetOk("domain_name_label")
		rfqdn, rfqdnOk := d.GetOk("reverse_fqdn")
		dnlc, dnlcOk := d.GetOk("domain_name_label_scope")

		if dnlOk || rfqdnOk || dnlcOk {
			dnsSettings := publicipaddresses.PublicIPAddressDnsSettings{}

			if rfqdnOk {
				dnsSettings.ReverseFqdn = pointer.To(rfqdn.(string))
			}

			if dnlOk {
				dnsSettings.DomainNameLabel = pointer.To(dnl.(string))
			}

			if dnlcOk {
				dnsSettings.DomainNameLabelScope = pointer.ToEnum[publicipaddresses.PublicIPAddressDnsSettingsDomainNameLabelScope](dnlc.(string))
			}

			payload.Properties.DnsSettings = &dnsSettings
		} else {
			payload.Properties.DnsSettings = nil
		}
	}

	if d.HasChanges("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourcePublicIpRead(d, meta)
}
