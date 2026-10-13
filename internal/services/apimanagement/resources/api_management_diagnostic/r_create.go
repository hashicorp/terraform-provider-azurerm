// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_diagnostic

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/diagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementDiagnosticCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnostic.NewDiagnosticID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("identifier").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_diagnostic", id.ID())
			}
		}
	}

	parameters := diagnostic.DiagnosticContract{
		Properties: &diagnostic.DiagnosticContractProperties{
			LoggerId: d.Get("api_management_logger_id").(string),
		},
	}

	if operationNameFormat, ok := d.GetOk("operation_name_format"); ok {
		if d.Get("identifier") == "applicationinsights" {
			parameters.Properties.OperationNameFormat = pointer.ToEnum[diagnostic.OperationNameFormat](operationNameFormat.(string))
		}
	}

	if samplingPercentage, ok := d.GetOk("sampling_percentage"); ok {
		parameters.Properties.Sampling = &diagnostic.SamplingSettings{
			SamplingType: pointer.To(diagnostic.SamplingTypeFixed),
			Percentage:   pointer.To(samplingPercentage.(float64)),
		}
	} else {
		parameters.Properties.Sampling = nil
	}

	if alwaysLogErrors, ok := d.GetOk("always_log_errors"); ok && alwaysLogErrors.(bool) {
		parameters.Properties.AlwaysLog = pointer.To(diagnostic.AlwaysLogAllErrors)
	}

	if verbosity, ok := d.GetOk("verbosity"); ok {
		parameters.Properties.Verbosity = pointer.ToEnum[diagnostic.Verbosity](verbosity.(string))
	}

	//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
	if logClientIP, exists := d.GetOkExists("log_client_ip"); exists { //nolint:staticcheck
		parameters.Properties.LogClientIP = pointer.To(logClientIP.(bool))
	}

	if httpCorrelationProtocol, ok := d.GetOk("http_correlation_protocol"); ok {
		parameters.Properties.HTTPCorrelationProtocol = pointer.ToEnum[diagnostic.HTTPCorrelationProtocol](httpCorrelationProtocol.(string))
	}

	frontendRequest, frontendRequestSet := d.GetOk("frontend_request")
	frontendResponse, frontendResponseSet := d.GetOk("frontend_response")
	if frontendRequestSet || frontendResponseSet {
		parameters.Properties.Frontend = &diagnostic.PipelineDiagnosticSettings{}
		if frontendRequestSet {
			parameters.Properties.Frontend.Request = expandApiManagementDiagnosticHTTPMessageDiagnostic(frontendRequest.([]any))
		}
		if frontendResponseSet {
			parameters.Properties.Frontend.Response = expandApiManagementDiagnosticHTTPMessageDiagnostic(frontendResponse.([]any))
		}
	}

	backendRequest, backendRequestSet := d.GetOk("backend_request")
	backendResponse, backendResponseSet := d.GetOk("backend_response")
	if backendRequestSet || backendResponseSet {
		parameters.Properties.Backend = &diagnostic.PipelineDiagnosticSettings{}
		if backendRequestSet {
			parameters.Properties.Backend.Request = expandApiManagementDiagnosticHTTPMessageDiagnostic(backendRequest.([]any))
		}
		if backendResponseSet {
			parameters.Properties.Backend.Response = expandApiManagementDiagnosticHTTPMessageDiagnostic(backendResponse.([]any))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, diagnostic.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementDiagnosticRead(d, meta)
}
