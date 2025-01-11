package streamingjobs

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

type CreateOrReplaceOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *StreamingJob
}

type CreateOrReplaceOperationOptions struct {
	IfMatch     *string
	IfNoneMatch *string
}

func DefaultCreateOrReplaceOperationOptions() CreateOrReplaceOperationOptions {
	return CreateOrReplaceOperationOptions{}
}

func (o CreateOrReplaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	if o.IfNoneMatch != nil {
		out.Append("If-None-Match", fmt.Sprintf("%v", *o.IfNoneMatch))
	}
	return &out
}

func (o CreateOrReplaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CreateOrReplaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// CreateOrReplace ...
func (c StreamingJobsClient) CreateOrReplace(ctx context.Context, id StreamingJobId, input StreamingJob, options CreateOrReplaceOperationOptions) (result CreateOrReplaceOperationResponse, err error) {
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// CreateOrReplaceThenPoll performs CreateOrReplace then polls until it's completed
func (c StreamingJobsClient) CreateOrReplaceThenPoll(ctx context.Context, id StreamingJobId, input StreamingJob, options CreateOrReplaceOperationOptions) error {
	result, err := c.CreateOrReplace(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing CreateOrReplace: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CreateOrReplace: %+v", err)
	}

	return nil
}
