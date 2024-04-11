package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Database
}

type UpdateOperationOptions struct {
	CallerRole *CallerRole
}

func DefaultUpdateOperationOptions() UpdateOperationOptions {
	return UpdateOperationOptions{}
}

func (o UpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o UpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o UpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CallerRole != nil {
		out.Append("callerRole", fmt.Sprintf("%v", *o.CallerRole))
	}
	return &out
}

// Update ...
func (c DatabasesClient) Update(ctx context.Context, id commonids.KustoDatabaseId, input Database, options UpdateOperationOptions) (result UpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusCreated,
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

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// UpdateThenPoll performs Update then polls until it's completed
func (c DatabasesClient) UpdateThenPoll(ctx context.Context, id commonids.KustoDatabaseId, input Database, options UpdateOperationOptions) error {
	result, err := c.Update(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing Update: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Update: %+v", err)
	}

	return nil
}
