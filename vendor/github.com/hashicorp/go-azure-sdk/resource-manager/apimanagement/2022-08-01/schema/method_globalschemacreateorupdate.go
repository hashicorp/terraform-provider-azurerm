package schema

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

type GlobalSchemaCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

type GlobalSchemaCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultGlobalSchemaCreateOrUpdateOperationOptions() GlobalSchemaCreateOrUpdateOperationOptions {
	return GlobalSchemaCreateOrUpdateOperationOptions{}
}

func (o GlobalSchemaCreateOrUpdateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o GlobalSchemaCreateOrUpdateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GlobalSchemaCreateOrUpdateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// GlobalSchemaCreateOrUpdate ...
func (c SchemaClient) GlobalSchemaCreateOrUpdate(ctx context.Context, id SchemaId, input GlobalSchemaContract, options GlobalSchemaCreateOrUpdateOperationOptions) (result GlobalSchemaCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
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

// GlobalSchemaCreateOrUpdateThenPoll performs GlobalSchemaCreateOrUpdate then polls until it's completed
func (c SchemaClient) GlobalSchemaCreateOrUpdateThenPoll(ctx context.Context, id SchemaId, input GlobalSchemaContract, options GlobalSchemaCreateOrUpdateOperationOptions) error {
	result, err := c.GlobalSchemaCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing GlobalSchemaCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GlobalSchemaCreateOrUpdate: %+v", err)
	}

	return nil
}
