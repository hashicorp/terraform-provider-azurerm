// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmCdnEndpointCustomDomainUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customdomains.ParseEndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	const (
		turnOn = iota
		turnOff
		update
		noChange
	)

	var (
		cdnManagedHTTPSStatus = noChange
		cdnManagedHTTPSParams cdn.BasicCustomDomainHTTPSParameters

		userManagedHTTPSStatus = noChange
		userManagedHTTPSParams cdn.BasicCustomDomainHTTPSParameters
	)

	if d.HasChange("cdn_managed_https") {
		props := resp.CustomDomainProperties
		if props == nil {
			return fmt.Errorf("unexpected nil of `CustomDomainProperties` in response")
		}

		cdnManagedHTTPSParams = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https").([]any))

		if props.CustomHTTPSParameters == nil {
			cdnManagedHTTPSStatus = turnOn
		} else {
			if cdnManagedHTTPSParams == nil {
				cdnManagedHTTPSStatus = turnOff
			} else {
				cdnManagedHTTPSStatus = update
			}
		}
	}

	if d.HasChange("user_managed_https") {
		props := resp.CustomDomainProperties
		if props == nil {
			return fmt.Errorf("unexpected nil of `CustomDomainProperties` in response")
		}

		var err error
		userManagedHTTPSParams, err = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(ctx, d.Get("user_managed_https").([]any), meta.(*clients.Client))
		if err != nil {
			return err
		}

		if props.CustomHTTPSParameters == nil {
			userManagedHTTPSStatus = turnOn
		} else {
			if userManagedHTTPSParams == nil {
				userManagedHTTPSStatus = turnOff
			} else {
				userManagedHTTPSStatus = update
			}
		}
	}

	// There are theoretically 16 (4x4) combinations of the cdn/user managed https status combinations.
	// While actually there are only following 8 combinations due to the exclusive nature of both settings.
	// +-----------------------------------+
	// |     	| n/a | on | off | update  |
	// |--------|--------------------------|
	// | n/a 	|     |  x |  x  |    x    |
	// | on  	|  x  |    |  x  |         |
	// | off    |  x  |  x |     |         |
	// | update |  x  |    |     |         |
	// +-----------------------------------+

	switch {
	case cdnManagedHTTPSStatus == turnOff || cdnManagedHTTPSStatus == update:
		if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
			return fmt.Errorf("disable CDN Managed HTTPS on %q: %+v", *id, err)
		}
	case userManagedHTTPSStatus == turnOff || userManagedHTTPSStatus == update:
		if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
			return fmt.Errorf("disable User Managed HTTPS on %q: %+v", *id, err)
		}
	}

	switch {
	case cdnManagedHTTPSStatus == turnOn || cdnManagedHTTPSStatus == update:
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id, cdnManagedHTTPSParams); err != nil {
			return fmt.Errorf("enable CDN Managed HTTPS on %q: %+v", *id, err)
		}
	case userManagedHTTPSStatus == turnOn || userManagedHTTPSStatus == update:
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id, userManagedHTTPSParams); err != nil {
			return fmt.Errorf("enable User Managed HTTPS on %q: %+v", *id, err)
		}
	}

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}
