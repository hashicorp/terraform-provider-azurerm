package namedvalue

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

type WorkspaceNamedValueCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NamedValueContract
}

type WorkspaceNamedValueCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions() WorkspaceNamedValueCreateOrUpdateOperationOptions {
	return WorkspaceNamedValueCreateOrUpdateOperationOptions{}
}

func (o WorkspaceNamedValueCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceNamedValueCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceNamedValueCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceNamedValueCreateOrUpdate ...
func (c NamedValueClient) WorkspaceNamedValueCreateOrUpdate(ctx context.Context, id WorkspaceNamedValueId, input NamedValueCreateContract, options WorkspaceNamedValueCreateOrUpdateOperationOptions) (result WorkspaceNamedValueCreateOrUpdateOperationResponse, err error) {
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

// WorkspaceNamedValueCreateOrUpdateThenPoll performs WorkspaceNamedValueCreateOrUpdate then polls until it's completed
func (c NamedValueClient) WorkspaceNamedValueCreateOrUpdateThenPoll(ctx context.Context, id WorkspaceNamedValueId, input NamedValueCreateContract, options WorkspaceNamedValueCreateOrUpdateOperationOptions) error {
	result, err := c.WorkspaceNamedValueCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing WorkspaceNamedValueCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WorkspaceNamedValueCreateOrUpdate: %+v", err)
	}

	return nil
}
