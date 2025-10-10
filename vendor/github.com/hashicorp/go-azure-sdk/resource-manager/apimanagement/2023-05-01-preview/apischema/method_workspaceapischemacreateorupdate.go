package apischema

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

type WorkspaceApiSchemaCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SchemaContract
}

type WorkspaceApiSchemaCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiSchemaCreateOrUpdateOperationOptions() WorkspaceApiSchemaCreateOrUpdateOperationOptions {
	return WorkspaceApiSchemaCreateOrUpdateOperationOptions{}
}

func (o WorkspaceApiSchemaCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiSchemaCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiSchemaCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiSchemaCreateOrUpdate ...
func (c ApiSchemaClient) WorkspaceApiSchemaCreateOrUpdate(ctx context.Context, id WorkspaceApiSchemaId, input SchemaContract, options WorkspaceApiSchemaCreateOrUpdateOperationOptions) (result WorkspaceApiSchemaCreateOrUpdateOperationResponse, err error) {
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

// WorkspaceApiSchemaCreateOrUpdateThenPoll performs WorkspaceApiSchemaCreateOrUpdate then polls until it's completed
func (c ApiSchemaClient) WorkspaceApiSchemaCreateOrUpdateThenPoll(ctx context.Context, id WorkspaceApiSchemaId, input SchemaContract, options WorkspaceApiSchemaCreateOrUpdateOperationOptions) error {
	result, err := c.WorkspaceApiSchemaCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing WorkspaceApiSchemaCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WorkspaceApiSchemaCreateOrUpdate: %+v", err)
	}

	return nil
}
