package updateruns

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

type SkipOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *UpdateRun
}

type SkipOperationOptions struct {
	IfMatch *string
}

func DefaultSkipOperationOptions() SkipOperationOptions {
	return SkipOperationOptions{}
}

func (o SkipOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o SkipOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SkipOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// Skip ...
func (c UpdateRunsClient) Skip(ctx context.Context, id UpdateRunId, input SkipProperties, options SkipOperationOptions) (result SkipOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/skip", id.ID()),
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

// SkipThenPoll performs Skip then polls until it's completed
func (c UpdateRunsClient) SkipThenPoll(ctx context.Context, id UpdateRunId, input SkipProperties, options SkipOperationOptions) error {
	result, err := c.Skip(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing Skip: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Skip: %+v", err)
	}

	return nil
}
