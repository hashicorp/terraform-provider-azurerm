// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorActivityLogAlert() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorActivityLogAlertCreateUpdate,
		Read:   resourceMonitorActivityLogAlertRead,
		Update: resourceMonitorActivityLogAlertCreateUpdate,
		Delete: resourceMonitorActivityLogAlertDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := activitylogalertsapis.ParseActivityLogAlertID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ActivityLogAlertUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"scopes": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: pluginsdk.HashString,
			},

			"criteria": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Administrative",
								"Autoscale",
								"Policy",
								"Recommendation",
								"ResourceHealth",
								"Security",
								"ServiceHealth",
							}, false),
						},
						"operation_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"caller": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"level": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Verbose",
								"Informational",
								"Warning",
								"Error",
								"Critical",
							}, false),
							ConflictsWith: []string{"criteria.0.levels"},
						},
						"levels": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Verbose",
									"Informational",
									"Warning",
									"Error",
									"Critical",
								}, false),
							},
							ConflictsWith: []string{"criteria.0.level"},
						},
						"resource_provider": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.resource_providers"},
						},
						"resource_providers": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.resource_provider"},
						},
						"resource_type": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.resource_types"},
						},
						"resource_types": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.resource_type"},
						},
						"resource_group": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.resource_groups"},
						},
						"resource_groups": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.resource_group"},
						},
						"resource_id": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  azure.ValidateResourceID,
							ConflictsWith: []string{"criteria.0.resource_ids"},
						},
						"resource_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.resource_id"},
						},
						"status": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.statuses"},
						},
						"statuses": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.status"},
						},
						"sub_status": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.sub_statuses"},
						},
						"sub_statuses": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							ConflictsWith: []string{"criteria.0.sub_status"},
						},
						"recommendation_category": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Cost",
								"Reliability",
								"OperationalExcellence",
								"Performance",
								"HighAvailability",
							},
								false,
							),
							ConflictsWith: []string{"criteria.0.recommendation_type"},
						},
						"recommendation_impact": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"High",
								"Medium",
								"Low",
							},
								false,
							),
							ConflictsWith: []string{"criteria.0.recommendation_type"},
						},
						"recommendation_type": {
							Type:          pluginsdk.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringIsNotEmpty,
							ConflictsWith: []string{"criteria.0.recommendation_category", "criteria.0.recommendation_impact"},
						},
						// lintignore:XS003
						"resource_health": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"current": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Available",
												"Degraded",
												"Unavailable",
												"Unknown",
											},
												false,
											),
										},
										Set: pluginsdk.HashString,
									},
									"previous": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Available",
												"Degraded",
												"Unavailable",
												"Unknown",
											},
												false,
											),
										},
										Set: pluginsdk.HashString,
									},
									"reason": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"PlatformInitiated",
												"UserInitiated",
												"Unknown",
											},
												false,
											),
										},
										Set: pluginsdk.HashString,
									},
								},
							},
							ConflictsWith: []string{"criteria.0.caller", "criteria.0.service_health"},
						},
						// lintignore:XS003
						"service_health": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"events": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Incident",
												"Maintenance",
												"Informational",
												"ActionRequired",
												"Security",
											},
												false,
											),
										},
										Set: pluginsdk.HashString,
									},
									"locations": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										Set: pluginsdk.HashString,
									},
									"services": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										Set: pluginsdk.HashString,
									},
								},
							},
							ConflictsWith: []string{"criteria.0.caller", "criteria.0.resource_health"},
						},
					},
				},
			},

			"action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_group_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"webhook_properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorActivityLogAlertCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActivityLogAlertsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := activitylogalertsapis.NewActivityLogAlertID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ActivityLogAlertsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_activity_log_alert", id.ID())
		}
	}

	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	scopesRaw := d.Get("scopes").(*pluginsdk.Set).List()
	criteriaRaw := d.Get("criteria").([]interface{})
	actionRaw := d.Get("action").([]interface{})

	t := d.Get("tags").(map[string]interface{})
	parameters := activitylogalertsapis.ActivityLogAlertResource{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &activitylogalertsapis.AlertRuleProperties{
			Enabled:     utils.Bool(enabled),
			Description: utils.String(description),
			Scopes:      expandStringValues(scopesRaw),
			Condition:   expandMonitorActivityLogAlertCriteria(criteriaRaw),
			Actions:     expandMonitorActivityLogAlertAction(actionRaw),
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err := client.ActivityLogAlertsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorActivityLogAlertRead(d, meta)
}

func resourceMonitorActivityLogAlertRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActivityLogAlertsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := activitylogalertsapis.ParseActivityLogAlertID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ActivityLogAlertsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.ActivityLogAlertName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("enabled", props.Enabled)
			d.Set("description", props.Description)

			var scopes []interface{}
			if props.Scopes != nil {
				scopes = utils.FlattenStringSlice(&props.Scopes)
			}
			if err := d.Set("scopes", scopes); err != nil {
				return fmt.Errorf("setting `scopes`: %+v", err)
			}

			if err := d.Set("criteria", flattenMonitorActivityLogAlertCriteria(props.Condition)); err != nil {
				return fmt.Errorf("setting `criteria`: %+v", err)
			}
			if err := d.Set("action", flattenMonitorActivityLogAlertAction(props.Actions)); err != nil {
				return fmt.Errorf("setting `action`: %+v", err)
			}
		}
		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceMonitorActivityLogAlertDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActivityLogAlertsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := activitylogalertsapis.ParseActivityLogAlertID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.ActivityLogAlertsDelete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandMonitorActivityLogAlertCriteria(input []interface{}) activitylogalertsapis.AlertRuleAllOfCondition {
	conditions := make([]activitylogalertsapis.AlertRuleAnyOfOrLeafCondition, 0)
	v := input[0].(map[string]interface{})

	if category := v["category"].(string); category != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("category"),
			Equals: utils.String(category),
		})
	}

	if op := v["operation_name"].(string); op != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("operationName"),
			Equals: utils.String(op),
		})
	}

	if caller := v["caller"].(string); caller != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("caller"),
			Equals: utils.String(caller),
		})
	}

	if level := v["level"].(string); level != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("level"),
			Equals: utils.String(level),
		})
	}

	if levels := v["levels"].([]interface{}); len(levels) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(levels, "level"),
		})
	}

	if resourceProvider := v["resource_provider"].(string); resourceProvider != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("resourceProvider"),
			Equals: utils.String(resourceProvider),
		})
	}

	if resourceProviders := v["resource_providers"].([]interface{}); len(resourceProviders) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceProviders, "resourceProvider"),
		})
	}

	if resourceType := v["resource_type"].(string); resourceType != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("resourceType"),
			Equals: utils.String(resourceType),
		})
	}

	if resourceTypes := v["resource_types"].([]interface{}); len(resourceTypes) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceTypes, "resourceType"),
		})
	}

	if resourceGroup := v["resource_group"].(string); resourceGroup != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("resourceGroup"),
			Equals: utils.String(resourceGroup),
		})
	}

	if resourceGroups := v["resource_groups"].([]interface{}); len(resourceGroups) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceGroups, "resourceGroup"),
		})
	}

	if id := v["resource_id"].(string); id != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("resourceId"),
			Equals: utils.String(id),
		})
	}

	if resourceIds := v["resource_ids"].([]interface{}); len(resourceIds) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceIds, "resourceId"),
		})
	}

	if status := v["status"].(string); status != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("status"),
			Equals: utils.String(status),
		})
	}

	if statuses := v["statuses"].([]interface{}); len(statuses) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(statuses, "status"),
		})
	}

	if subStatus := v["sub_status"].(string); subStatus != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("subStatus"),
			Equals: utils.String(subStatus),
		})
	}

	if statuses := v["sub_statuses"].([]interface{}); len(statuses) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(statuses, "subStatus"),
		})
	}

	if recommendationType := v["recommendation_type"].(string); recommendationType != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("properties.recommendationType"),
			Equals: utils.String(recommendationType),
		})
	}

	if recommendationCategory := v["recommendation_category"].(string); recommendationCategory != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("properties.recommendationCategory"),
			Equals: utils.String(recommendationCategory),
		})
	}

	if recommendationImpact := v["recommendation_impact"].(string); recommendationImpact != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  utils.String("properties.recommendationImpact"),
			Equals: utils.String(recommendationImpact),
		})
	}

	if resourceHealth := v["resource_health"].([]interface{}); len(resourceHealth) > 0 {
		conditions = expandResourceHealth(resourceHealth, conditions)
	}

	if serviceHealth := v["service_health"].([]interface{}); len(serviceHealth) > 0 {
		conditions = expandServiceHealth(serviceHealth, conditions)
	}

	return activitylogalertsapis.AlertRuleAllOfCondition{
		AllOf: conditions,
	}
}

