package resourceguardproxybaseresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DppResourceGuardProxyUnlockDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *UnlockDeleteResponse
}

type DppResourceGuardProxyUnlockDeleteOperationOptions struct {
	XMsAuthorizationAuxiliary *string
}

func DefaultDppResourceGuardProxyUnlockDeleteOperationOptions() DppResourceGuardProxyUnlockDeleteOperationOptions {
	return DppResourceGuardProxyUnlockDeleteOperationOptions{}
}

func (o DppResourceGuardProxyUnlockDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsAuthorizationAuxiliary != nil {
		out.Append("x-ms-authorization-auxiliary", fmt.Sprintf("%v", *o.XMsAuthorizationAuxiliary))
	}
	return &out
}

func (o DppResourceGuardProxyUnlockDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DppResourceGuardProxyUnlockDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// DppResourceGuardProxyUnlockDelete ...
func (c ResourceGuardProxyBaseResourcesClient) DppResourceGuardProxyUnlockDelete(ctx context.Context, id BackupResourceGuardProxyId, input UnlockDeleteRequest, options DppResourceGuardProxyUnlockDeleteOperationOptions) (result DppResourceGuardProxyUnlockDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/unlockDelete", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

	var model UnlockDeleteResponse
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
