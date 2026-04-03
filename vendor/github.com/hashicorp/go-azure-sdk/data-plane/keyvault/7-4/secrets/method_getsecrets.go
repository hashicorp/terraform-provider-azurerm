package secrets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSecretsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecretItem
}

type GetSecretsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecretItem
}

type GetSecretsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetSecretsOperationOptions() GetSecretsOperationOptions {
	return GetSecretsOperationOptions{}
}

func (o GetSecretsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetSecretsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetSecretsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetSecretsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetSecretsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetSecrets ...
func (c SecretsClient) GetSecrets(ctx context.Context, options GetSecretsOperationOptions) (result GetSecretsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetSecretsCustomPager{},
		Path:          "/secrets",
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
		Values *[]SecretItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetSecretsComplete retrieves all the results into a single object
func (c SecretsClient) GetSecretsComplete(ctx context.Context, options GetSecretsOperationOptions) (GetSecretsCompleteResult, error) {
	return c.GetSecretsCompleteMatchingPredicate(ctx, options, SecretItemOperationPredicate{})
}

// GetSecretsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SecretsClient) GetSecretsCompleteMatchingPredicate(ctx context.Context, options GetSecretsOperationOptions, predicate SecretItemOperationPredicate) (result GetSecretsCompleteResult, err error) {
	items := make([]SecretItem, 0)

	resp, err := c.GetSecrets(ctx, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = GetSecretsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
