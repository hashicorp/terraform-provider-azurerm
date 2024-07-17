package apiversionset

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiVersionSetCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApiVersionSetContract
}

type WorkspaceApiVersionSetCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions() WorkspaceApiVersionSetCreateOrUpdateOperationOptions {
	return WorkspaceApiVersionSetCreateOrUpdateOperationOptions{}
}

func (o WorkspaceApiVersionSetCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiVersionSetCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiVersionSetCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiVersionSetCreateOrUpdate ...
func (c ApiVersionSetClient) WorkspaceApiVersionSetCreateOrUpdate(ctx context.Context, id WorkspaceApiVersionSetId, input ApiVersionSetContract, options WorkspaceApiVersionSetCreateOrUpdateOperationOptions) (result WorkspaceApiVersionSetCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
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

	var model ApiVersionSetContract
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