func expandAnyOfCondition(input []interface{}, field string) *[]activitylogalertsapis.AlertRuleLeafCondition {
	conditions := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
	for _, v := range input {
		conditions = append(conditions, activitylogalertsapis.AlertRuleLeafCondition{
			Field:  utils.String(field),
			Equals: utils.String(v.(string)),
		})
	}
	return &conditions
}

func expandResourceHealth(resourceHealth []interface{}, conditions []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition) []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition {
	for _, serviceItem := range resourceHealth {
		if serviceItem == nil {
			continue
		}
		vs := serviceItem.(map[string]interface{})

		cv := vs["current"].(*pluginsdk.Set)
		if len(cv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range cv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  utils.String("properties.currentHealthStatus"),
					Equals: utils.String(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		pv := vs["previous"].(*pluginsdk.Set)
		if len(pv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range pv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  utils.String("properties.previousHealthStatus"),
					Equals: utils.String(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		rv := vs["reason"].(*pluginsdk.Set)
		if len(rv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range rv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  utils.String("properties.cause"),
					Equals: utils.String(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}
	}
	return conditions
}

func expandServiceHealth(serviceHealth []interface{}, conditions []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition) []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition {
	for _, serviceItem := range serviceHealth {
		if serviceItem == nil {
			continue
		}
		vs := serviceItem.(map[string]interface{})
		rv := vs["locations"].(*pluginsdk.Set)
		if len(rv.List()) > 0 {
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				Field:       utils.String("properties.impactedServices[*].ImpactedRegions[*].RegionName"),
				ContainsAny: utils.ExpandStringSlice(rv.List()),
			})
		}

		ev := vs["events"].(*pluginsdk.Set)
		if len(ev.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range ev.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  utils.String("properties.incidentType"),
					Equals: utils.String(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		sv := vs["services"].(*pluginsdk.Set)
		if len(sv.List()) > 0 {
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				Field:       utils.String("properties.impactedServices[*].ServiceName"),
				ContainsAny: utils.ExpandStringSlice(sv.List()),
			})
		}
	}
	return conditions
}

func expandMonitorActivityLogAlertAction(input []interface{}) activitylogalertsapis.ActionList {
	actions := make([]activitylogalertsapis.ActionGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = pv.(string)
				}
			}

			actions = append(actions, activitylogalertsapis.ActionGroup{
				ActionGroupId:     agID,
				WebhookProperties: &props,
			})
		}
	}
	return activitylogalertsapis.ActionList{
		ActionGroups: &actions,
	}
}

