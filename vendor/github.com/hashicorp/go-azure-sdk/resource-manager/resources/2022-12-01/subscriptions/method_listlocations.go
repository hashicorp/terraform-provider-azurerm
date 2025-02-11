package subscriptions

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

type ListLocationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *LocationListResult
}

type ListLocationsOperationOptions struct {
	IncludeExtendedLocations *bool
}

func DefaultListLocationsOperationOptions() ListLocationsOperationOptions {
	return ListLocationsOperationOptions{}
}

func (o ListLocationsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListLocationsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListLocationsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludeExtendedLocations != nil {
		out.Append("includeExtendedLocations", fmt.Sprintf("%v", *o.IncludeExtendedLocations))
	}
	return &out
}

// ListLocations ...
func (c SubscriptionsClient) ListLocations(ctx context.Context, id commonids.SubscriptionId, options ListLocationsOperationOptions) (result ListLocationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/locations", id.ID()),
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

	var model LocationListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
