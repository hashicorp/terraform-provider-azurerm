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
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceEventGridDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventGridDomainCreate,
		Read:   resourceEventGridDomainRead,
		Update: resourceEventGridDomainUpdate,
		Delete: resourceEventGridDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := domains.ParseDomainID(id)
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
						"EventGrid domain name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"input_schema": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(domains.InputSchemaEventGridSchema),
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(domains.PossibleValuesForInputSchema(), false),
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

			"auto_create_topic_with_first_subscription": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"auto_delete_topic_with_last_subscription": {
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
							Default:  string(domains.IPActionTypeAllow),
							ValidateFunc: validation.StringInSlice([]string{
								string(domains.IPActionTypeAllow),
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

func resourceEventGridDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Domains
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := domains.NewDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_eventgrid_domain", id.ID())
	}

	inboundIPRules := expandDomainInboundIPRules(d.Get("inbound_ip_rule").([]interface{}))
	publicNetworkAccess := domains.PublicNetworkAccessDisabled
	if v, ok := d.GetOk("public_network_access_enabled"); ok && v.(bool) {
		publicNetworkAccess = domains.PublicNetworkAccessEnabled
	}

	domain := domains.Domain{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &domains.DomainProperties{
			AutoCreateTopicWithFirstSubscription: utils.Bool(d.Get("auto_create_topic_with_first_subscription").(bool)),
			AutoDeleteTopicWithLastSubscription:  utils.Bool(d.Get("auto_delete_topic_with_last_subscription").(bool)),
			DisableLocalAuth:                     utils.Bool(!d.Get("local_auth_enabled").(bool)),
			InboundIPRules:                       inboundIPRules,
			InputSchema:                          pointer.To(domains.InputSchema(d.Get("input_schema").(string))),
			InputSchemaMapping:                   expandDomainInputMapping(d),
			PublicNetworkAccess:                  pointer.To(publicNetworkAccess),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("identity"); ok {
		identity, err := identity.ExpandSystemAndUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		domain.Identity = identity
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, domain); err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceEventGridDomainRead(d, meta)
}

func resourceEventGridDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Domains
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domains.ParseDomainID(d.Id())
	if err != nil {
		return err
	}

	payload := domains.DomainUpdateParameters{Properties: &domains.DomainUpdateParameterProperties{}}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := domains.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = domains.PublicNetworkAccessEnabled
		}

		payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("local_auth_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
	}

	if d.HasChange("auto_create_topic_with_first_subscription") {
		payload.Properties.AutoCreateTopicWithFirstSubscription = pointer.To(d.Get("auto_create_topic_with_first_subscription").(bool))
	}

	if d.HasChange("auto_delete_topic_with_last_subscription") {
		payload.Properties.AutoDeleteTopicWithLastSubscription = pointer.To(d.Get("auto_delete_topic_with_last_subscription").(bool))
	}

	if d.HasChange("inbound_ip_rule") {
		inboundIpRule := d.Get("inbound_ip_rule").([]interface{})

		if len(inboundIpRule) == 0 {
			payload.Properties.InboundIPRules = pointer.To([]domains.InboundIPRule{})
		} else {
			payload.Properties.InboundIPRules = expandDomainInboundIPRules(inboundIpRule)
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %s", id, err)
	}

	return resourceEventGridDomainRead(d, meta)
}

func resourceEventGridDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

			inputMappingFields := flattenDomainInputMapping(props.InputSchemaMapping)
			if err := d.Set("input_mapping_fields", inputMappingFields); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			inputMappingDefaultValues := flattenDomainInputMappingDefaultValues(props.InputSchemaMapping)
			if err := d.Set("input_mapping_default_values", inputMappingDefaultValues); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields`: %+v", err)
			}

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == domains.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			inboundIPRules := flattenDomainInboundIPRules(props.InboundIPRules)
			if err := d.Set("inbound_ip_rule", inboundIPRules); err != nil {
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

func resourceEventGridDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Domains
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domains.ParseDomainID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandDomainInputMapping(d *pluginsdk.ResourceData) *domains.JsonInputSchemaMapping {
	imf, imfok := d.GetOk("input_mapping_fields")

	imdv, imdvok := d.GetOk("input_mapping_default_values")

	if !imfok && !imdvok {
		return nil
	}

	jismp := domains.JsonInputSchemaMappingProperties{}

	if imfok {
		mappings := imf.([]interface{})
		if len(mappings) > 0 && mappings[0] != nil {
			mapping := mappings[0].(map[string]interface{})

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
		mappings := imdv.([]interface{})
		if len(mappings) > 0 && mappings[0] != nil {
			mapping := mappings[0].(map[string]interface{})

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

func flattenDomainInputMapping(input domains.InputSchemaMapping) []interface{} {
	output := make([]interface{}, 0)
	val, ok := input.(domains.JsonInputSchemaMapping)
	if ok {
		if props := val.Properties; props != nil {
			eventTime := ""
			if props.EventTime != nil && props.EventTime.SourceField != nil {
				eventTime = *props.EventTime.SourceField
			}

			id := ""
			if props.Id != nil && props.Id.SourceField != nil {
				id = *props.Id.SourceField
			}

			topic := ""
			if props.Topic != nil && props.Topic.SourceField != nil {
				topic = *props.Topic.SourceField
			}

			dataVersion := ""
			if props.DataVersion != nil && props.DataVersion.SourceField != nil {
				dataVersion = *props.DataVersion.SourceField
			}

			eventType := ""
			if props.EventType != nil && props.EventType.SourceField != nil {
				eventType = *props.EventType.SourceField
			}

			subject := ""
			if props.Subject != nil && props.Subject.SourceField != nil {
				subject = *props.Subject.SourceField
			}

			output = append(output, map[string]interface{}{
				"data_version": dataVersion,
				"event_type":   eventType,
				"event_time":   eventTime,
				"id":           id,
				"topic":        topic,
				"subject":      subject,
			})
		}
	}
	return output
}

func flattenDomainInputMappingDefaultValues(input domains.InputSchemaMapping) []interface{} {
	output := make([]interface{}, 0)
	val, ok := input.(domains.JsonInputSchemaMapping)
	if ok {
		if props := val.Properties; props != nil {
			dataVersion := ""
			if props.DataVersion != nil && props.DataVersion.DefaultValue != nil {
				dataVersion = *props.DataVersion.DefaultValue
			}

			eventType := ""
			if props.EventType != nil && props.EventType.DefaultValue != nil {
				eventType = *props.EventType.DefaultValue
			}

			subject := ""
			if props.Subject != nil && props.Subject.DefaultValue != nil {
				subject = *props.Subject.DefaultValue
			}

			output = append(output, map[string]interface{}{
				"data_version": dataVersion,
				"event_type":   eventType,
				"subject":      subject,
			})
		}
	}

	return output
}

func expandDomainInboundIPRules(input []interface{}) *[]domains.InboundIPRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]domains.InboundIPRule, 0)
	for _, item := range input {
		rawRule := item.(map[string]interface{})
		rules = append(rules, domains.InboundIPRule{
			Action: pointer.To(domains.IPActionType(rawRule["action"].(string))),
			IPMask: pointer.To(rawRule["ip_mask"].(string)),
		})
	}
	return &rules
}

func flattenDomainInboundIPRules(input *[]domains.InboundIPRule) []interface{} {
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