func flattenMonitorActivityLogAlertCriteria(input activitylogalertsapis.AlertRuleAllOfCondition) []interface{} {
	result := make(map[string]interface{})
	if input.AllOf == nil {
		return []interface{}{result}
	}
	for _, condition := range input.AllOf {
		if condition.Field != nil && condition.Equals != nil {
			switch strings.ToLower(*condition.Field) {
			case "operationname":
				result["operation_name"] = *condition.Equals
			case "resourceprovider":
				result["resource_provider"] = *condition.Equals
			case "resourcetype":
				result["resource_type"] = *condition.Equals
			case "resourcegroup":
				result["resource_group"] = *condition.Equals
			case "resourceid":
				result["resource_id"] = *condition.Equals
			case "substatus":
				result["sub_status"] = *condition.Equals
			case "properties.recommendationtype":
				result["recommendation_type"] = *condition.Equals
			case "properties.recommendationcategory":
				result["recommendation_category"] = *condition.Equals
			case "properties.recommendationimpact":
				result["recommendation_impact"] = *condition.Equals
			case "caller", "category", "level", "status":
				result[*condition.Field] = *condition.Equals
			}
		}

		if condition.Field != nil && condition.ContainsAny != nil && len(*condition.ContainsAny) > 0 {
			switch strings.ToLower(*condition.Field) {
			case "resourceprovider":
				result["resource_providers"] = *condition.ContainsAny
			case "resourcetype":
				result["resource_types"] = *condition.ContainsAny
			case "resourcegroup":
				result["resource_groups"] = *condition.ContainsAny
			case "resourceid":
				result["resource_ids"] = *condition.ContainsAny
			case "substatus":
				result["sub_statuses"] = *condition.ContainsAny
			case "level":
				result["levels"] = *condition.ContainsAny
			case "status":
				result["statuses"] = *condition.ContainsAny
			}
		}

		if condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			values := make([]string, 0)
			for _, leafCondition := range *condition.AnyOf {
				if leafCondition.Field != nil && leafCondition.Equals != nil {
					values = append(values, *leafCondition.Equals)
				}
				switch strings.ToLower(*leafCondition.Field) {
				case "resourceprovider":
					result["resource_providers"] = values
				case "resourcetype":
					result["resource_types"] = values
				case "resourcegroup":
					result["resource_groups"] = values
				case "resourceid":
					result["resource_ids"] = values
				case "substatus":
					result["sub_statuses"] = values
				case "level":
					result["levels"] = values
				case "status":
					result["statuses"] = values
				}
			}
		}
	}

	if result["category"] == "ResourceHealth" {
		flattenMonitorActivityLogAlertResourceHealth(input, result)
	}

	if result["category"] == "ServiceHealth" {
		flattenMonitorActivityLogAlertServiceHealth(input, result)
	}

	return []interface{}{result}
}

