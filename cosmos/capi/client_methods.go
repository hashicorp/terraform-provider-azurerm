package capi

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"net/http"
)

func (c BaseClient) Create(ctx context.Context, path string, preparers []autorest.PrepareDecorator, responders []autorest.RespondDecorator) (response *autorest.Response, err error) {
	req, err := c.Preparer(ctx, "Create", http.MethodPost, path, preparers...)
	if err != nil {
		return nil, err
	}

	resp, err := c.Sender("Create", req)
	if err == nil {
		resp, err = c.Responder("Create", resp.Response, http.StatusCreated, responders...)
	}

	return &resp, err
}

func (c BaseClient) Get(ctx context.Context, path string, responders ...autorest.RespondDecorator) (response *autorest.Response, err error) {

	req, err := c.Preparer(ctx, "Get", http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := c.Sender("Get", req)
	if err == nil {
		resp, err = c.Responder("Get", resp.Response, http.StatusOK, responders...)
	}

	return &resp, err
}

func (c BaseClient) Delete(ctx context.Context, path string) (response *autorest.Response, err error) {

	req, err := c.Preparer(ctx, "Delete", http.MethodDelete, path)
	if err != nil {
		return nil, err
	}

	resp, err := c.Sender("Delete", req)
	if err == nil {
		resp, err = c.Responder("Delete", resp.Response, http.StatusNoContent)
	}

	return &resp, err
}
