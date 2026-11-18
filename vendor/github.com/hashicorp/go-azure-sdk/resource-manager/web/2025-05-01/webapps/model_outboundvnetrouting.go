package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundVnetRouting struct {
	AllTraffic           *bool `json:"allTraffic,omitempty"`
	ApplicationTraffic   *bool `json:"applicationTraffic,omitempty"`
	BackupRestoreTraffic *bool `json:"backupRestoreTraffic,omitempty"`
	ContentShareTraffic  *bool `json:"contentShareTraffic,omitempty"`
	ImagePullTraffic     *bool `json:"imagePullTraffic,omitempty"`
}
