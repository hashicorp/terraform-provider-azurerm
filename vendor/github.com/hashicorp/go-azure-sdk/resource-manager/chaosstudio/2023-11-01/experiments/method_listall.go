package experiments

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

type ListAllOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Experiment
}

type ListAllCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Experiment
}

type ListAllOperationOptions struct {
	ContinuationToken *string
	Running           *bool
}

func DefaultListAllOperationOptions() ListAllOperationOptions {
	return ListAllOperationOptions{}
}

func (o ListAllOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAllOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListAllOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ContinuationToken != nil {
		out.Append("continuationToken", fmt.Sprintf("%v", *o.ContinuationToken))
	}
	if o.Running != nil {
		out.Append("running", fmt.Sprintf("%v", *o.Running))
	}
	return &out
}

// ListAll ...
func (c ExperimentsClient) ListAll(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions) (result ListAllOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/providers/Microsoft.Chaos/experiments", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]Experiment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAllComplete retrieves all the results into a single object
func (c ExperimentsClient) ListAllComplete(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions) (ListAllCompleteResult, error) {
	return c.ListAllCompleteMatchingPredicate(ctx, id, options, ExperimentOperationPredicate{})
}

// ListAllCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ExperimentsClient) ListAllCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListAllOperationOptions, predicate ExperimentOperationPredicate) (result ListAllCompleteResult, err error) {
	items := make([]Experiment, 0)

	resp, err := c.ListAll(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListAllCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
