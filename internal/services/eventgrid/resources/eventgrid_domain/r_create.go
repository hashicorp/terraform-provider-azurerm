// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Domains
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := domains.NewDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventgrid_domain", id.ID())
		}
	}

	inboundIPRules := expandDomainInboundIPRules(d.Get("inbound_ip_rule").([]any))
	publicNetworkAccess := domains.PublicNetworkAccessDisabled
	if v, ok := d.GetOk("public_network_access_enabled"); ok && v.(bool) {
		publicNetworkAccess = domains.PublicNetworkAccessEnabled
	}

	domain := domains.Domain{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &domains.DomainProperties{
			AutoCreateTopicWithFirstSubscription: pointer.To(d.Get("auto_create_topic_with_first_subscription").(bool)),
			AutoDeleteTopicWithLastSubscription:  pointer.To(d.Get("auto_delete_topic_with_last_subscription").(bool)),
			DisableLocalAuth:                     pointer.To(!d.Get("local_auth_enabled").(bool)),
			InboundIPRules:                       inboundIPRules,
			InputSchema:                          pointer.ToEnum[domains.InputSchema](d.Get("input_schema").(string)),
			InputSchemaMapping:                   expandDomainInputMapping(d),
			PublicNetworkAccess:                  pointer.To(publicNetworkAccess),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("identity"); ok {
		identity, err := identity.ExpandSystemAndUserAssignedMap(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		domain.Identity = identity
	}

	if tls, ok := d.GetOk("minimum_tls_version"); ok {
		domain.Properties.MinimumTlsVersionAllowed = pointer.ToEnum[domains.TlsVersion](tls.(string))
	}

	if boundary, ok := d.GetOk("data_residency_boundary"); ok {
		domain.Properties.DataResidencyBoundary = pointer.ToEnum[domains.DataResidencyBoundary](boundary.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, domain, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridDomainRead(d, meta)
}

func expandDomainInputMapping(d *pluginsdk.ResourceData) *domains.JsonInputSchemaMapping {
	imf, imfok := d.GetOk("input_mapping_fields")

	imdv, imdvok := d.GetOk("input_mapping_default_values")

	if !imfok && !imdvok {
		return nil
	}

	jismp := domains.JsonInputSchemaMappingProperties{}

	if imfok {
		mappings := imf.([]any)
		if len(mappings) > 0 && mappings[0] != nil {
			mapping := mappings[0].(map[string]any)

			if id := mapping["id"].(string); id != "" {
				jismp.Id = &domains.JsonField{SourceField: &id}
			}

			if eventTime := mapping["event_time"].(string); eventTime != "" {
				jismp.EventTime = &domains.JsonField{
					SourceField: &eventTime,
				}
			}

			if topic := mapping["topic"].(string); topic != "" {
				jismp.Topic = &domains.JsonField{
					SourceField: &topic,
				}
			}

			if dataVersion := mapping["data_version"].(string); dataVersion != "" {
				jismp.DataVersion = &domains.JsonFieldWithDefault{
					SourceField: &dataVersion,
				}
			}

			if subject := mapping["subject"].(string); subject != "" {
				jismp.Subject = &domains.JsonFieldWithDefault{
					SourceField: &subject,
				}
			}

			if eventType := mapping["event_type"].(string); eventType != "" {
				jismp.EventType = &domains.JsonFieldWithDefault{
					SourceField: &eventType,
				}
			}
		}
	}

	if imdvok {
		mappings := imdv.([]any)
		if len(mappings) > 0 && mappings[0] != nil {
			mapping := mappings[0].(map[string]any)

			if dataVersion := mapping["data_version"].(string); dataVersion != "" {
				jismp.DataVersion = &domains.JsonFieldWithDefault{
					DefaultValue: &dataVersion,
				}
			}

			if subject := mapping["subject"].(string); subject != "" {
				jismp.Subject = &domains.JsonFieldWithDefault{
					DefaultValue: &subject,
				}
			}

			if eventType := mapping["event_type"].(string); eventType != "" {
				jismp.EventType = &domains.JsonFieldWithDefault{
					DefaultValue: &eventType,
				}
			}
		}
	}

	return &domains.JsonInputSchemaMapping{
		Properties: &jismp,
	}
}

func expandDomainInboundIPRules(input []any) *[]domains.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]domains.InboundIPRule, 0)
	for _, item := range input {
		rawRule := item.(map[string]any)
		rules = append(rules, domains.InboundIPRule{
			Action: pointer.ToEnum[domains.IPActionType](rawRule["action"].(string)),
			IPMask: pointer.To(rawRule["ip_mask"].(string)),
		})
	}
	return &rules
}
