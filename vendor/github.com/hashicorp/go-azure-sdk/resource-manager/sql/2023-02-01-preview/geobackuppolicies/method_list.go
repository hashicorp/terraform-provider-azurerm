package geobackuppolicies

import (
	"context"
	"fmt"
	"net/http"

<<<<<<< HEAD
=======
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GeoBackupPolicy
}

type ListCompleteResult struct {
	Items []GeoBackupPolicy
}

// List ...
<<<<<<< HEAD
func (c GeoBackupPoliciesClient) List(ctx context.Context, id DatabaseId) (result ListOperationResponse, err error) {
=======
func (c GeoBackupPoliciesClient) List(ctx context.Context, id commonids.SqlDatabaseId) (result ListOperationResponse, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/geoBackupPolicies", id.ID()),
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
		Values *[]GeoBackupPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListComplete retrieves all the results into a single object
<<<<<<< HEAD
func (c GeoBackupPoliciesClient) ListComplete(ctx context.Context, id DatabaseId) (ListCompleteResult, error) {
=======
func (c GeoBackupPoliciesClient) ListComplete(ctx context.Context, id commonids.SqlDatabaseId) (ListCompleteResult, error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	return c.ListCompleteMatchingPredicate(ctx, id, GeoBackupPolicyOperationPredicate{})
}

// ListCompleteMatchingPredicate retrieves all the results and then applies the predicate
<<<<<<< HEAD
func (c GeoBackupPoliciesClient) ListCompleteMatchingPredicate(ctx context.Context, id DatabaseId, predicate GeoBackupPolicyOperationPredicate) (result ListCompleteResult, err error) {
=======
func (c GeoBackupPoliciesClient) ListCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, predicate GeoBackupPolicyOperationPredicate) (result ListCompleteResult, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	items := make([]GeoBackupPolicy, 0)

	resp, err := c.List(ctx, id)
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

	result = ListCompleteResult{
		Items: items,
	}
	return
}
