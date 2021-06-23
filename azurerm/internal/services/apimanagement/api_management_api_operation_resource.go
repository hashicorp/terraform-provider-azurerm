package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiOperation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiOperationCreateUpdate,
		Read:   resourceApiManagementApiOperationRead,
		Update: resourceApiManagementApiOperationCreateUpdate,
		Delete: resourceApiManagementApiOperationDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"operation_id": schemaz.SchemaApiManagementChildName(),

			"api_name": schemaz.SchemaApiManagementApiName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"method": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"url_template": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"request": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"header": schemaz.SchemaApiManagementOperationParameterContract(),

						"query_parameter": schemaz.SchemaApiManagementOperationParameterContract(),

						"representation": schemaz.SchemaApiManagementOperationRepresentation(),
					},
				},
			},

			"response": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"status_code": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"header": schemaz.SchemaApiManagementOperationParameterContract(),

						"representation": schemaz.SchemaApiManagementOperationRepresentation(),
					},
				},
			},

			"template_parameter": schemaz.SchemaApiManagementOperationParameterContract(),
		},
	}
}

func resourceApiManagementApiOperationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiId := d.Get("api_name").(string)
	operationId := d.Get("operation_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiId, operationId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Operation %q (API %q / API Management Service %q / Resource Group %q): %s", operationId, apiId, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	method := d.Get("method").(string)
	urlTemplate := d.Get("url_template").(string)

	requestContractRaw := d.Get("request").([]interface{})
	requestContract, err := expandApiManagementOperationRequestContract(requestContractRaw)
	if err != nil {
		return err
	}

	responseContractsRaw := d.Get("response").([]interface{})
	responseContracts, err := expandApiManagementOperationResponseContract(responseContractsRaw)
	if err != nil {
		return err
	}

	templateParametersRaw := d.Get("template_parameter").([]interface{})
	templateParameters := schemaz.ExpandApiManagementOperationParameterContract(templateParametersRaw)

	parameters := apimanagement.OperationContract{
		OperationContractProperties: &apimanagement.OperationContractProperties{
			Description:        utils.String(description),
			DisplayName:        utils.String(displayName),
			Method:             utils.String(method),
			Request:            requestContract,
			Responses:          responseContracts,
			TemplateParameters: templateParameters,
			URLTemplate:        utils.String(urlTemplate),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiId, operationId, parameters, ""); err != nil {
		return fmt.Errorf("creating/updating API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiId, operationId)
	if err != nil {
		return fmt.Errorf("retrieving API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceApiManagementApiOperationRead(d, meta)
}

func resourceApiManagementApiOperationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiOperationID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiId := id.ApiName
	operationId := id.OperationName

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiId, operationId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Operation %q (API %q / API Management Service %q / Resource Group %q) was not found - removing from state!", operationId, apiId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
	}

	d.Set("operation_id", operationId)
	d.Set("api_name", apiId)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.OperationContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("method", props.Method)
		d.Set("url_template", props.URLTemplate)

		flattenedRequest := flattenApiManagementOperationRequestContract(props.Request)
		if err := d.Set("request", flattenedRequest); err != nil {
			return fmt.Errorf("flattening `request`: %+v", err)
		}

		flattenedResponse := flattenApiManagementOperationResponseContract(props.Responses)
		if err := d.Set("response", flattenedResponse); err != nil {
			return fmt.Errorf("flattening `response`: %+v", err)
		}

		flattenedTemplateParams := schemaz.FlattenApiManagementOperationParameterContract(props.TemplateParameters)
		if err := d.Set("template_parameter", flattenedTemplateParams); err != nil {
			return fmt.Errorf("flattening `template_parameter`: %+v", err)
		}
	}

	return nil
}

func resourceApiManagementApiOperationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiOperationID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiId := id.ApiName
	operationId := id.OperationName

	resp, err := client.Delete(ctx, resourceGroup, serviceName, apiId, operationId, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting API Operation %q (API %q / API Management Service %q / Resource Group %q): %+v", operationId, apiId, serviceName, resourceGroup, err)
		}
	}

	return nil
}

func expandApiManagementOperationRequestContract(input []interface{}) (*apimanagement.RequestContract, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	vs := input[0].(map[string]interface{})
	if vs == nil {
		return nil, nil
	}
	description := vs["description"].(string)

	headersRaw := vs["header"].([]interface{})
	if headersRaw == nil {
		headersRaw = []interface{}{}
	}
	headers := schemaz.ExpandApiManagementOperationParameterContract(headersRaw)

	queryParametersRaw := vs["query_parameter"].([]interface{})
	if queryParametersRaw == nil {
		queryParametersRaw = []interface{}{}
	}
	queryParameters := schemaz.ExpandApiManagementOperationParameterContract(queryParametersRaw)

	representationsRaw := vs["representation"].([]interface{})
	if representationsRaw == nil {
		representationsRaw = []interface{}{}
	}
	representations, err := schemaz.ExpandApiManagementOperationRepresentation(representationsRaw)
	if err != nil {
		return nil, err
	}

	return &apimanagement.RequestContract{
		Description:     utils.String(description),
		Headers:         headers,
		QueryParameters: queryParameters,
		Representations: representations,
	}, nil
}

func flattenApiManagementOperationRequestContract(input *apimanagement.RequestContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if input.Description != nil {
		output["description"] = *input.Description
	}

	output["header"] = schemaz.FlattenApiManagementOperationParameterContract(input.Headers)
	output["query_parameter"] = schemaz.FlattenApiManagementOperationParameterContract(input.QueryParameters)
	output["representation"] = schemaz.FlattenApiManagementOperationRepresentation(input.Representations)

	return []interface{}{output}
}

func expandApiManagementOperationResponseContract(input []interface{}) (*[]apimanagement.ResponseContract, error) {
	if len(input) == 0 {
		return &[]apimanagement.ResponseContract{}, nil
	}

	outputs := make([]apimanagement.ResponseContract, 0)

	for _, v := range input {
		vs := v.(map[string]interface{})

		description := vs["description"].(string)
		statusCode := vs["status_code"].(int)

		headersRaw := vs["header"].([]interface{})
		headers := schemaz.ExpandApiManagementOperationParameterContract(headersRaw)

		representationsRaw := vs["representation"].([]interface{})
		representations, err := schemaz.ExpandApiManagementOperationRepresentation(representationsRaw)
		if err != nil {
			return nil, err
		}

		output := apimanagement.ResponseContract{
			Description:     utils.String(description),
			Headers:         headers,
			Representations: representations,
			StatusCode:      utils.Int32(int32(statusCode)),
		}

		outputs = append(outputs, output)
	}

	return &outputs, nil
}

func flattenApiManagementOperationResponseContract(input *[]apimanagement.ResponseContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Description != nil {
			output["description"] = *v.Description
		}

		if v.StatusCode != nil {
			output["status_code"] = int(*v.StatusCode)
		}

		output["header"] = schemaz.FlattenApiManagementOperationParameterContract(v.Headers)
		output["representation"] = schemaz.FlattenApiManagementOperationRepresentation(v.Representations)

		outputs = append(outputs, output)
	}

	return outputs
}
