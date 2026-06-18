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

type DeleteRoleAssignmentByIdOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type DeleteRoleAssignmentByIdOperationOptions struct {
	Scope *string
}

func DefaultDeleteRoleAssignmentByIdOperationOptions() DeleteRoleAssignmentByIdOperationOptions {
	return DeleteRoleAssignmentByIdOperationOptions{}
}

func (o DeleteRoleAssignmentByIdOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DeleteRoleAssignmentByIdOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DeleteRoleAssignmentByIdOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Scope != nil {
		out.Append("scope", fmt.Sprintf("%v", *o.Scope))
	}
	return &out
}

// DeleteRoleAssignmentById ...
func (c RoleAssignmentsClient) DeleteRoleAssignmentById(ctx context.Context, id RoleAssignmentIdId, options DeleteRoleAssignmentByIdOperationOptions) (result DeleteRoleAssignmentByIdOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: options,
		Path:          id.Path(),
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

	return
}
