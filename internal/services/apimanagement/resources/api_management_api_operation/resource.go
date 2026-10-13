// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationCreateUpdate,
		Read:   resourceApiManagementApiOperationRead,
		Update: resourceApiManagementApiOperationCreateUpdate,
		Delete: resourceApiManagementApiOperationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apioperation.ParseOperationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: apiManagementApiOperationSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
			// Get the parameters used in url_template
			urlTemplate := diff.Get("url_template").(string)
			re := regexp.MustCompile(`\{([^}]+)\}`)
			matches := re.FindAllStringSubmatch(urlTemplate, -1)
			urlTemplateParamSet := make(map[string]struct{})
			for _, match := range matches {
				if len(match) > 1 {
					// Since Azure API Management's `url_template` supports two formats: {name} and {*name}, `*` should be removed when getting name.
					urlTemplateParamSet[strings.TrimPrefix(match[1], "*")] = struct{}{}
				}
			}

			// Get the parameters defined in template_parameter
			templateParametersRaw := diff.Get("template_parameter").([]any)
			templateParameterSet := make(map[string]struct{})
			for _, p := range templateParametersRaw {
				paramValue := p.(map[string]any)
				templateParameterSet[paramValue["name"].(string)] = struct{}{}
			}

			for key := range urlTemplateParamSet {
				if _, found := templateParameterSet[key]; !found {
					return fmt.Errorf("template parameter `%s` used in `url_template` is not defined in `template_parameter`", key)
				}
			}

			for key := range templateParameterSet {
				if _, found := urlTemplateParamSet[key]; !found {
					return fmt.Errorf("template parameter `%s` defined in `template_parameter` is not used in `url_template`", key)
				}
			}

			return nil
		}),
	}
}
