package servicetags

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagInformationListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceTagInformation
}

type ServiceTagInformationListCompleteResult struct {
	Items []ServiceTagInformation
}

type ServiceTagInformationListOperationOptions struct {
	NoAddressPrefixes *bool
	TagName           *string
}

func DefaultServiceTagInformationListOperationOptions() ServiceTagInformationListOperationOptions {
	return ServiceTagInformationListOperationOptions{}
}

func (o ServiceTagInformationListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ServiceTagInformationListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ServiceTagInformationListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.NoAddressPrefixes != nil {
		out.Append("noAddressPrefixes", fmt.Sprintf("%v", *o.NoAddressPrefixes))
	}
	if o.TagName != nil {
		out.Append("tagName", fmt.Sprintf("%v", *o.TagName))
	}
	return &out
}

// ServiceTagInformationList ...
func (c ServiceTagsClient) ServiceTagInformationList(ctx context.Context, id LocationId, options ServiceTagInformationListOperationOptions) (result ServiceTagInformationListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/serviceTagDetails", id.ID()),
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
		Values *[]ServiceTagInformation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServiceTagInformationListComplete retrieves all the results into a single object
func (c ServiceTagsClient) ServiceTagInformationListComplete(ctx context.Context, id LocationId, options ServiceTagInformationListOperationOptions) (ServiceTagInformationListCompleteResult, error) {
	return c.ServiceTagInformationListCompleteMatchingPredicate(ctx, id, options, ServiceTagInformationOperationPredicate{})
}

// ServiceTagInformationListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceTagsClient) ServiceTagInformationListCompleteMatchingPredicate(ctx context.Context, id LocationId, options ServiceTagInformationListOperationOptions, predicate ServiceTagInformationOperationPredicate) (result ServiceTagInformationListCompleteResult, err error) {
	items := make([]ServiceTagInformation, 0)

	resp, err := c.ServiceTagInformationList(ctx, id, options)
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

	result = ServiceTagInformationListCompleteResult{
		Items: items,
	}
	return
}
