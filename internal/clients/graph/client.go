// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graph

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/msgraph"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type options struct {
	query odata.Query
}

func (o options) ToHeaders() *client.Headers {
	h := client.Headers{}
	h.AppendHeader(o.query.Headers())
	return &h
}

func (o options) ToOData() *odata.Query {
	return &o.query
}

func (o options) ToQuery() *client.QueryParams {
	q := client.QueryParams{}
	q.AppendValues(o.query.Values())
	return &q
}

type directoryObjectModel struct {
	ID *string `json:"id"`
}

func graphClient(authorizer auth.Authorizer, environment environments.Environment) (*msgraph.Client, error) {
	client, err := msgraph.NewMsGraphClient(environment.MicrosoftGraph, "Graph", msgraph.VersionOnePointZero)
	if err != nil {
		return nil, fmt.Errorf("building client: %+v", err)
	}

	client.Authorizer = authorizer

	return client, nil
}

func ServicePrincipalObjectID(ctx context.Context, authorizer auth.Authorizer, environment environments.Environment, clientId string) (*string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(5*time.Minute))
		defer cancel()
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: options{
			query: odata.Query{
				Filter: fmt.Sprintf("appId eq '%s'", clientId),
			},
		},
		Path: "/servicePrincipals",
	}

	client, err := graphClient(authorizer, environment)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building new request: %+v", err)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("executing request: %+v", err)
	}

	model := struct {
		ServicePrincipals []directoryObjectModel `json:"value"`
	}{}
	if err := resp.Unmarshal(&model); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %+v", err)
	}

	if len(model.ServicePrincipals) != 1 {
		return nil, fmt.Errorf("unexpected number of results, expected 1, received %d", len(model.ServicePrincipals))
	}

	id := model.ServicePrincipals[0].ID
	if id == nil {
		return nil, fmt.Errorf("returned object ID was nil")
	}

	return id, nil
}

func UserPrincipalObjectID(ctx context.Context, authorizer auth.Authorizer, environment environments.Environment) (*string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(5*time.Minute))
		defer cancel()
	}

	opts := client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: nil,
		Path:          "/me",
	}

	client, err := graphClient(authorizer, environment)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building new request: %+v", err)
	}

	resp, err := req.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("executing request: %+v", err)
	}

	model := directoryObjectModel{}
	if err := resp.Unmarshal(&model); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %+v", err)
	}

	if model.ID == nil {
		return nil, fmt.Errorf("returned object ID was nil")
	}

	return model.ID, nil
}
