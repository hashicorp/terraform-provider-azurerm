// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

// Subscriptions is a list of AAD Subscriptions  which can be used for test purposes
type Subscriptions struct {
	// Primary is the Primary/Default AAD Subscription. This is specified via "ARM_SUBSCRIPTION_ID"
	Primary string

	// Secondary is the Secondary AAD Subscrption which should be used for testing. This is specified via "ARM_TEST_SUBSCRIPTION_ID_ALT"
	Secondary string
}
