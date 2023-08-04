// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph" // nolint: staticcheck
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceResources() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceResourcesRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Resource to search for",
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
			"resource_group_name": {
				Description:  "The name of the Resource group where the Resources are located",
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
			"type": {
				Description:  "The Resource Type of the Resources you want to list",
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
			"location": {
				Description:  "Only return resources that deployed at a specific location",
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
			"required_tags": tags.Schema(),
			"subscription_ids": {
				Description:   `Azure subscription ids against which to execute the query`,
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"management_group_ids"},
				MinItems:      1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"management_group_ids": {
				Description:   `Azure management group names against which to execute the query`,
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"subscription_ids"},
				MinItems:      1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
			"resources": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"resource_group_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"location": commonschema.LocationComputed(),
						"tags":     tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceResourcesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ResourceGraphClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	resourceName := d.Get("name").(string)
	resourceType := d.Get("type").(string)
	requiredTags := d.Get("required_tags").(map[string]interface{})
	resourceLocation := d.Get("location").(string)

	if resourceGroupName == "" && resourceName == "" && resourceType == "" {
		return fmt.Errorf("at least one of `name`, `resource_group_name` or `type` must be specified")
	}

	filterClauses := []string{}

	if resourceName != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("name == '%s'", resourceName))
	}
	if resourceType != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("type =~ '%s' ", resourceType))
	}
	if resourceGroupName != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("resourceGroup =~ '%s'", resourceGroupName))
	}
	if resourceLocation != "" {
		filterClauses = append(filterClauses, fmt.Sprintf("location =~ '%s'", resourceLocation))
	}
	if len(requiredTags) > 0 {
		for requiredTagName, requiredTagVal := range requiredTags {
			filterClauses = append(filterClauses, fmt.Sprintf("tags['%s'] == '%s'", requiredTagName, requiredTagVal))
		}
	}

	filter := strings.Join(filterClauses, " and ")
	query := fmt.Sprintf(`resources
| project name, id, subscriptionId, resourceGroup, type, location, tags
| where %s`, filter)

	opts := resourcegraph.QueryRequestOptions{
		ResultFormat: resourcegraph.ResultFormatObjectArray,
	}

	queryRequest := resourcegraph.QueryRequest{
		Options: &opts,
		Query:   &query,
	}

	if data, ok := d.GetOk("subscription_ids"); ok {
		subscriptions := data.(*schema.Set).List()
		subs := make([]string, len(subscriptions))
		for i, v := range subscriptions {
			subs[i] = v.(string)
		}
		queryRequest.Subscriptions = &subs
	}

	if data, ok := d.GetOk("management_group_ids"); ok {
		managementgroups := data.(*schema.Set).List()
		grps := make([]string, len(managementgroups))
		for i, v := range managementgroups {
			grps[i] = v.(string)
		}
		queryRequest.ManagementGroups = &grps
	}

	data, err := doResourceQuery(ctx, client, queryRequest)
	if err != nil {
		return err
	}
	d.SetId("resource-" + uuid.New().String())
	if err := d.Set("resources", flattenResult(data)); err != nil {
		return fmt.Errorf("failed to set resources property: %w", err)
	}

	return nil
}

func doResourceQuery(ctx context.Context, client *resourcegraph.BaseClient, queryRequest resourcegraph.QueryRequest) (data []map[string]interface{}, err error) {
	data = make([]map[string]interface{}, 0)
	for {
		resp, err := client.Resources(ctx, queryRequest)
		if err != nil {
			err = fmt.Errorf("query failed: %w", err)
			return nil, err
		}
		if resp.Count != nil && *resp.Count > 0 && resp.Data != nil {
			for _, v := range resp.Data.([]interface{}) {
				item := v.(map[string]interface{})
				data = append(data, item)
			}
		}
		if resp.SkipToken == nil {
			break
		}
		queryRequest.Options.SkipToken = resp.SkipToken
	}
	return
}

func flattenResult(data []map[string]interface{}) (result []map[string]interface{}) {
	result = make([]map[string]interface{}, len(data))
	for i, v := range data {
		result[i] = map[string]interface{}{
			"name":                v["name"],
			"id":                  v["id"],
			"subscription_id":     v["subscriptionId"],
			"resource_group_name": v["resourceGroup"],
			"type":                v["type"],
			"location":            v["location"],
			"tags":                v["tags"],
		}
	}
	return
}
