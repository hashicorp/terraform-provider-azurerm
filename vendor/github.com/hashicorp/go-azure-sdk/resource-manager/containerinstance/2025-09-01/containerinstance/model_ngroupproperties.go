package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NGroupProperties struct {
	ContainerGroupProfiles *[]ContainerGroupProfileStub `json:"containerGroupProfiles,omitempty"`
	ElasticProfile         *ElasticProfile              `json:"elasticProfile,omitempty"`
	PlacementProfile       *PlacementProfile            `json:"placementProfile,omitempty"`
	ProvisioningState      *NGroupProvisioningState     `json:"provisioningState,omitempty"`
	UpdateProfile          *UpdateProfile               `json:"updateProfile,omitempty"`
}
