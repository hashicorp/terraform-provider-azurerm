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

type GetSecretVersionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecretItem
}

type GetSecretVersionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecretItem
}

type GetSecretVersionsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetSecretVersionsOperationOptions() GetSecretVersionsOperationOptions {
	return GetSecretVersionsOperationOptions{}
}

func (o GetSecretVersionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetSecretVersionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetSecretVersionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetSecretVersionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetSecretVersionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetSecretVersions ...
func (c SecretsClient) GetSecretVersions(ctx context.Context, id SecretId, options GetSecretVersionsOperationOptions) (result GetSecretVersionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetSecretVersionsCustomPager{},
		Path:          fmt.Sprintf("%s/versions", id.Path()),
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

// GetSecretVersionsComplete retrieves all the results into a single object
func (c SecretsClient) GetSecretVersionsComplete(ctx context.Context, id SecretId, options GetSecretVersionsOperationOptions) (GetSecretVersionsCompleteResult, error) {
	return c.GetSecretVersionsCompleteMatchingPredicate(ctx, id, options, SecretItemOperationPredicate{})
}

// GetSecretVersionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SecretsClient) GetSecretVersionsCompleteMatchingPredicate(ctx context.Context, id SecretId, options GetSecretVersionsOperationOptions, predicate SecretItemOperationPredicate) (result GetSecretVersionsCompleteResult, err error) {
	items := make([]SecretItem, 0)

	resp, err := c.GetSecretVersions(ctx, id, options)
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

	result = GetSecretVersionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
