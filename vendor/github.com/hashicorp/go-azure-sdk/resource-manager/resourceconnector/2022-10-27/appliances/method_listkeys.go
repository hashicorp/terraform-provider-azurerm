package appliances

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ApplianceListKeysResults
}

type ListKeysOperationOptions struct {
	ArtifactType *string
}

func DefaultListKeysOperationOptions() ListKeysOperationOptions {
	return ListKeysOperationOptions{}
}

func (o ListKeysOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListKeysOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListKeysOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.ArtifactType != nil {
		out.Append("artifactType", fmt.Sprintf("%v", *o.ArtifactType))
	}
	return &out
}

// ListKeys ...
func (c AppliancesClient) ListKeys(ctx context.Context, id ApplianceId, options ListKeysOperationOptions) (result ListKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/listKeys", id.ID()),
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

	var model ApplianceListKeysResults
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
