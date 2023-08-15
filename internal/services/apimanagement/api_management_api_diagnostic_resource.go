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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apidiagnostic"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiDiagnostic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiDiagnosticCreateUpdate,
		Read:   resourceApiManagementApiDiagnosticRead,
		Update: resourceApiManagementApiDiagnosticCreateUpdate,
		Delete: resourceApiManagementApiDiagnosticDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apidiagnostic.ParseApiDiagnosticID(id)
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

			"api_name": schemaz.SchemaApiManagementApiName(),

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
					string(apidiagnostic.VerbosityVerbose),
					string(apidiagnostic.VerbosityInformation),
					string(apidiagnostic.VerbosityError),
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
					string(apidiagnostic.HTTPCorrelationProtocolNone),
					string(apidiagnostic.HTTPCorrelationProtocolLegacy),
					string(apidiagnostic.HTTPCorrelationProtocolWThreeC),
				}, false),
			},

			"frontend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"frontend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"operation_name_format": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(apidiagnostic.OperationNameFormatName),
				ValidateFunc: validation.StringInSlice([]string{
					string(apidiagnostic.OperationNameFormatName),
					string(apidiagnostic.OperationNameFormatUrl),
				}, false),
			},
		},
	}
}

func resourceApiManagementApiDiagnosticAdditionalContentSchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"body_bytes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 8192),
				},
				"headers_to_log": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Set: pluginsdk.HashString,
				},
				"data_masking": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"query_params": schemaApiManagementDataMaskingEntityList(),
							"headers":      schemaApiManagementDataMaskingEntityList(),
						},
					},
				},
			},
		},
	}
}

func resourceApiManagementApiDiagnosticCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apidiagnostic.NewApiDiagnosticID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("identifier").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Diagnostic %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_diagnostic", id.ID())
		}
	}

	parameters := apidiagnostic.DiagnosticContract{
		Properties: &apidiagnostic.DiagnosticContractProperties{
			LoggerId:            d.Get("api_management_logger_id").(string),
			OperationNameFormat: pointer.To(apidiagnostic.OperationNameFormat(d.Get("operation_name_format").(string))),
		},
	}

	samplingPercentage := d.GetRawConfig().AsValueMap()["sampling_percentage"]
	if !samplingPercentage.IsNull() {
		parameters.Properties.Sampling = &apidiagnostic.SamplingSettings{
			SamplingType: pointer.To(apidiagnostic.SamplingTypeFixed),
			Percentage:   pointer.To(d.Get("sampling_percentage").(float64)),
		}
	} else {
		parameters.Properties.Sampling = nil
	}

	if alwaysLogErrors, ok := d.GetOk("always_log_errors"); ok && alwaysLogErrors.(bool) {
		parameters.Properties.AlwaysLog = pointer.To(apidiagnostic.AlwaysLogAllErrors)
	}

	if verbosity, ok := d.GetOk("verbosity"); ok {
		parameters.Properties.Verbosity = pointer.To(apidiagnostic.Verbosity(verbosity.(string)))
	}

	//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
	if logClientIP, exists := d.GetOkExists("log_client_ip"); exists { //nolint:staticcheck
		parameters.Properties.LogClientIP = pointer.To(logClientIP.(bool))
	}

	if httpCorrelationProtocol, ok := d.GetOk("http_correlation_protocol"); ok {
		parameters.Properties.HTTPCorrelationProtocol = pointer.To(apidiagnostic.HTTPCorrelationProtocol(httpCorrelationProtocol.(string)))
	}

	frontendRequest, frontendRequestSet := d.GetOk("frontend_request")
	frontendResponse, frontendResponseSet := d.GetOk("frontend_response")
	if frontendRequestSet || frontendResponseSet {
		parameters.Properties.Frontend = &apidiagnostic.PipelineDiagnosticSettings{}
		if frontendRequestSet {
			parameters.Properties.Frontend.Request = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(frontendRequest.([]interface{}))
		}
		if frontendResponseSet {
			parameters.Properties.Frontend.Response = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(frontendResponse.([]interface{}))
		}
	}

	backendRequest, backendRequestSet := d.GetOk("backend_request")
	backendResponse, backendResponseSet := d.GetOk("backend_response")
	if backendRequestSet || backendResponseSet {
		parameters.Properties.Backend = &apidiagnostic.PipelineDiagnosticSettings{}
		if backendRequestSet {
			parameters.Properties.Backend.Request = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(backendRequest.([]interface{}))
		}
		if backendResponseSet {
			parameters.Properties.Backend.Response = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(backendResponse.([]interface{}))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apidiagnostic.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating Diagnostic %s: %+v", id, err)
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving Diagnostic %s: %+v", id, err)
	}
	if resp.Model != nil && pointer.From(resp.Model.Id) == "" {
		return fmt.Errorf("reading ID for Diagnostic %s: ID is empty", id)
	}
	d.SetId(id.ID())

	return resourceApiManagementApiDiagnosticRead(d, meta)
}

func resourceApiManagementApiDiagnosticRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := apidiagnostic.ParseApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(diagnosticId.ApiId)

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
				format = string(pointer.From(props.OperationNameFormat))
			}
			d.Set("operation_name_format", format)
		}
	}

	return nil
}

func resourceApiManagementApiDiagnosticDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := apidiagnostic.ParseApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(diagnosticId.ApiId)

	newId := apidiagnostic.NewApiDiagnosticID(diagnosticId.SubscriptionId, diagnosticId.ResourceGroupName, diagnosticId.ServiceName, name, diagnosticId.DiagnosticId)
	if resp, err := client.Delete(ctx, newId, apidiagnostic.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}

func expandApiManagementApiDiagnosticHTTPMessageDiagnostic(input []interface{}) *apidiagnostic.HTTPMessageDiagnostic {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &apidiagnostic.HTTPMessageDiagnostic{
		Body: &apidiagnostic.BodyDiagnosticSettings{},
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

	result.DataMasking = expandApiManagementApiDiagnosticDataMasking(v["data_masking"].([]interface{}))

	return result
}

func flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(input *apidiagnostic.HTTPMessageDiagnostic) []interface{} {
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

	diagnostic["data_masking"] = flattenApiManagementApiDiagnosticDataMasking(input.DataMasking)

	result = append(result, diagnostic)

	return result
}

func schemaApiManagementDataMaskingEntityList() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mode": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(apidiagnostic.DataMaskingModeHide),
						string(apidiagnostic.DataMaskingModeMask),
					}, false),
				},

				"value": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func expandApiManagementApiDiagnosticDataMasking(input []interface{}) *apidiagnostic.DataMasking {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	inputRaw := input[0].(map[string]interface{})
	return &apidiagnostic.DataMasking{
		QueryParams: expandApiManagementApiDiagnosticDataMaskingEntityList(inputRaw["query_params"].([]interface{})),
		Headers:     expandApiManagementApiDiagnosticDataMaskingEntityList(inputRaw["headers"].([]interface{})),
	}
}

func expandApiManagementApiDiagnosticDataMaskingEntityList(input []interface{}) *[]apidiagnostic.DataMaskingEntity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]apidiagnostic.DataMaskingEntity, 0)
	for _, v := range input {
		entity := v.(map[string]interface{})
		result = append(result, apidiagnostic.DataMaskingEntity{
			Mode:  pointer.To(apidiagnostic.DataMaskingMode(entity["mode"].(string))),
			Value: pointer.To(entity["value"].(string)),
		})
	}
	return &result
}

func flattenApiManagementApiDiagnosticDataMasking(dataMasking *apidiagnostic.DataMasking) []interface{} {
	if dataMasking == nil {
		return []interface{}{}
	}

	var queryParams, headers []interface{}
	if dataMasking.QueryParams != nil {
		queryParams = flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMasking.QueryParams)
	}
	if dataMasking.Headers != nil {
		headers = flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMasking.Headers)
	}

	return []interface{}{
		map[string]interface{}{
			"query_params": queryParams,
			"headers":      headers,
		},
	}
}

func flattenApiManagementApiDiagnosticDataMaskingEntityList(dataMaskingList *[]apidiagnostic.DataMaskingEntity) []interface{} {
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
