// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// This `azuresdkhack` only exists because the API Specification for `Microsoft.Web/connections` (2016-06-01) does not define
// the `parameterValueSet` and `parameterValueType` properties, which are supported by the API and are required to configure
// API Connections that authenticate using a Managed Identity.
// If these properties are added to the API Specification, this can be removed.

type ConnectionsClient struct {
	client *connections.ConnectionsClient
}

func NewConnectionsWorkaroundClient(client *connections.ConnectionsClient) ConnectionsClient {
	return ConnectionsClient{
		client: client,
	}
}

type GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApiConnectionDefinition
}

func (c ConnectionsClient) Get(ctx context.Context, id connections.ConnectionId) (result GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model ApiConnectionDefinition
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

type CreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApiConnectionDefinition
}

func (c ConnectionsClient) CreateOrUpdate(ctx context.Context, id connections.ConnectionId, input ApiConnectionDefinition) (result CreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.client.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model ApiConnectionDefinition
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}

type ApiConnectionDefinition struct {
	Etag       *string                            `json:"etag,omitempty"`
	Id         *string                            `json:"id,omitempty"`
	Kind       *string                            `json:"kind,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Properties *ApiConnectionDefinitionProperties `json:"properties,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}

type ApiConnectionDefinitionProperties struct {
	Api                      *connections.ApiReference                 `json:"api,omitempty"`
	ChangedTime              *string                                   `json:"changedTime,omitempty"`
	CreatedTime              *string                                   `json:"createdTime,omitempty"`
	CustomParameterValues    *map[string]any                           `json:"customParameterValues,omitempty"`
	DisplayName              *string                                   `json:"displayName,omitempty"`
	NonSecretParameterValues *map[string]any                           `json:"nonSecretParameterValues,omitempty"`
	ParameterValues          *map[string]any                           `json:"parameterValues,omitempty"`
	Statuses                 *[]connections.ConnectionStatusDefinition `json:"statuses,omitempty"`
	TestLinks                *[]connections.ApiConnectionTestLink      `json:"testLinks,omitempty"`

	// The following properties are not defined in the API Specification
	ParameterValueSet  *ParameterValueSet `json:"parameterValueSet,omitempty"`
	ParameterValueType *string            `json:"parameterValueType,omitempty"`
}

type ParameterValueSet struct {
	Name   *string         `json:"name,omitempty"`
	Values *map[string]any `json:"values,omitempty"`
}
