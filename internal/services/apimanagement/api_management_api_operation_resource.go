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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperation() *pluginsdk.Resource {
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

		Schema: map[string]*pluginsdk.Schema{
			"operation_id": schemaz.SchemaApiManagementChildName(),

			"api_name": schemaz.SchemaApiManagementApiName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apioperation.NewOperationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("operation_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation", id.ID())
		}
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	method := d.Get("method").(string)
	urlTemplate := d.Get("url_template").(string)

	requestContractRaw := d.Get("request").([]interface{})
	requestContract, err := expandApiManagementOperationRequestContract(d, "request", requestContractRaw)
	if err != nil {
		return err
	}

	responseContractsRaw := d.Get("response").([]interface{})
	responseContracts, err := expandApiManagementOperationResponseContract(d, "response", responseContractsRaw)
	if err != nil {
		return err
	}

	templateParametersRaw := d.Get("template_parameter").([]interface{})
	templateParameters := schemaz.ExpandApiManagementOperationParameterContract(d, "template_parameter", templateParametersRaw)

	parameters := apioperation.OperationContract{
		Properties: &apioperation.OperationContractProperties{
			Description:        pointer.To(description),
			DisplayName:        displayName,
			Method:             method,
			Request:            requestContract,
			Responses:          responseContracts,
			TemplateParameters: templateParameters,
			UrlTemplate:        urlTemplate,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, apioperation.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementApiOperationRead(d, meta)
}

func resourceApiManagementApiOperationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperation.ParseOperationID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)

	newId := apioperation.NewOperationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.OperationId)
	resp, err := client.Get(ctx, newId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", newId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	d.Set("operation_id", id.OperationId)
	d.Set("api_name", apiName)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("method", props.Method)
			d.Set("url_template", props.UrlTemplate)

			flattenedRequest, err := flattenApiManagementOperationRequestContract(props.Request)
			if err != nil {
				return err
			}
			if err := d.Set("request", flattenedRequest); err != nil {
				return fmt.Errorf("flattening `request`: %+v", err)
			}

			flattenedResponse, err := flattenApiManagementOperationResponseContract(props.Responses)
			if err != nil {
				return err
			}
			if err := d.Set("response", flattenedResponse); err != nil {
				return fmt.Errorf("flattening `response`: %+v", err)
			}

			flattenedTemplateParams, err := schemaz.FlattenApiManagementOperationParameterContract(props.TemplateParameters)
			if err != nil {
				return err
			}

			if err := d.Set("template_parameter", flattenedTemplateParams); err != nil {
				return fmt.Errorf("flattening `template_parameter`: %+v", err)
			}
		}
	}
	return nil
}

func resourceApiManagementApiOperationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperation.ParseOperationID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)

	newId := apioperation.NewOperationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name, id.OperationId)
	resp, err := client.Delete(ctx, newId, apioperation.DeleteOperationOptions{})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}

func expandApiManagementOperationRequestContract(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) (*apioperation.RequestContract, error) {
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
	headers := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.0.header", schemaPath), headersRaw)

	queryParametersRaw := vs["query_parameter"].([]interface{})
	if queryParametersRaw == nil {
		queryParametersRaw = []interface{}{}
	}
	queryParameters := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.0.query_parameter", schemaPath), queryParametersRaw)

	representationsRaw := vs["representation"].([]interface{})
	if representationsRaw == nil {
		representationsRaw = []interface{}{}
	}
	representations, err := schemaz.ExpandApiManagementOperationRepresentation(d, fmt.Sprintf("%s.0.representation", schemaPath), representationsRaw)
	if err != nil {
		return nil, err
	}

	return &apioperation.RequestContract{
		Description:     pointer.To(description),
		Headers:         headers,
		QueryParameters: queryParameters,
		Representations: representations,
	}, nil
}

func flattenApiManagementOperationRequestContract(input *apioperation.RequestContract) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	output := make(map[string]interface{})

	output["description"] = pointer.From(input.Description)

	header, err := schemaz.FlattenApiManagementOperationParameterContract(input.Headers)
	if err != nil {
		return nil, err
	}
	output["header"] = header

	queryParameter, err := schemaz.FlattenApiManagementOperationParameterContract(input.QueryParameters)
	if err != nil {
		return nil, err
	}
	output["query_parameter"] = queryParameter

	representation, err := schemaz.FlattenApiManagementOperationRepresentation(input.Representations)
	if err != nil {
		return nil, err
	}
	output["representation"] = representation

	return []interface{}{output}, nil
}

func expandApiManagementOperationResponseContract(d *pluginsdk.ResourceData, schemaPath string, input []interface{}) (*[]apioperation.ResponseContract, error) {
	if len(input) == 0 {
		return &[]apioperation.ResponseContract{}, nil
	}

	outputs := make([]apioperation.ResponseContract, 0)

	for i, v := range input {
		vs := v.(map[string]interface{})

		description := vs["description"].(string)
		statusCode := vs["status_code"].(int)

		headersRaw := vs["header"].([]interface{})
		headers := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.%d.header", schemaPath, i), headersRaw)

		representationsRaw := vs["representation"].([]interface{})
		representations, err := schemaz.ExpandApiManagementOperationRepresentation(d, fmt.Sprintf("%s.%d.representation", schemaPath, i), representationsRaw)
		if err != nil {
			return nil, err
		}

		output := apioperation.ResponseContract{
			Description:     pointer.To(description),
			Headers:         headers,
			Representations: representations,
			StatusCode:      int64(statusCode),
		}

		outputs = append(outputs, output)
	}

	return &outputs, nil
}

func flattenApiManagementOperationResponseContract(input *[]apioperation.ResponseContract) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	outputs := make([]interface{}, 0)

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Description != nil {
			output["description"] = *v.Description
		}

		output["status_code"] = int(v.StatusCode)

		header, err := schemaz.FlattenApiManagementOperationParameterContract(v.Headers)
		if err != nil {
			return nil, err
		}
		output["header"] = header

		representation, err := schemaz.FlattenApiManagementOperationRepresentation(v.Representations)
		if err != nil {
			return nil, err
		}
		output["representation"] = representation

		outputs = append(outputs, output)
	}

	return outputs, nil
}
