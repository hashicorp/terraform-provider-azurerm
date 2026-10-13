// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apioperation.NewOperationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("operation_id").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	method := d.Get("method").(string)
	urlTemplate := d.Get("url_template").(string)

	requestContractRaw := d.Get("request").([]any)
	requestContract, err := expandApiManagementOperationRequestContract(d, "request", requestContractRaw)
	if err != nil {
		return err
	}

	responseContractsRaw := d.Get("response").([]any)
	responseContracts, err := expandApiManagementOperationResponseContract(d, "response", responseContractsRaw)
	if err != nil {
		return err
	}

	templateParametersRaw := d.Get("template_parameter").([]any)
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

func expandApiManagementOperationRequestContract(d *pluginsdk.ResourceData, schemaPath string, input []any) (*apioperation.RequestContract, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	vs := input[0].(map[string]any)
	if vs == nil {
		return nil, nil
	}
	description := vs["description"].(string)

	headersRaw := vs["header"].([]any)
	if headersRaw == nil {
		headersRaw = []any{}
	}
	headers := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.0.header", schemaPath), headersRaw)

	queryParametersRaw := vs["query_parameter"].([]any)
	if queryParametersRaw == nil {
		queryParametersRaw = []any{}
	}
	queryParameters := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.0.query_parameter", schemaPath), queryParametersRaw)

	representationsRaw := vs["representation"].([]any)
	if representationsRaw == nil {
		representationsRaw = []any{}
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

func expandApiManagementOperationResponseContract(d *pluginsdk.ResourceData, schemaPath string, input []any) (*[]apioperation.ResponseContract, error) {
	if len(input) == 0 {
		return &[]apioperation.ResponseContract{}, nil
	}

	outputs := make([]apioperation.ResponseContract, 0)

	for i, v := range input {
		vs := v.(map[string]any)

		description := vs["description"].(string)
		statusCode := vs["status_code"].(int)

		headersRaw := vs["header"].([]any)
		headers := schemaz.ExpandApiManagementOperationParameterContract(d, fmt.Sprintf("%s.%d.header", schemaPath, i), headersRaw)

		representationsRaw := vs["representation"].([]any)
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
