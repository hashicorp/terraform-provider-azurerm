package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeProfile struct {
	Subscription *EdgeProfileSubscription `json:"subscription,omitempty"`
}
