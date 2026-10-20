// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_topic

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridTopicRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Topics
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topics.ParseTopicID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := client.ListSharedAccessKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Shared Access Keys for %s: %+v", *id, err)
	}

	d.Set("name", id.TopicName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("endpoint", props.Endpoint)
			d.Set("input_schema", pointer.FromEnum(props.InputSchema))

			if err := d.Set("input_mapping_fields", flattenTopicInputMapping(props.InputSchemaMapping)); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			if err := d.Set("input_mapping_default_values", flattenTopicInputMappingDefaultValues(props.InputSchemaMapping)); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == topics.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
			d.Set("minimum_tls_version", pointer.From(props.MinimumTlsVersionAllowed))
			d.Set("data_residency_boundary", pointer.From(props.DataResidencyBoundary))

			if err := d.Set("inbound_ip_rule", flattenTopicInboundIPRules(props.InboundIPRules)); err != nil {
				return fmt.Errorf("setting `inbound_ip_rule`: %+v", err)
			}

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !*props.DisableLocalAuth
			}

			d.Set("local_auth_enabled", localAuthEnabled)
		}

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.Key1)
		d.Set("secondary_access_key", model.Key2)
	}

	return nil
}
