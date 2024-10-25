package firewallpolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleCollectionGroupDraftsDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

// FirewallPolicyRuleCollectionGroupDraftsDelete ...
func (c FirewallPoliciesClient) FirewallPolicyRuleCollectionGroupDraftsDelete(ctx context.Context, id RuleCollectionGroupId) (result FirewallPolicyRuleCollectionGroupDraftsDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod: http.MethodDelete,
		Path:       fmt.Sprintf("%s/ruleCollectionGroupDrafts/default", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	return
}
