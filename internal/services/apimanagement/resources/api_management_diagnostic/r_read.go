// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementDiagnosticRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := diagnostic.ParseDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *diagnosticId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *diagnosticId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *diagnosticId, err)
	}

	d.Set("resource_group_name", diagnosticId.ResourceGroupName)
	d.Set("api_management_name", diagnosticId.ServiceName)
	if model := resp.Model; model != nil {
		d.Set("identifier", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("api_management_logger_id", props.LoggerId)
			if props.Sampling != nil && props.Sampling.Percentage != nil {
				d.Set("sampling_percentage", pointer.From(props.Sampling.Percentage))
			}
			d.Set("always_log_errors", pointer.From(props.AlwaysLog) == diagnostic.AlwaysLogAllErrors)
			d.Set("verbosity", pointer.From(props.Verbosity))
			d.Set("log_client_ip", pointer.From(props.LogClientIP))
			d.Set("http_correlation_protocol", pointer.From(props.HTTPCorrelationProtocol))
			if frontend := props.Frontend; frontend != nil {
				d.Set("frontend_request", flattenApiManagementDiagnosticHTTPMessageDiagnostic(frontend.Request))
				d.Set("frontend_response", flattenApiManagementDiagnosticHTTPMessageDiagnostic(frontend.Response))
			} else {
				d.Set("frontend_request", nil)
				d.Set("frontend_response", nil)
			}
			if backend := props.Backend; backend != nil {
				d.Set("backend_request", flattenApiManagementDiagnosticHTTPMessageDiagnostic(backend.Request))
				d.Set("backend_response", flattenApiManagementDiagnosticHTTPMessageDiagnostic(backend.Response))
			} else {
				d.Set("backend_request", nil)
				d.Set("backend_response", nil)
			}

			format := ""
			if props.OperationNameFormat != nil {
				format = pointer.FromEnum(props.OperationNameFormat)
			}
			d.Set("operation_name_format", format)
		}
	}

	return nil
}
