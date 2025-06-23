package standbycontainergrouppools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyContainerGroupPoolElasticityProfile struct {
	MaxReadyCapacity int64         `json:"maxReadyCapacity"`
	RefillPolicy     *RefillPolicy `json:"refillPolicy,omitempty"`
}
