// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ResourceProviderId{}

type ResourceProviderId struct {
	SubscriptionId   string
	ResourceProvider string
}

func (id ResourceProviderId) ID() string {
	fmtString := "/subscriptions/%s/providers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceProvider)
}

func (id ResourceProviderId) String() string {
	return fmt.Sprintf("Resource Provider %q", id.ResourceProvider)
}
