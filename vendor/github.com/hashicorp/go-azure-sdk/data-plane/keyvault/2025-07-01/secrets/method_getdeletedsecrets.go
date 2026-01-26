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

type GetDeletedSecretsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedSecretItem
}

type GetDeletedSecretsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedSecretItem
}

type GetDeletedSecretsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetDeletedSecretsOperationOptions() GetDeletedSecretsOperationOptions {
	return GetDeletedSecretsOperationOptions{}
}

func (o GetDeletedSecretsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetDeletedSecretsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetDeletedSecretsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetDeletedSecretsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDeletedSecretsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDeletedSecrets ...
func (c SecretsClient) GetDeletedSecrets(ctx context.Context, options GetDeletedSecretsOperationOptions) (result GetDeletedSecretsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetDeletedSecretsCustomPager{},
		Path:          "/deletedsecrets",
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
		Values *[]DeletedSecretItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDeletedSecretsComplete retrieves all the results into a single object
func (c SecretsClient) GetDeletedSecretsComplete(ctx context.Context, options GetDeletedSecretsOperationOptions) (GetDeletedSecretsCompleteResult, error) {
	return c.GetDeletedSecretsCompleteMatchingPredicate(ctx, options, DeletedSecretItemOperationPredicate{})
}

// GetDeletedSecretsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SecretsClient) GetDeletedSecretsCompleteMatchingPredicate(ctx context.Context, options GetDeletedSecretsOperationOptions, predicate DeletedSecretItemOperationPredicate) (result GetDeletedSecretsCompleteResult, err error) {
	items := make([]DeletedSecretItem, 0)

	resp, err := c.GetDeletedSecrets(ctx, options)
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

	result = GetDeletedSecretsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
