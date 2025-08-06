// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceEventGridTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridTopicCreate,
		Read:   resourceEventGridTopicRead,
		Update: resourceEventGridTopicUpdate,
		Delete: resourceEventGridTopicDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := topics.ParseTopicID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
						"EventGrid topic name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"input_schema": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(topics.InputSchemaEventGridSchema),
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(topics.PossibleValuesForInputSchema(), false),
			},

			// lintignore:XS003
			"input_mapping_fields": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"topic": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"event_time": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"event_type": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"subject": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"data_version": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},

			// lintignore:XS003
			"input_mapping_default_values": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"event_type": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"subject": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"data_version": {
							Type:     pluginsdk.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"inbound_ip_rule": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				MaxItems:   128,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_mask": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(topics.IPActionTypeAllow),
							ValidateFunc: validation.StringInSlice([]string{
								string(topics.IPActionTypeAllow),
							}, false),
						},
					},
				},
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceEventGridTopicCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Topics
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := topics.NewTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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

	inboundIPRules := expandTopicInboundIPRules(d.Get("inbound_ip_rule").([]interface{}))
	publicNetworkAccess := topics.PublicNetworkAccessDisabled
	if v, ok := d.GetOk("public_network_access_enabled"); ok && v.(bool) {
		publicNetworkAccess = topics.PublicNetworkAccessEnabled
	}

	topic := topics.Topic{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &topics.TopicProperties{
			InputSchemaMapping:  expandTopicInputMapping(d),
			InputSchema:         pointer.To(topics.InputSchema(d.Get("input_schema").(string))),
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
			InboundIPRules:      inboundIPRules,
			DisableLocalAuth:    utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("identity"); ok {
		identityRaw := v.([]interface{})
		identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		topic.Identity = identity
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, topic); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridTopicRead(d, meta)
}

func resourceEventGridTopicUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Topics
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topics.ParseTopicID(d.Id())
	if err != nil {
		return err
	}

	payload := topics.TopicUpdateParameters{Properties: &topics.TopicUpdateParameterProperties{}}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := topics.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = topics.PublicNetworkAccessEnabled
		}

		payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("local_auth_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
	}

	if d.HasChange("inbound_ip_rule") {
		inboundIpRule := d.Get("inbound_ip_rule").([]interface{})

		if len(inboundIpRule) == 0 {
			payload.Properties.InboundIPRules = pointer.To([]topics.InboundIPRule{})
		} else {
			payload.Properties.InboundIPRules = expandTopicInboundIPRules(inboundIpRule)
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceEventGridTopicRead(d, meta)
}

func resourceEventGridTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("input_schema", string(pointer.From(props.InputSchema)))

			inputMappingFields := flattenTopicInputMapping(props.InputSchemaMapping)
			if err := d.Set("input_mapping_fields", inputMappingFields); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			inputMappingDefaultValues := flattenTopicInputMappingDefaultValues(props.InputSchemaMapping)
			if err := d.Set("input_mapping_default_values", inputMappingDefaultValues); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == topics.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			inboundIPRules := flattenTopicInboundIPRules(props.InboundIPRules)
			if err := d.Set("inbound_ip_rule", inboundIPRules); err != nil {
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

func resourceEventGridTopicDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Topics
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topics.ParseTopicID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandTopicInputMapping(d *pluginsdk.ResourceData) *topics.JsonInputSchemaMapping {
	imf, imfok := d.GetOk("input_mapping_fields")

	imdv, imdvok := d.GetOk("input_mapping_default_values")

	if !imfok && !imdvok {
		return nil
	}

	jismp := topics.JsonInputSchemaMappingProperties{}

	if imfok {
		mappings := imf.([]interface{})
		if len(mappings) > 0 && mappings[0] != nil {
			if mapping := mappings[0].(map[string]interface{}); mapping != nil {
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
		mappings := imdv.([]interface{})
		if len(mappings) > 0 && mappings[0] != nil {
			if mapping := mappings[0].(map[string]interface{}); mapping != nil {
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

func flattenTopicInputMapping(input topics.InputSchemaMapping) []interface{} {
	val, ok := input.(topics.JsonInputSchemaMapping)
	if !ok {
		return []interface{}{}
	}

	dataVersion := ""
	eventTime := ""
	eventType := ""
	id := ""
	subject := ""
	topic := ""
	if props := val.Properties; props != nil {
		if props.EventTime != nil && props.EventTime.SourceField != nil {
			eventTime = *props.EventTime.SourceField
		}

		if props.Id != nil && props.Id.SourceField != nil {
			id = *props.Id.SourceField
		}

		if props.Topic != nil && props.Topic.SourceField != nil {
			topic = *props.Topic.SourceField
		}

		if props.DataVersion != nil && props.DataVersion.SourceField != nil {
			dataVersion = *props.DataVersion.SourceField
		}

		if props.EventType != nil && props.EventType.SourceField != nil {
			eventType = *props.EventType.SourceField
		}

		if props.Subject != nil && props.Subject.SourceField != nil {
			subject = *props.Subject.SourceField
		}
	}

	return []interface{}{
		map[string]interface{}{
			"data_version": dataVersion,
			"event_time":   eventTime,
			"event_type":   eventType,
			"id":           id,
			"subject":      subject,
			"topic":        topic,
		},
	}
}

func flattenTopicInputMappingDefaultValues(input topics.InputSchemaMapping) []interface{} {
	val, ok := input.(topics.JsonInputSchemaMapping)
	if !ok || val.Properties == nil {
		return []interface{}{}
	}

	dataVersion := ""
	eventType := ""
	subject := ""
	if val.Properties != nil {
		if val.Properties.DataVersion != nil && val.Properties.DataVersion.DefaultValue != nil {
			dataVersion = *val.Properties.DataVersion.DefaultValue
		}
		if val.Properties.EventType != nil && val.Properties.EventType.DefaultValue != nil {
			eventType = *val.Properties.EventType.DefaultValue
		}
		if val.Properties.Subject != nil && val.Properties.Subject.DefaultValue != nil {
			subject = *val.Properties.Subject.DefaultValue
		}
	}

	return []interface{}{
		map[string]interface{}{
			"data_version": dataVersion,
			"event_type":   eventType,
			"subject":      subject,
		},
	}
}

func expandTopicInboundIPRules(input []interface{}) *[]topics.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]topics.InboundIPRule, 0)
	for _, item := range input {
		rawRule := item.(map[string]interface{})
		rules = append(rules, topics.InboundIPRule{
			Action: pointer.To(topics.IPActionType(rawRule["action"].(string))),
			IPMask: utils.String(rawRule["ip_mask"].(string)),
		})
	}
	return &rules
}

func flattenTopicInboundIPRules(input *[]topics.InboundIPRule) []interface{} {
	rules := make([]interface{}, 0)
	if input == nil {
		return rules
	}

	for _, r := range *input {
		action := ""
		if r.Action != nil {
			action = string(*r.Action)
		}

		rules = append(rules, map[string]interface{}{
			"action":  action,
			"ip_mask": pointer.From(r.IPMask),
		})
	}
	return rules
}
