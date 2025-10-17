package topicrecords

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicRecordsClient struct {
	Client *resourcemanager.Client
}

func NewTopicRecordsClientWithBaseURI(sdkApi sdkEnv.Api) (*TopicRecordsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "topicrecords", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TopicRecordsClient: %+v", err)
	}

	return &TopicRecordsClient{
		Client: client,
	}, nil
}
