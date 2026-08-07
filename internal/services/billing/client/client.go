// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

const ptuReservationAPIVersion = "2022-11-01"

// Client holds the HTTP client for the Microsoft.Capacity/reservationOrders API.
// This API is not yet in go-azure-sdk so we use the base resourcemanager client directly.
type Client struct {
	ReservationOrdersClient *resourcemanager.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	reservationOrdersClient, err := resourcemanager.NewClient(o.Environment.ResourceManager, "reservationOrders", ptuReservationAPIVersion)
	if err != nil {
		return nil, fmt.Errorf("building PTU Reservation Orders client: %+v", err)
	}
	o.Configure(reservationOrdersClient, o.Authorizers.ResourceManager)

	return &Client{
		ReservationOrdersClient: reservationOrdersClient,
	}, nil
}

// -------------------------------------------------------------------
// Request / Response types for Microsoft.Capacity/reservationOrders
// -------------------------------------------------------------------

type ReservationOrderSku struct {
	Name string `json:"name"`
}

type ReservedResourceProperties struct {
	InstanceFlexibility string `json:"instanceFlexibility"`
}

type ReservationOrderPurchaseProperties struct {
	DisplayName                string                      `json:"displayName"`
	ReservedResourceType       string                      `json:"reservedResourceType"`
	ReservedResourceProperties *ReservedResourceProperties `json:"reservedResourceProperties,omitempty"`
	Term                       string                      `json:"term"`
	BillingPlan                string                      `json:"billingPlan"`
	BillingScopeId             string                      `json:"billingScopeId"`
	AppliedScopeType           string                      `json:"appliedScopeType"`
	Quantity                   int64                       `json:"quantity"`
	Renew                      bool                        `json:"renew"`
	ReservationOrderId         string                      `json:"reservationOrderId,omitempty"`
}

type ReservationOrderPurchaseRequest struct {
	Sku        *ReservationOrderSku                `json:"sku,omitempty"`
	Location   string                              `json:"location"`
	Properties *ReservationOrderPurchaseProperties `json:"properties,omitempty"`
}

type ReservationOrderResponseProperties struct {
	DisplayName       string `json:"displayName"`
	ProvisioningState string `json:"provisioningState"`
	// Azure returns originalQuantity (not quantity) on the reservationOrder GET response.
	OriginalQuantity int64  `json:"originalQuantity"`
	Term             string `json:"term"`
	BillingPlan      string `json:"billingPlan"`
	BillingScopeId   string `json:"billingScopeId"`
	AppliedScopeType string `json:"appliedScopeType"`
	Renew            bool   `json:"renew"`
}

type ReservationOrderResponse struct {
	Id         string                              `json:"id"`
	Name       string                              `json:"name"`
	Location   string                              `json:"location,omitempty"`
	Sku        *ReservationOrderSku                `json:"sku,omitempty"`
	Properties *ReservationOrderResponseProperties `json:"properties,omitempty"`
}

type ReservationReturnProperties struct {
	ReturnReason string `json:"returnReason,omitempty"`
	Message      string `json:"message,omitempty"`
}

type ReservationReturnRequest struct {
	Properties *ReservationReturnProperties `json:"properties,omitempty"`
}

// CalculatePriceResponseProperties is the subset of the calculatePrice response we need.
type CalculatePriceResponseProperties struct {
	ReservationOrderId string `json:"reservationOrderId"`
}

type CalculatePriceResponse struct {
	Properties *CalculatePriceResponseProperties `json:"properties,omitempty"`
}

// -------------------------------------------------------------------
// Client methods
// -------------------------------------------------------------------

