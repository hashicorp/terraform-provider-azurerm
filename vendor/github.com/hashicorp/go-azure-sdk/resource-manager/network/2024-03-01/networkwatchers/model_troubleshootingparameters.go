package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TroubleshootingParameters struct {
	Properties       TroubleshootingProperties `json:"properties"`
	TargetResourceId string                    `json:"targetResourceId"`
}
