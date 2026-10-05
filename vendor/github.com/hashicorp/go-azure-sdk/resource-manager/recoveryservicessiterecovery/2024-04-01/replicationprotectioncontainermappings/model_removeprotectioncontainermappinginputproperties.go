package replicationprotectioncontainermappings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoveProtectionContainerMappingInputProperties struct {
	ProviderSpecificInput *ReplicationProviderContainerUnmappingInput `json:"providerSpecificInput,omitempty"`
}
