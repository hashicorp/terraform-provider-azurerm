package apioperation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiOperationUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *OperationContract
}

type WorkspaceApiOperationUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiOperationUpdateOperationOptions() WorkspaceApiOperationUpdateOperationOptions {
	return WorkspaceApiOperationUpdateOperationOptions{}
}

func (o WorkspaceApiOperationUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiOperationUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiOperationUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiOperationUpdate ...
func (c ApiOperationClient) WorkspaceApiOperationUpdate(ctx context.Context, id ApiOperationId, input OperationUpdateContract, options WorkspaceApiOperationUpdateOperationOptions) (result WorkspaceApiOperationUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPatch,
		OptionsObject: options,
		Path:          id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
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

	var model OperationContract
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
