package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Context struct {
	ContextType        *string `json:"contextType,omitempty"`
	NotificationSource *string `json:"notificationSource,omitempty"`
}
