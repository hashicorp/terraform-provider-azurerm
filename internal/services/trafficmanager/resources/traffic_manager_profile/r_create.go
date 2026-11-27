// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmTrafficManagerProfileCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).TrafficManager.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := profiles.NewTrafficManagerProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s", id)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_traffic_manager_profile", id.ID())
		}
	}

	trafficRoutingMethod := profiles.TrafficRoutingMethod(d.Get("traffic_routing_method").(string))
	// No existing profile - start from a new struct.
	profile := profiles.Profile{
		Name:     pointer.To(id.TrafficManagerProfileName),
		Location: pointer.To("global"), // must be provided in request
		Properties: &profiles.ProfileProperties{
			TrafficRoutingMethod:        &trafficRoutingMethod,
			TrafficViewEnrollmentStatus: expandArmTrafficManagerTrafficView(d.Get("traffic_view_enabled").(bool)),
			DnsConfig:                   expandArmTrafficManagerDNSConfig(d),
			MonitorConfig:               expandArmTrafficManagerMonitorConfig(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if maxReturn, ok := d.GetOk("max_return"); ok {
		profile.Properties.MaxReturn = pointer.To(int64(maxReturn.(int)))
	}

	if status, ok := d.GetOk("profile_status"); ok {
		profile.Properties.ProfileStatus = pointer.ToEnum[profiles.ProfileStatus](status.(string))
	}

	trafficRoutingMethodPtr := profile.Properties.TrafficRoutingMethod
	if *trafficRoutingMethodPtr == profiles.TrafficRoutingMethodMultiValue &&
		profile.Properties.MaxReturn == nil {
		return errors.New("`max_return` must be specified when `traffic_routing_method` is set to `MultiValue`")
	}

	if *profile.Properties.MonitorConfig.IntervalInSeconds == int64(10) &&
		*profile.Properties.MonitorConfig.TimeoutInSeconds == int64(10) {
		return errors.New("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`")
	}

	if _, err := client.CreateOrUpdate(ctx, id, profile); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceArmTrafficManagerProfileRead(d, meta)
}
