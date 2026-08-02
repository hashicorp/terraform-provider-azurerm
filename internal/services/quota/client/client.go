// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GroupQuotaLimitsClient            *groupquotalimits.GroupQuotaLimitsClient
	GroupQuotasClient                 *groupquotas.GroupQuotasClient
	GroupQuotasSubscriptionsClient    *groupquotassubscriptions.GroupQuotasSubscriptionsClient
	SubscriptionQuotaAllocationClient *subscriptionquotaallocation.SubscriptionQuotaAllocationClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	groupQuotaLimitsClient, err := groupquotalimits.NewGroupQuotaLimitsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GroupQuotaLimits client: %+v", err)
	}
	o.Configure(groupQuotaLimitsClient.Client, o.Authorizers.ResourceManager)

	groupQuotasClient, err := groupquotas.NewGroupQuotasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GroupQuotas client: %+v", err)
	}
	o.Configure(groupQuotasClient.Client, o.Authorizers.ResourceManager)

	groupQuotasSubscriptionsClient, err := groupquotassubscriptions.NewGroupQuotasSubscriptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GroupQuotasSubscriptions client: %+v", err)
	}
	o.Configure(groupQuotasSubscriptionsClient.Client, o.Authorizers.ResourceManager)

	subscriptionQuotaAllocationClient, err := subscriptionquotaallocation.NewSubscriptionQuotaAllocationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SubscriptionQuotaAllocation client: %+v", err)
	}
	o.Configure(subscriptionQuotaAllocationClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		GroupQuotaLimitsClient:            groupQuotaLimitsClient,
		GroupQuotasClient:                 groupQuotasClient,
		GroupQuotasSubscriptionsClient:    groupQuotasSubscriptionsClient,
		SubscriptionQuotaAllocationClient: subscriptionQuotaAllocationClient,
	}, nil
}