func flattenMonitorActivityLogAlertResourceHealth(input activitylogalertsapis.AlertRuleAllOfCondition, result map[string]interface{}) {
	rhResult := make(map[string]interface{})

	for _, condition := range input.AllOf {
		if condition.Field == nil && condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			values := []string{}
			for _, cond := range *condition.AnyOf {
				if cond.Field != nil && cond.Equals != nil {
					values = append(values, *cond.Equals)
				}
				switch strings.ToLower(*cond.Field) {
				case "properties.currenthealthstatus":
					rhResult["current"] = values
				case "properties.previoushealthstatus":
					rhResult["previous"] = values
				case "properties.cause":
					rhResult["reason"] = values
				}
			}
		}
	}

	result["resource_health"] = []interface{}{rhResult}
}

func flattenMonitorActivityLogAlertServiceHealth(input activitylogalertsapis.AlertRuleAllOfCondition, result map[string]interface{}) {
	shResult := make(map[string]interface{})
	for _, condition := range input.AllOf {
		if condition.Field != nil && condition.ContainsAny != nil && len(*condition.ContainsAny) > 0 {
			switch strings.ToLower(*condition.Field) {
			case "properties.impactedservices[*].impactedregions[*].regionname":
				shResult["locations"] = *condition.ContainsAny
			case "properties.impactedservices[*].servicename":
				shResult["services"] = *condition.ContainsAny
			}
		}
		if condition.Field == nil && condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			events := []string{}
			for _, evCond := range *condition.AnyOf {
				if evCond.Field != nil && evCond.Equals != nil {
					events = append(events, *evCond.Equals)
				}
			}
			shResult["events"] = events
		}
	}

	result["service_health"] = []interface{}{shResult}
}

func flattenMonitorActivityLogAlertAction(input activitylogalertsapis.ActionList) (result []interface{}) {
	result = make([]interface{}, 0)
	if input.ActionGroups == nil {
		return
	}
	for _, action := range *input.ActionGroups {
		v := make(map[string]interface{})

		v["action_group_id"] = action.ActionGroupId

		props := make(map[string]interface{})
		if action.WebhookProperties != nil {
			for pk, pv := range *action.WebhookProperties {
				props[pk] = pv
			}
		}
		v["webhook_properties"] = props

		result = append(result, v)
	}
	return result
}
