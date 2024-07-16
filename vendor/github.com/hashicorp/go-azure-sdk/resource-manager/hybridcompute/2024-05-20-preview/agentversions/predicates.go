package agentversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentVersionOperationPredicate struct {
	AgentVersion *string
	DownloadLink *string
	OsType       *string
}

func (p AgentVersionOperationPredicate) Matches(input AgentVersion) bool {

	if p.AgentVersion != nil && (input.AgentVersion == nil || *p.AgentVersion != *input.AgentVersion) {
		return false
	}

	if p.DownloadLink != nil && (input.DownloadLink == nil || *p.DownloadLink != *input.DownloadLink) {
		return false
	}

	if p.OsType != nil && (input.OsType == nil || *p.OsType != *input.OsType) {
		return false
	}

	return true
}
