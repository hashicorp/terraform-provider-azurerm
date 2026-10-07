package blobservices

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

type BlobContainersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ListContainerItem
}

type BlobContainersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ListContainerItem
}

type BlobContainersListOperationOptions struct {
	Filter      *string
	Include     *ListContainersInclude
	Maxpagesize *string
}

func DefaultBlobContainersListOperationOptions() BlobContainersListOperationOptions {
	return BlobContainersListOperationOptions{}
}

func (o BlobContainersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o BlobContainersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o BlobContainersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Include != nil {
		out.Append("$include", fmt.Sprintf("%v", *o.Include))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type BlobContainersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BlobContainersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BlobContainersList ...
func (c BlobServicesClient) BlobContainersList(ctx context.Context, id commonids.StorageAccountId, options BlobContainersListOperationOptions) (result BlobContainersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &BlobContainersListCustomPager{},
		Path:          fmt.Sprintf("%s/blobServices/default/containers", id.ID()),
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
		Values *[]ListContainerItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BlobContainersListComplete retrieves all the results into a single object
func (c BlobServicesClient) BlobContainersListComplete(ctx context.Context, id commonids.StorageAccountId, options BlobContainersListOperationOptions) (BlobContainersListCompleteResult, error) {
	return c.BlobContainersListCompleteMatchingPredicate(ctx, id, options, ListContainerItemOperationPredicate{})
}

// BlobContainersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BlobServicesClient) BlobContainersListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, options BlobContainersListOperationOptions, predicate ListContainerItemOperationPredicate) (result BlobContainersListCompleteResult, err error) {
	items := make([]ListContainerItem, 0)

	resp, err := c.BlobContainersList(ctx, id, options)
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

	result = BlobContainersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
