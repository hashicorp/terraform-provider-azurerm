package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApiContract
}

type WorkspaceApiCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiCreateOrUpdateOperationOptions() WorkspaceApiCreateOrUpdateOperationOptions {
	return WorkspaceApiCreateOrUpdateOperationOptions{}
}

func (o WorkspaceApiCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiCreateOrUpdate ...
func (c ApiClient) WorkspaceApiCreateOrUpdate(ctx context.Context, id WorkspaceApiId, input ApiCreateOrUpdateParameter, options WorkspaceApiCreateOrUpdateOperationOptions) (result WorkspaceApiCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// WorkspaceApiCreateOrUpdateThenPoll performs WorkspaceApiCreateOrUpdate then polls until it's completed
func (c ApiClient) WorkspaceApiCreateOrUpdateThenPoll(ctx context.Context, id WorkspaceApiId, input ApiCreateOrUpdateParameter, options WorkspaceApiCreateOrUpdateOperationOptions) error {
	result, err := c.WorkspaceApiCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing WorkspaceApiCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WorkspaceApiCreateOrUpdate: %+v", err)
	}

	return nil
}
