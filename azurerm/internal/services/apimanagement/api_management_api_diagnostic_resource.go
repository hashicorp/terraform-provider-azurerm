package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiDiagnostic() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementApiDiagnosticCreateUpdate,
		Read:   resourceApiManagementApiDiagnosticRead,
		Update: resourceApiManagementApiDiagnosticCreateUpdate,
		Delete: resourceApiManagementApiDiagnosticDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ApiDiagnosticID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"applicationinsights",
					"azuremonitor",
				}, false),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"api_name": azure.SchemaApiManagementApiName(),

			"api_management_logger_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.LoggerID,
			},

			"sampling_percentage": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatBetween(0.0, 100.0),
			},

			"always_log_errors": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"verbosity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.Verbose),
					string(apimanagement.Information),
					string(apimanagement.Error),
				}, false),
			},

			"log_client_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"http_correlation_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.HTTPCorrelationProtocolNone),
					string(apimanagement.HTTPCorrelationProtocolLegacy),
					string(apimanagement.HTTPCorrelationProtocolW3C),
				}, false),
			},

			"frontend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"frontend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_request": resourceApiManagementApiDiagnosticAdditionalContentSchema(),

			"backend_response": resourceApiManagementApiDiagnosticAdditionalContentSchema(),
		},
	}
}

func resourceApiManagementApiDiagnosticAdditionalContentSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"body_bytes": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 8192),
				},
				"headers_to_log": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Set: schema.HashString,
				},
			},
		},
	}
}

func resourceApiManagementApiDiagnosticCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId := d.Get("identifier").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, diagnosticId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %s", diagnosticId, resourceGroup, serviceName, apiName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_diagnostic", *existing.ID)
		}
	}

	parameters := apimanagement.DiagnosticContract{
		DiagnosticContractProperties: &apimanagement.DiagnosticContractProperties{
			LoggerID: utils.String(d.Get("api_management_logger_id").(string)),
		},
	}

	if samplingPercentage, ok := d.GetOk("sampling_percentage"); ok {
		parameters.Sampling = &apimanagement.SamplingSettings{
			SamplingType: apimanagement.Fixed,
			Percentage:   utils.Float(samplingPercentage.(float64)),
		}
	} else {
		parameters.Sampling = nil
	}

	if alwaysLogErrors, ok := d.GetOk("always_log_errors"); ok && alwaysLogErrors.(bool) {
		parameters.AlwaysLog = apimanagement.AllErrors
	}

	if verbosity, ok := d.GetOk("verbosity"); ok {
		switch verbosity.(string) {
		case string(apimanagement.Verbose):
			parameters.Verbosity = apimanagement.Verbose
		case string(apimanagement.Information):
			parameters.Verbosity = apimanagement.Information
		case string(apimanagement.Error):
			parameters.Verbosity = apimanagement.Error
		}
	}

	logClientIP, _ := d.GetOk("log_client_ip")
	if logClientIP != nil {
		parameters.LogClientIP = utils.Bool(logClientIP.(bool))
	}

	if httpCorrelationProtocol, ok := d.GetOk("http_correlation_protocol"); ok {
		switch httpCorrelationProtocol.(string) {
		case string(apimanagement.HTTPCorrelationProtocolNone):
			parameters.HTTPCorrelationProtocol = apimanagement.HTTPCorrelationProtocolNone
		case string(apimanagement.HTTPCorrelationProtocolLegacy):
			parameters.HTTPCorrelationProtocol = apimanagement.HTTPCorrelationProtocolLegacy
		case string(apimanagement.HTTPCorrelationProtocolW3C):
			parameters.HTTPCorrelationProtocol = apimanagement.HTTPCorrelationProtocolW3C
		}
	}

	frontendRequest, frontendRequestSet := d.GetOk("frontend_request")
	frontendResponse, frontendResponseSet := d.GetOk("frontend_response")
	if frontendRequestSet || frontendResponseSet {
		parameters.Frontend = &apimanagement.PipelineDiagnosticSettings{}
		if frontendRequestSet {
			parameters.Frontend.Request = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(frontendRequest.([]interface{}))
		}
		if frontendResponseSet {
			parameters.Frontend.Response = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(frontendResponse.([]interface{}))
		}
	}

	backendRequest, backendRequestSet := d.GetOk("backend_request")
	backendResponse, backendResponseSet := d.GetOk("backend_response")
	if backendRequestSet || backendResponseSet {
		parameters.Backend = &apimanagement.PipelineDiagnosticSettings{}
		if backendRequestSet {
			parameters.Backend.Request = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(backendRequest.([]interface{}))
		}
		if backendResponseSet {
			parameters.Backend.Response = expandApiManagementApiDiagnosticHTTPMessageDiagnostic(backendResponse.([]interface{}))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, diagnosticId, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId, resourceGroup, serviceName, apiName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, diagnosticId)
	if err != nil {
		return fmt.Errorf("retrieving Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId, resourceGroup, serviceName, apiName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading ID for Diagnostic %q (Resource Group %q / API Management Service %q / API %q): ID is empty", diagnosticId, resourceGroup, serviceName, apiName)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementApiDiagnosticRead(d, meta)
}

func resourceApiManagementApiDiagnosticRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := parse.ApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.DiagnosticName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Diagnostic %q (Resource Group %q / API Management Service %q / API %q) was not found - removing from state!", diagnosticId.DiagnosticName, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId.DiagnosticName, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, err)
	}

	d.Set("api_name", diagnosticId.ApiName)
	d.Set("identifier", resp.Name)
	d.Set("resource_group_name", diagnosticId.ResourceGroup)
	d.Set("api_management_name", diagnosticId.ServiceName)
	if props := resp.DiagnosticContractProperties; props != nil {
		d.Set("api_management_logger_id", props.LoggerID)
		if props.Sampling != nil && props.Sampling.Percentage != nil {
			d.Set("sampling_percentage", props.Sampling.Percentage)
		}
		d.Set("always_log_errors", props.AlwaysLog == apimanagement.AllErrors)
		d.Set("verbosity", props.Verbosity)
		d.Set("log_client_ip", props.LogClientIP)
		d.Set("http_correlation_protocol", props.HTTPCorrelationProtocol)
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
	}

	return nil
}

func resourceApiManagementApiDiagnosticDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := parse.ApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.DiagnosticName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId.DiagnosticName, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, err)
		}
	}

	return nil
}

func expandApiManagementApiDiagnosticHTTPMessageDiagnostic(input []interface{}) *apimanagement.HTTPMessageDiagnostic {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := &apimanagement.HTTPMessageDiagnostic{
		Body: &apimanagement.BodyDiagnosticSettings{},
	}

	if bodyBytes, ok := v["body_bytes"]; ok {
		result.Body.Bytes = utils.Int32(int32(bodyBytes.(int)))
	}
	if headersSetRaw, ok := v["headers_to_log"]; ok {
		headersSet := headersSetRaw.(*schema.Set).List()
		headers := []string{}
		for _, header := range headersSet {
			headers = append(headers, header.(string))
		}
		result.Headers = &headers
	}

	return result
}

func flattenApiManagementApiDiagnosticHTTPMessageDiagnostic(input *apimanagement.HTTPMessageDiagnostic) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	diagnostic := map[string]interface{}{}

	if input.Body != nil && input.Body.Bytes != nil {
		diagnostic["body_bytes"] = input.Body.Bytes
	}

	if input.Headers != nil {
		diagnostic["headers_to_log"] = set.FromStringSlice(*input.Headers)
	}
	result = append(result, diagnostic)

	return result
}
