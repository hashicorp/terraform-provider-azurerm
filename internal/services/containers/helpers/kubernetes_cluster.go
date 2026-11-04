// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
)

func Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func Base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}

func RetryNodePoolCreation(ctx context.Context, client *agentpools.AgentPoolsClient, id agentpools.AgentPoolId, profile agentpools.AgentPool) error {
	// retries the creation of a node pool 3 times
	var err error
	for range 3 {
		if err = client.CreateOrUpdateThenPoll(ctx, id, profile, agentpools.DefaultCreateOrUpdateOperationOptions()); err == nil {
			return nil
		}
	}

	return err
}
