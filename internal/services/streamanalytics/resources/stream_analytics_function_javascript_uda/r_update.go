// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_function_javascript_uda

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsFunctionUDAUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).FunctionsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := functions.ParseFunctionID(d.Id())
	if err != nil {
		return err
	}

	props := functions.Function{
		Properties: &functions.AggregateFunctionProperties{
			Properties: &functions.FunctionConfiguration{
				Binding: &functions.JavaScriptFunctionBinding{
					Properties: &functions.JavaScriptFunctionBindingProperties{
						Script: pointer.To(d.Get("script").(string)),
					},
				},
				Inputs: expandStreamAnalyticsFunctionUDAInputs(d.Get("input").([]any)),
				Output: expandStreamAnalyticsFunctionUDAOutput(d.Get("output").([]any)),
			},
		},
	}

	var opts functions.UpdateOperationOptions
	if _, err := client.Update(ctx, *id, props, opts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsFunctionUDARead(d, meta)
}
