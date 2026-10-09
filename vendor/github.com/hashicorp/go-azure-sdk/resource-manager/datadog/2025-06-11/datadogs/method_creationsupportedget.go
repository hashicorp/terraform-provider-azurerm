package datadogs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreationSupportedGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *CreateResourceSupportedResponse
}

type CreationSupportedGetOperationOptions struct {
	DatadogOrganizationId *string
}

func DefaultCreationSupportedGetOperationOptions() CreationSupportedGetOperationOptions {
	return CreationSupportedGetOperationOptions{}
}

func (o CreationSupportedGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CreationSupportedGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CreationSupportedGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DatadogOrganizationId != nil {
		out.Append("datadogOrganizationId", fmt.Sprintf("%v", *o.DatadogOrganizationId))
	}
	return &out
}

// CreationSupportedGet ...
func (c DatadogsClient) CreationSupportedGet(ctx context.Context, id commonids.SubscriptionId, options CreationSupportedGetOperationOptions) (result CreationSupportedGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Datadog/subscriptionStatuses/default", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
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

	var model CreateResourceSupportedResponse
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
