// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Domains
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domains.ParseDomainID(d.Id())
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

	keys, err := client.ListSharedAccessKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Shared Access Keys for %s: %+v", *id, err)
	}

	d.Set("name", id.DomainName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("endpoint", props.Endpoint)

			inputSchema := ""
			if props.InputSchema != nil {
				inputSchema = string(*props.InputSchema)
			}
			d.Set("input_schema", inputSchema)

			if err := d.Set("input_mapping_fields", flattenDomainInputMapping(props.InputSchemaMapping)); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			if err := d.Set("input_mapping_default_values", flattenDomainInputMappingDefaultValues(props.InputSchemaMapping)); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == domains.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			if err := d.Set("inbound_ip_rule", flattenDomainInboundIPRules(props.InboundIPRules)); err != nil {
				return fmt.Errorf("setting `inbound_ip_rule`: %+v", err)
			}

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_auth_enabled", localAuthEnabled)

			autoCreateTopicWithFirstSubscription := true
			if props.AutoCreateTopicWithFirstSubscription != nil {
				autoCreateTopicWithFirstSubscription = *props.AutoCreateTopicWithFirstSubscription
			}
			d.Set("auto_create_topic_with_first_subscription", autoCreateTopicWithFirstSubscription)

			d.Set("data_residency_boundary", pointer.From(props.DataResidencyBoundary))
			d.Set("minimum_tls_version", pointer.From(props.MinimumTlsVersionAllowed))

			autoDeleteTopicWithLastSubscription := true
			if props.AutoDeleteTopicWithLastSubscription != nil {
				autoDeleteTopicWithLastSubscription = *props.AutoDeleteTopicWithLastSubscription
			}
			d.Set("auto_delete_topic_with_last_subscription", autoDeleteTopicWithLastSubscription)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.Key1)
		d.Set("secondary_access_key", model.Key2)
	}

	return nil
}
