// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_topic

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridTopicCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Topics
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := topics.NewTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventgrid_topic", id.ID())
		}
	}

	inboundIPRules := expandTopicInboundIPRules(d.Get("inbound_ip_rule").([]any))
	publicNetworkAccess := topics.PublicNetworkAccessDisabled
	if v, ok := d.GetOk("public_network_access_enabled"); ok && v.(bool) {
		publicNetworkAccess = topics.PublicNetworkAccessEnabled
	}

	topic := topics.Topic{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &topics.TopicProperties{
			InputSchemaMapping:  expandTopicInputMapping(d),
			InputSchema:         pointer.ToEnum[topics.InputSchema](d.Get("input_schema").(string)),
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
			InboundIPRules:      inboundIPRules,
			DisableLocalAuth:    pointer.To(!d.Get("local_auth_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("identity"); ok {
		identityRaw := v.([]any)
		identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		topic.Identity = identity
	}

	if tls, ok := d.GetOk("minimum_tls_version"); ok {
		topic.Properties.MinimumTlsVersionAllowed = pointer.ToEnum[topics.TlsVersion](tls.(string))
	}

	if boundary, ok := d.GetOk("data_residency_boundary"); ok {
		topic.Properties.DataResidencyBoundary = pointer.ToEnum[topics.DataResidencyBoundary](boundary.(string))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, topic, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridTopicRead(d, meta)
}

func expandTopicInputMapping(d *pluginsdk.ResourceData) *topics.JsonInputSchemaMapping {
	imf, imfok := d.GetOk("input_mapping_fields")

	imdv, imdvok := d.GetOk("input_mapping_default_values")

	if !imfok && !imdvok {
		return nil
	}

	jismp := topics.JsonInputSchemaMappingProperties{}

	if imfok {
		mappings := imf.([]any)
		if len(mappings) > 0 && mappings[0] != nil {
			if mapping := mappings[0].(map[string]any); mapping != nil {
				if id := mapping["id"].(string); id != "" {
					jismp.Id = &topics.JsonField{
						SourceField: &id,
					}
				}

				if eventTime := mapping["event_time"].(string); eventTime != "" {
					jismp.EventTime = &topics.JsonField{
						SourceField: &eventTime,
					}
				}

				if topic := mapping["topic"].(string); topic != "" {
					jismp.Topic = &topics.JsonField{
						SourceField: &topic,
					}
				}

				if dataVersion := mapping["data_version"].(string); dataVersion != "" {
					jismp.DataVersion = &topics.JsonFieldWithDefault{
						SourceField: &dataVersion,
					}
				}

				if subject := mapping["subject"].(string); subject != "" {
					jismp.Subject = &topics.JsonFieldWithDefault{
						SourceField: &subject,
					}
				}

				if eventType := mapping["event_type"].(string); eventType != "" {
					jismp.EventType = &topics.JsonFieldWithDefault{
						SourceField: &eventType,
					}
				}
			}
		}
	}

	if imdvok {
		mappings := imdv.([]any)
		if len(mappings) > 0 && mappings[0] != nil {
			if mapping := mappings[0].(map[string]any); mapping != nil {
				if dataVersion := mapping["data_version"].(string); dataVersion != "" {
					if v := jismp.DataVersion; v != nil && v.SourceField != nil {
						jismp.DataVersion = &topics.JsonFieldWithDefault{
							SourceField:  v.SourceField,
							DefaultValue: &dataVersion,
						}
					} else {
						jismp.DataVersion = &topics.JsonFieldWithDefault{
							DefaultValue: &dataVersion,
						}
					}
				}

				if subject := mapping["subject"].(string); subject != "" {
					if v := jismp.Subject; v != nil && v.SourceField != nil {
						jismp.Subject = &topics.JsonFieldWithDefault{
							DefaultValue: &subject,
							SourceField:  v.SourceField,
						}
					} else {
						jismp.Subject = &topics.JsonFieldWithDefault{
							DefaultValue: &subject,
						}
					}
				}

				if eventType := mapping["event_type"].(string); eventType != "" {
					if v := jismp.EventType; v != nil && v.SourceField != nil {
						jismp.EventType = &topics.JsonFieldWithDefault{
							DefaultValue: &eventType,
							SourceField:  v.SourceField,
						}
					} else {
						jismp.EventType = &topics.JsonFieldWithDefault{
							DefaultValue: &eventType,
						}
					}
				}
			}
		}
	}

	return &topics.JsonInputSchemaMapping{
		Properties: &jismp,
	}
}

func expandTopicInboundIPRules(input []any) *[]topics.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]topics.InboundIPRule, 0)
	for _, item := range input {
		rawRule := item.(map[string]any)
		rules = append(rules, topics.InboundIPRule{
			Action: pointer.ToEnum[topics.IPActionType](rawRule["action"].(string)),
			IPMask: pointer.To(rawRule["ip_mask"].(string)),
		})
	}
	return &rules
}
