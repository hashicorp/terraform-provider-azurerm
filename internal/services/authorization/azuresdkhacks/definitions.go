// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2018-01-01-preview/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RoleDefinitionsWorkaroundClient struct {
	client *resourcemanager.Client
}

func NewRoleDefinitionsWorkaroundClient(resourcemanagerClient *resourcemanager.Client) RoleDefinitionsWorkaroundClient {
	return RoleDefinitionsWorkaroundClient{
		client: resourcemanagerClient,
	}
}

// CreateOrUpdate ...
func (c RoleDefinitionsWorkaroundClient) CreateOrUpdate(ctx context.Context, id roledefinitions.ScopedRoleDefinitionId, input RoleDefinition) (result CreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.client.NewRequest(ctx, opts)
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}

func (c RoleDefinitionsWorkaroundClient) Get(ctx context.Context, id roledefinitions.ScopedRoleDefinitionId) (result GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       id.ID(),
	}

	req, err := c.client.NewRequest(ctx, opts)
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}

type CreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RoleDefinition
}

type GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RoleDefinition
}

type RoleDefinition struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *RoleDefinitionProperties `json:"properties,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}

type RoleDefinitionProperties struct {
	AssignableScopes *[]string                     `json:"assignableScopes,omitempty"`
	Description      *string                       `json:"description,omitempty"`
	Permissions      *[]roledefinitions.Permission `json:"permissions,omitempty"`
	RoleName         *string                       `json:"roleName,omitempty"`
	Type             *string                       `json:"type,omitempty"`
	// not exposed in the sdk
	CreatedOn *string `json:"createdOn,omitempty"`
	UpdatedOn *string `json:"updatedOn,omitempty"`
	CreatedBy *string `json:"createdBy,omitempty"`
	UpdatedBy *string `json:"updatedBy,omitempty"`
}
