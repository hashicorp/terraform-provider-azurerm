package loadtests

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadTestResourcePatchRequestBody struct {
	Identity   *identity.SystemAssigned                    `json:"identity,omitempty"`
	Properties *LoadTestResourcePatchRequestBodyProperties `json:"properties,omitempty"`
	Tags       *interface{}                                `json:"tags,omitempty"`
}
