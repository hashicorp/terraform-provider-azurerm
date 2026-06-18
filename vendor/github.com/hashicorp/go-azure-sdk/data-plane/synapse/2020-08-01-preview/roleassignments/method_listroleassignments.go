package roleassignments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRoleAssignmentsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RoleAssignmentDetailsList
}

type ListRoleAssignmentsOperationOptions struct {
	PrincipalId     *string
	RoleId          *string
	Scope           *string
	XMsContinuation *string
}

func DefaultListRoleAssignmentsOperationOptions() ListRoleAssignmentsOperationOptions {
	return ListRoleAssignmentsOperationOptions{}
}

func (o ListRoleAssignmentsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsContinuation != nil {
		out.Append("x-ms-continuation", fmt.Sprintf("%v", *o.XMsContinuation))
	}
	return &out
}

func (o ListRoleAssignmentsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListRoleAssignmentsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.PrincipalId != nil {
		out.Append("principalId", fmt.Sprintf("%v", *o.PrincipalId))
	}
	if o.RoleId != nil {
		out.Append("roleId", fmt.Sprintf("%v", *o.RoleId))
	}
	if o.Scope != nil {
		out.Append("scope", fmt.Sprintf("%v", *o.Scope))
	}
	return &out
}

// ListRoleAssignments ...
func (c RoleAssignmentsClient) ListRoleAssignments(ctx context.Context, options ListRoleAssignmentsOperationOptions) (result ListRoleAssignmentsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/roleAssignments",
	}

	req, err := c.Client.NewRequest(ctx, opts)
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

	var model RoleAssignmentDetailsList
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
