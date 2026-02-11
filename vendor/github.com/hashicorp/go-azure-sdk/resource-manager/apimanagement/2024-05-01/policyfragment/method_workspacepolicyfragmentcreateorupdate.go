package policyfragment

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

type WorkspacePolicyFragmentCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PolicyFragmentContract
}

type WorkspacePolicyFragmentCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspacePolicyFragmentCreateOrUpdateOperationOptions() WorkspacePolicyFragmentCreateOrUpdateOperationOptions {
	return WorkspacePolicyFragmentCreateOrUpdateOperationOptions{}
}

func (o WorkspacePolicyFragmentCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspacePolicyFragmentCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspacePolicyFragmentCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspacePolicyFragmentCreateOrUpdate ...
func (c PolicyFragmentClient) WorkspacePolicyFragmentCreateOrUpdate(ctx context.Context, id WorkspacePolicyFragmentId, input PolicyFragmentContract, options WorkspacePolicyFragmentCreateOrUpdateOperationOptions) (result WorkspacePolicyFragmentCreateOrUpdateOperationResponse, err error) {
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

// WorkspacePolicyFragmentCreateOrUpdateThenPoll performs WorkspacePolicyFragmentCreateOrUpdate then polls until it's completed
func (c PolicyFragmentClient) WorkspacePolicyFragmentCreateOrUpdateThenPoll(ctx context.Context, id WorkspacePolicyFragmentId, input PolicyFragmentContract, options WorkspacePolicyFragmentCreateOrUpdateOperationOptions) error {
	result, err := c.WorkspacePolicyFragmentCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing WorkspacePolicyFragmentCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WorkspacePolicyFragmentCreateOrUpdate: %+v", err)
	}

	return nil
}
