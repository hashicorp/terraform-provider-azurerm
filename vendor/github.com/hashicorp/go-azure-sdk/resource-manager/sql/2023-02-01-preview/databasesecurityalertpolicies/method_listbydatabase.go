package databasesecurityalertpolicies

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

type ListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatabaseSecurityAlertPolicy
}

type ListByDatabaseCompleteResult struct {
	Items []DatabaseSecurityAlertPolicy
}

// ListByDatabase ...
<<<<<<< HEAD
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabase(ctx context.Context, id DatabaseId) (result ListByDatabaseOperationResponse, err error) {
=======
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabase(ctx context.Context, id commonids.SqlDatabaseId) (result ListByDatabaseOperationResponse, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/securityAlertPolicies", id.ID()),
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
		Values *[]DatabaseSecurityAlertPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDatabaseComplete retrieves all the results into a single object
<<<<<<< HEAD
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabaseComplete(ctx context.Context, id DatabaseId) (ListByDatabaseCompleteResult, error) {
=======
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabaseComplete(ctx context.Context, id commonids.SqlDatabaseId) (ListByDatabaseCompleteResult, error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	return c.ListByDatabaseCompleteMatchingPredicate(ctx, id, DatabaseSecurityAlertPolicyOperationPredicate{})
}

// ListByDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
<<<<<<< HEAD
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabaseCompleteMatchingPredicate(ctx context.Context, id DatabaseId, predicate DatabaseSecurityAlertPolicyOperationPredicate) (result ListByDatabaseCompleteResult, err error) {
=======
func (c DatabaseSecurityAlertPoliciesClient) ListByDatabaseCompleteMatchingPredicate(ctx context.Context, id commonids.SqlDatabaseId, predicate DatabaseSecurityAlertPolicyOperationPredicate) (result ListByDatabaseCompleteResult, err error) {
>>>>>>> 5e957238fca9519400c2479c7d1f73e3d1b0871c
	items := make([]DatabaseSecurityAlertPolicy, 0)

	resp, err := c.ListByDatabase(ctx, id)
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

	result = ListByDatabaseCompleteResult{
		Items: items,
	}
	return
}
