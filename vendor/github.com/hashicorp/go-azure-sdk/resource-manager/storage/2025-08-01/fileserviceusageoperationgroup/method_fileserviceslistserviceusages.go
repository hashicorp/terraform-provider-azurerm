package fileserviceusageoperationgroup

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

type FileServicesListServiceUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FileServiceUsage
}

type FileServicesListServiceUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FileServiceUsage
}

type FileServicesListServiceUsagesOperationOptions struct {
	Maxpagesize *int64
}

func DefaultFileServicesListServiceUsagesOperationOptions() FileServicesListServiceUsagesOperationOptions {
	return FileServicesListServiceUsagesOperationOptions{}
}

func (o FileServicesListServiceUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o FileServicesListServiceUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o FileServicesListServiceUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type FileServicesListServiceUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FileServicesListServiceUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FileServicesListServiceUsages ...
func (c FileServiceUsageOperationGroupClient) FileServicesListServiceUsages(ctx context.Context, id commonids.StorageAccountId, options FileServicesListServiceUsagesOperationOptions) (result FileServicesListServiceUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &FileServicesListServiceUsagesCustomPager{},
		Path:          fmt.Sprintf("%s/fileServices/default/usages", id.ID()),
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
		Values *[]FileServiceUsage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FileServicesListServiceUsagesComplete retrieves all the results into a single object
func (c FileServiceUsageOperationGroupClient) FileServicesListServiceUsagesComplete(ctx context.Context, id commonids.StorageAccountId, options FileServicesListServiceUsagesOperationOptions) (FileServicesListServiceUsagesCompleteResult, error) {
	return c.FileServicesListServiceUsagesCompleteMatchingPredicate(ctx, id, options, FileServiceUsageOperationPredicate{})
}

// FileServicesListServiceUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FileServiceUsageOperationGroupClient) FileServicesListServiceUsagesCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, options FileServicesListServiceUsagesOperationOptions, predicate FileServiceUsageOperationPredicate) (result FileServicesListServiceUsagesCompleteResult, err error) {
	items := make([]FileServiceUsage, 0)

	resp, err := c.FileServicesListServiceUsages(ctx, id, options)
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

	result = FileServicesListServiceUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
