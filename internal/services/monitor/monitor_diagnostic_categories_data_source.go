// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettingscategories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMonitorDiagnosticCategories() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceMonitorDiagnosticCategoriesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"log_category_types": {
				Type:     pluginsdk.TypeSet,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
				Computed: true,
			},

			"log_category_groups": {
				Type:     pluginsdk.TypeSet,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
				Computed: true,
			},

			"metrics": {
				Type:     pluginsdk.TypeSet,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
				Computed: true,
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["logs"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeSet,
			Elem:       &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:        pluginsdk.HashString,
			Computed:   true,
			Deprecated: "`logs` will be removed in favour of the property `log_category_types` in version 4.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func dataSourceMonitorDiagnosticCategoriesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	categoriesClient := meta.(*clients.Client).Monitor.DiagnosticSettingsCategoryClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	actualResourceId := commonids.NewScopeID(d.Get("resource_id").(string))
	// trim off the leading `/` since the CheckExistenceByID / List methods don't expect it
	resourceId := strings.TrimPrefix(actualResourceId.Scope, "/")
	resourceIdToList, err := commonids.ParseScopeID(resourceId)
	if err != nil {
		return fmt.Errorf("parsing resource id error: %+v", err)
	}

	// then retrieve the possible Diagnostics Categories for this Resource
	categories, err := categoriesClient.DiagnosticSettingsCategoryList(ctx, *resourceIdToList)
	if err != nil {
		return fmt.Errorf("retrieving Diagnostics Categories for Resource %q: %+v", actualResourceId, err)
	}

	if categories.Model == nil && categories.Model.Value == nil {
		return fmt.Errorf("retrieving Diagnostics Categories for Resource %q: `categories.Value` was nil", actualResourceId)
	}

	d.SetId(actualResourceId.ID())
	val := *categories.Model.Value

	metrics := make([]string, 0)
	logs := make([]string, 0)
	categoryGroups := make([]string, 0)

	for _, v := range val {
		if v.Name == nil {
			continue
		}

		if category := v.Properties; category != nil {
			if category.CategoryGroups != nil {
				categoryGroups = append(categoryGroups, *category.CategoryGroups...)
			}
			if category.CategoryType != nil {
				switch *category.CategoryType {
				case diagnosticsettingscategories.CategoryTypeLogs:
					logs = append(logs, *v.Name)
				case diagnosticsettingscategories.CategoryTypeMetrics:
					metrics = append(metrics, *v.Name)
				default:
					return fmt.Errorf("Unsupported category type %q", string(*category.CategoryType))
				}
			}
		}
	}

	if err := d.Set("log_category_types", logs); err != nil {
		return fmt.Errorf("setting `log_category_types`: %+v", err)
	}

	if !features.FourPointOhBeta() {
		if err := d.Set("logs", logs); err != nil {
			return fmt.Errorf("setting `log`: %+v", err)
		}
	}

	if err := d.Set("metrics", metrics); err != nil {
		return fmt.Errorf("setting `metrics`: %+v", err)
	}

	if err := d.Set("log_category_groups", categoryGroups); err != nil {
		return fmt.Errorf("setting `log_category_groups`: %+v", err)
	}
	return nil
}
