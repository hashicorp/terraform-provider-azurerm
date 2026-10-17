// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperation.ParseOperationID(d.Id())
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(id.ApiId)

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

func flattenApiManagementOperationRequestContract(input *apioperation.RequestContract) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	output := make(map[string]any)

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

	return []any{output}, nil
}

func flattenApiManagementOperationResponseContract(input *[]apioperation.ResponseContract) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	outputs := make([]any, 0)

	for _, v := range *input {
		output := make(map[string]any)

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
