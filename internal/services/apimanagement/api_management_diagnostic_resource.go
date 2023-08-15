// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/diagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementDiagnostic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementDiagnosticCreateUpdate,
		Read:   resourceApiManagementDiagnosticRead,
		Update: resourceApiManagementDiagnosticCreateUpdate,
		Delete: resourceApiManagementDiagnosticDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := diagnostic.ParseDiagnosticID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"identifier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"applicationinsights",
					"azuremonitor",
				}, false),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"api_management_logger_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: logger.ValidateLoggerID,
			},

			"sampling_percentage": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatBetween(0.0, 100.0),
			},

			"always_log_errors": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"verbosity": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(diagnostic.VerbosityVerbose),
					string(diagnostic.VerbosityInformation),
					string(diagnostic.VerbosityError),
				}, false),
			},

			"log_client_ip": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"http_correlation_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(diagnostic.HTTPCorrelationProtocolNone),
					string(diagnostic.HTTPCorrelationProtocolLegacy),
					string(diagnostic.HTTPCorrelationProtocolWThreeC),
				}, false),
			},

			"frontend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"frontend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"operation_name_format": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(diagnostic.OperationNameFormatName),
				ValidateFunc: validation.StringInSlice([]string{
					string(diagnostic.OperationNameFormatName),
					string(diagnostic.OperationNameFormatUrl),
				}, false),
			},
		},
	}
}

func resourceApiManagementDiagnosticCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnostic.NewDiagnosticID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("identifier").(string))

	if d.IsNewResource() {
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

	parameters := diagnostic.DiagnosticContract{
		Properties: &diagnostic.DiagnosticContractProperties{
			LoggerId:            d.Get("api_management_logger_id").(string),
			OperationNameFormat: pointer.To(diagnostic.OperationNameFormat(d.Get("operation_name_format").(string))),
		},
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
		parameters.Properties.Verbosity = pointer.To(diagnostic.Verbosity(verbosity.(string)))
	}

	//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
	if logClientIP, exists := d.GetOkExists("log_client_ip"); exists { //nolint:staticcheck
		parameters.Properties.LogClientIP = pointer.To(logClientIP.(bool))
	}

	if httpCorrelationProtocol, ok := d.GetOk("http_correlation_protocol"); ok {
		parameters.Properties.HTTPCorrelationProtocol = pointer.To(diagnostic.HTTPCorrelationProtocol(httpCorrelationProtocol.(string)))
	}

	frontendRequest, frontendRequestSet := d.GetOk("frontend_request")
	frontendResponse, frontendResponseSet := d.GetOk("frontend_response")
	if frontendRequestSet || frontendResponseSet {
		parameters.Properties.Frontend = &diagnostic.PipelineDiagnosticSettings{}
		if frontendRequestSet {
			parameters.Properties.Frontend.Request = expandApiManagementDiagnosticHTTPMessageDiagnostic(frontendRequest.([]interface{}))
		}
		if frontendResponseSet {
			parameters.Properties.Frontend.Response = expandApiManagementDiagnosticHTTPMessageDiagnostic(frontendResponse.([]interface{}))
		}
	}

	backendRequest, backendRequestSet := d.GetOk("backend_request")
	backendResponse, backendResponseSet := d.GetOk("backend_response")
	if backendRequestSet || backendResponseSet {
		parameters.Properties.Backend = &diagnostic.PipelineDiagnosticSettings{}
		if backendRequestSet {
			parameters.Properties.Backend.Request = expandApiManagementDiagnosticHTTPMessageDiagnostic(backendRequest.([]interface{}))
		}
		if backendResponseSet {
			parameters.Properties.Backend.Response = expandApiManagementDiagnosticHTTPMessageDiagnostic(backendResponse.([]interface{}))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, diagnostic.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementDiagnosticRead(d, meta)
}

func resourceApiManagementDiagnosticRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			format := string(diagnostic.OperationNameFormatName)
			if props.OperationNameFormat != nil {
				format = string(pointer.From(props.OperationNameFormat))
			}
			d.Set("operation_name_format", format)
		}
	}

	return nil
}

func resourceApiManagementDiagnosticDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := diagnostic.ParseDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *diagnosticId, diagnostic.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *diagnosticId, err)
		}
	}

	return nil
}

func expandApiManagementDiagnosticHTTPMessageDiagnostic(input []interface{}) *diagnostic.HTTPMessageDiagnostic {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &diagnostic.HTTPMessageDiagnostic{
		Body: &diagnostic.BodyDiagnosticSettings{},
	}

	if bodyBytes, ok := v["body_bytes"]; ok {
		result.Body.Bytes = pointer.To(int64(bodyBytes.(int)))
	}
	if headersSetRaw, ok := v["headers_to_log"]; ok {
		headersSet := headersSetRaw.(*pluginsdk.Set).List()
		headers := []string{}
		for _, header := range headersSet {
			headers = append(headers, header.(string))
		}
		result.Headers = pointer.To(headers)
	}

	result.DataMasking = expandApiManagementDataMasking(v["data_masking"].([]interface{}))

	return result
}

func flattenApiManagementDiagnosticHTTPMessageDiagnostic(input *diagnostic.HTTPMessageDiagnostic) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	diagnostic := map[string]interface{}{}

	if input.Body != nil && input.Body.Bytes != nil {
		diagnostic["body_bytes"] = pointer.From(input.Body.Bytes)
	}

	if input.Headers != nil {
		diagnostic["headers_to_log"] = set.FromStringSlice(*input.Headers)
	}

	diagnostic["data_masking"] = flattenApiManagementDataMasking(input.DataMasking)

	result = append(result, diagnostic)

	return result
}

func expandApiManagementDataMasking(input []interface{}) *diagnostic.DataMasking {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	inputRaw := input[0].(map[string]interface{})
	return &diagnostic.DataMasking{
		QueryParams: expandApiManagementDataMaskingEntityList(inputRaw["query_params"].([]interface{})),
		Headers:     expandApiManagementDataMaskingEntityList(inputRaw["headers"].([]interface{})),
	}
}

func expandApiManagementDataMaskingEntityList(input []interface{}) *[]diagnostic.DataMaskingEntity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]diagnostic.DataMaskingEntity, 0)
	for _, v := range input {
		entity := v.(map[string]interface{})
		result = append(result, diagnostic.DataMaskingEntity{
			Mode:  pointer.To(diagnostic.DataMaskingMode(entity["mode"].(string))),
			Value: pointer.To(entity["value"].(string)),
		})
	}
	return &result
}

func flattenApiManagementDataMasking(dataMasking *diagnostic.DataMasking) []interface{} {
	if dataMasking == nil {
		return []interface{}{}
	}

	var queryParams, headers []interface{}
	if dataMasking.QueryParams != nil {
		queryParams = flattenApiManagementDataMaskingEntityList(dataMasking.QueryParams)
	}
	if dataMasking.Headers != nil {
		headers = flattenApiManagementDataMaskingEntityList(dataMasking.Headers)
	}

	return []interface{}{
		map[string]interface{}{
			"query_params": queryParams,
			"headers":      headers,
		},
	}
}

func flattenApiManagementDataMaskingEntityList(dataMaskingList *[]diagnostic.DataMaskingEntity) []interface{} {
	if dataMaskingList == nil || len(*dataMaskingList) == 0 {
		return []interface{}{}
	}

	result := []interface{}{}

	for _, entity := range *dataMaskingList {
		result = append(result, map[string]interface{}{
			"mode":  pointer.From(entity.Mode),
			"value": pointer.From(entity.Value),
		})
	}

	return result
}
