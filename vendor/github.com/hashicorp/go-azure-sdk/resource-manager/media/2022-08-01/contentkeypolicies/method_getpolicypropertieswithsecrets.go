package contentkeypolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetPolicyPropertiesWithSecretsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ContentKeyPolicyProperties
}

// GetPolicyPropertiesWithSecrets ...
func (c ContentKeyPoliciesClient) GetPolicyPropertiesWithSecrets(ctx context.Context, id ContentKeyPolicyId) (result GetPolicyPropertiesWithSecretsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/getPolicyPropertiesWithSecrets", id.ID()),
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

	var model ContentKeyPolicyProperties
	result.Model = &model

	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
