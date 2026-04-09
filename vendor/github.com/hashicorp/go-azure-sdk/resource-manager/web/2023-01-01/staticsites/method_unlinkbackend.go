package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnlinkBackendOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type UnlinkBackendOperationOptions struct {
	IsCleaningAuthConfig *bool
}

func DefaultUnlinkBackendOperationOptions() UnlinkBackendOperationOptions {
	return UnlinkBackendOperationOptions{}
}

func (o UnlinkBackendOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o UnlinkBackendOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o UnlinkBackendOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IsCleaningAuthConfig != nil {
		out.Append("isCleaningAuthConfig", fmt.Sprintf("%v", *o.IsCleaningAuthConfig))
	}
	return &out
}

// UnlinkBackend ...
func (c StaticSitesClient) UnlinkBackend(ctx context.Context, id LinkedBackendId, options UnlinkBackendOperationOptions) (result UnlinkBackendOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: options,
		Path:          id.ID(),
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
