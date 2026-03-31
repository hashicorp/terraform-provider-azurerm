package managementgroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionsCreateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SubscriptionUnderManagementGroup
}

type SubscriptionsCreateOperationOptions struct {
	CacheControl *string
}

func DefaultSubscriptionsCreateOperationOptions() SubscriptionsCreateOperationOptions {
	return SubscriptionsCreateOperationOptions{}
}

func (o SubscriptionsCreateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.CacheControl != nil {
		out.Append("Cache-Control", fmt.Sprintf("%v", *o.CacheControl))
	}
	return &out
}

func (o SubscriptionsCreateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SubscriptionsCreateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// SubscriptionsCreate ...
func (c ManagementGroupsClient) SubscriptionsCreate(ctx context.Context, id SubscriptionId, options SubscriptionsCreateOperationOptions) (result SubscriptionsCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model SubscriptionUnderManagementGroup
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
