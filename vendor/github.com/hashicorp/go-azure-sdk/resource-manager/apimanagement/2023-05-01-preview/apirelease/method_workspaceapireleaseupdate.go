package apirelease

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiReleaseUpdateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApiReleaseContract
}

type WorkspaceApiReleaseUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiReleaseUpdateOperationOptions() WorkspaceApiReleaseUpdateOperationOptions {
	return WorkspaceApiReleaseUpdateOperationOptions{}
}

func (o WorkspaceApiReleaseUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiReleaseUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiReleaseUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiReleaseUpdate ...
func (c ApiReleaseClient) WorkspaceApiReleaseUpdate(ctx context.Context, id ApiReleaseId, input ApiReleaseContract, options WorkspaceApiReleaseUpdateOperationOptions) (result WorkspaceApiReleaseUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPatch,
		Path:          id.ID(),
		OptionsObject: options,
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

	var model ApiReleaseContract
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
