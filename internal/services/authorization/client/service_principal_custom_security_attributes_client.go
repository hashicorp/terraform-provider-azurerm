// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/msgraph"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

const customSecurityAttributeValueODataType = "#Microsoft.DirectoryServices.CustomSecurityAttributeValue"

type graphRequestOptions struct {
	query odata.Query
}

func (o graphRequestOptions) ToHeaders() *client.Headers {
	h := client.Headers{}
	h.AppendHeader(o.query.Headers())
	return &h
}

func (o graphRequestOptions) ToOData() *odata.Query {
	return &o.query
}

func (o graphRequestOptions) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.AppendValues(o.query.Values())
	return &q
}

type ServicePrincipalCustomSecurityAttributesClient struct {
	msGraphClient *msgraph.Client
}

type servicePrincipalModel struct {
	ID                       *string                           `json:"id"`
	CustomSecurityAttributes map[string]map[string]interface{} `json:"customSecurityAttributes"`
}

type servicePrincipalPatchModel struct {
	CustomSecurityAttributes map[string]map[string]interface{} `json:"customSecurityAttributes"`
}

func NewServicePrincipalCustomSecurityAttributesClient(o *common.ClientOptions) (*ServicePrincipalCustomSecurityAttributesClient, error) {
	graphClient, err := msgraph.NewClient(o.Environment.MicrosoftGraph, "ServicePrincipalCustomSecurityAttributes", msgraph.VersionOnePointZero)
	if err != nil {
		return nil, fmt.Errorf("building microsoft graph client: %+v", err)
	}

	graphAuthorizer, err := o.Authorizers.AuthorizerFunc(o.Environment.MicrosoftGraph)
	if err != nil {
		return nil, fmt.Errorf("building microsoft graph authorizer: %+v", err)
	}
	o.Configure(graphClient, graphAuthorizer)

	return &ServicePrincipalCustomSecurityAttributesClient{
		msGraphClient: graphClient,
	}, nil
}

func (c *ServicePrincipalCustomSecurityAttributesClient) Get(ctx context.Context, servicePrincipalObjectId string) (*servicePrincipalModel, error) {
	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusNotFound,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: graphRequestOptions{
			query: odata.Query{
				Select: []string{"id", "customSecurityAttributes"},
			},
		},
		Path: fmt.Sprintf("/servicePrincipals/%s", servicePrincipalObjectId),
	}

	req, err := c.msGraphClient.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building request to get service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	model := servicePrincipalModel{}
	if err := resp.Unmarshal(&model); err != nil {
		return nil, fmt.Errorf("unmarshaling service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}

	if model.CustomSecurityAttributes == nil {
		model.CustomSecurityAttributes = map[string]map[string]interface{}{}
	}

	return &model, nil
}

func (c *ServicePrincipalCustomSecurityAttributesClient) Update(ctx context.Context, servicePrincipalObjectId string, customSecurityAttributes map[string]map[string]interface{}) error {
	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod:    http.MethodPatch,
		OptionsObject: nil,
		Path:          fmt.Sprintf("/servicePrincipals/%s", servicePrincipalObjectId),
	}

	req, err := c.msGraphClient.NewRequest(ctx, opts)
	if err != nil {
		return fmt.Errorf("building request to update service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}

	payload := servicePrincipalPatchModel{
		CustomSecurityAttributes: normalizeCustomSecurityAttributesForPatch(customSecurityAttributes),
	}
	if err := req.Marshal(payload); err != nil {
		return fmt.Errorf("marshaling request payload to update service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}

	if _, err := req.Execute(ctx); err != nil {
		return fmt.Errorf("updating service principal custom attributes for %q: %+v", servicePrincipalObjectId, err)
	}

	return nil
}

func normalizeCustomSecurityAttributesForPatch(input map[string]map[string]interface{}) map[string]map[string]interface{} {
	if input == nil {
		return map[string]map[string]interface{}{}
	}

	output := make(map[string]map[string]interface{}, len(input))
	for setName, attributes := range input {
		normalized := map[string]interface{}{
			"@odata.type": customSecurityAttributeValueODataType,
		}

		for key, value := range attributes {
			if strings.EqualFold(key, "@odata.type") {
				continue
			}
			normalized[key] = value
		}

		output[setName] = normalized
	}

	return output
}
