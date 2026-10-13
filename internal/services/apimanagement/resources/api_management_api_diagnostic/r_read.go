// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_diagnostic

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiDiagnosticRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := apidiagnostic.ParseApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(diagnosticId.ApiId)

	newId := apidiagnostic.NewApiDiagnosticID(diagnosticId.SubscriptionId, diagnosticId.ResourceGroupName, diagnosticId.ServiceName, apiName, diagnosticId.DiagnosticId)
	resp, err := client.Get(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", newId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}
	d.Set("api_name", apiName)
	d.Set("resource_group_name", diagnosticId.ResourceGroupName)
	d.Set("api_management_name", diagnosticId.ServiceName)
	if model := resp.Model; model != nil {
		d.Set("identifier", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("api_management_logger_id", props.LoggerId)
			if props.Sampling != nil && props.Sampling.Percentage != nil {
				d.Set("sampling_percentage", pointer.From(props.Sampling.Percentage))
			}
			d.Set("always_log_errors", pointer.From(props.AlwaysLog) == apidiagnostic.AlwaysLogAllErrors)
			d.Set("verbosity", pointer.From(props.Verbosity))
			d.Set("log_client_ip", pointer.From(props.LogClientIP))
			d.Set("http_correlation_protocol", pointer.From(props.HTTPCorrelationProtocol))
			if frontend := props.Frontend; frontend != nil {
				d.Set("frontend_request", flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(frontend.Request))
				d.Set("frontend_response", flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(frontend.Response))
			} else {
				d.Set("frontend_request", nil)
				d.Set("frontend_response", nil)
			}
			if backend := props.Backend; backend != nil {
				d.Set("backend_request", flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(backend.Request))
				d.Set("backend_response", flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(backend.Response))
			} else {
				d.Set("backend_request", nil)
				d.Set("backend_response", nil)
			}

			format := string(apidiagnostic.OperationNameFormatName)
			if props.OperationNameFormat != nil {
				format = pointer.FromEnum(props.OperationNameFormat)
			}
			d.Set("operation_name_format", format)
		}
	}

	return nil
}
