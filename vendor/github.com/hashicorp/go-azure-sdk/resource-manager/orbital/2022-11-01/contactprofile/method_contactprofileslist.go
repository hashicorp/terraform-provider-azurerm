package contactprofile

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

type ContactProfilesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContactProfile
}

type ContactProfilesListCompleteResult struct {
	Items []ContactProfile
}

// ContactProfilesList ...
func (c ContactProfileClient) ContactProfilesList(ctx context.Context, id commonids.ResourceGroupId) (result ContactProfilesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Orbital/contactProfiles", id.ID()),
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
		Values *[]ContactProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ContactProfilesListComplete retrieves all the results into a single object
func (c ContactProfileClient) ContactProfilesListComplete(ctx context.Context, id commonids.ResourceGroupId) (ContactProfilesListCompleteResult, error) {
	return c.ContactProfilesListCompleteMatchingPredicate(ctx, id, ContactProfileOperationPredicate{})
}

// ContactProfilesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContactProfileClient) ContactProfilesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ContactProfileOperationPredicate) (result ContactProfilesListCompleteResult, err error) {
	items := make([]ContactProfile, 0)

	resp, err := c.ContactProfilesList(ctx, id)
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

	result = ContactProfilesListCompleteResult{
		Items: items,
	}
	return
}
