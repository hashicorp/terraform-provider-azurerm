package scalingplan

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScalingPlanClient struct {
	Client  autorest.Client
	baseUri string
}

func NewScalingPlanClientWithBaseURI(endpoint string) ScalingPlanClient {
	return ScalingPlanClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