// CalculatePrice calls POST /providers/Microsoft.Capacity/calculatePrice to open
// a purchase session and obtain the reservation order ID that must be used in the
// subsequent CreateOrUpdate call. Both calls must use the same order ID and be
// issued in quick succession (Azure expires the session after a few minutes).
func (c *Client) CalculatePrice(ctx context.Context, input ReservationOrderPurchaseRequest) (string, error) {
	opts := sdkclient.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK},
		HttpMethod:          http.MethodPost,
		Path:                "/providers/Microsoft.Capacity/calculatePrice",
	}

	req, err := c.ReservationOrdersClient.NewRequest(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("building calculatePrice request: %+v", err)
	}
	if err = req.Marshal(input); err != nil {
		return "", fmt.Errorf("marshaling calculatePrice request body: %+v", err)
	}

	resp, execErr := req.Execute(ctx)
	if execErr != nil {
		return "", fmt.Errorf("executing calculatePrice request: %+v", execErr)
	}

	var result CalculatePriceResponse
	if err = resp.Unmarshal(&result); err != nil {
		return "", fmt.Errorf("unmarshaling calculatePrice response: %+v", err)
	}

	if result.Properties == nil || result.Properties.ReservationOrderId == "" {
		return "", fmt.Errorf("calculatePrice response did not include a reservationOrderId")
	}

	return result.Properties.ReservationOrderId, nil
}

// CreateOrUpdate purchases (PUT) a new PTU reservation order.
// Returns the final state after creation completes.
func (c *Client) CreateOrUpdate(ctx context.Context, path string, input ReservationOrderPurchaseRequest) (*ReservationOrderResponse, error) {
	opts := sdkclient.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK, http.StatusCreated, http.StatusAccepted},
		HttpMethod:          http.MethodPut,
		Path:                path,
	}

	req, err := c.ReservationOrdersClient.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building create request: %+v", err)
	}
	if err = req.Marshal(input); err != nil {
		return nil, fmt.Errorf("marshaling create request body: %+v", err)
	}

	resp, execErr := req.Execute(ctx)
	if execErr != nil {
		return nil, fmt.Errorf("executing create request: %+v", execErr)
	}

	if resp.Response != nil && resp.Response.StatusCode == http.StatusAccepted {
		poller, pollerErr := resourcemanager.PollerFromResponse(resp, c.ReservationOrdersClient)
		if pollerErr != nil {
			return nil, fmt.Errorf("building create poller: %+v", pollerErr)
		}
		if pollerErr = poller.PollUntilDone(ctx); pollerErr != nil {
			return nil, fmt.Errorf("polling after create: %+v", pollerErr)
		}
	}

	return c.Get(ctx, path)
}

// Get retrieves a PTU reservation order by its ARM path.
// Returns (nil, nil) when the reservation order is not found (404).
func (c *Client) Get(ctx context.Context, path string) (*ReservationOrderResponse, error) {
	opts := sdkclient.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK},
		HttpMethod:          http.MethodGet,
		Path:                path,
	}

	req, err := c.ReservationOrdersClient.NewRequest(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("building get request: %+v", err)
	}

	resp, execErr := req.Execute(ctx)
	if resp != nil && resp.Response != nil && response.WasNotFound(resp.Response) {
		return nil, nil
	}
	if execErr != nil {
		return nil, fmt.Errorf("executing get request: %+v", execErr)
	}

	var result ReservationOrderResponse
	if err = resp.Unmarshal(&result); err != nil {
		return nil, fmt.Errorf("unmarshaling get response: %+v", err)
	}
	return &result, nil
}

// Return attempts to return (cancel/refund) a PTU reservation order.
// Azure only allows returns within the first 30 days of purchase.
func (c *Client) Return(ctx context.Context, path string) error {
	opts := sdkclient.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
		HttpMethod:          http.MethodPost,
		Path:                path + "/return",
	}

	req, err := c.ReservationOrdersClient.NewRequest(ctx, opts)
	if err != nil {
		return fmt.Errorf("building return request: %+v", err)
	}
	if err = req.Marshal(ReservationReturnRequest{
		Properties: &ReservationReturnProperties{
			ReturnReason: "Other",
			Message:      "Destroyed by Terraform",
		},
	}); err != nil {
		return fmt.Errorf("marshaling return request body: %+v", err)
	}

	resp, execErr := req.Execute(ctx)
	if execErr != nil {
		return fmt.Errorf("executing return request: %+v", execErr)
	}

	if resp.Response != nil && resp.Response.StatusCode == http.StatusAccepted {
		poller, pollerErr := resourcemanager.PollerFromResponse(resp, c.ReservationOrdersClient)
		if pollerErr != nil {
			return fmt.Errorf("building return poller: %+v", pollerErr)
		}
		if pollerErr = poller.PollUntilDone(ctx); pollerErr != nil {
			return fmt.Errorf("polling after return: %+v", pollerErr)
		}
	}

	return nil
}
