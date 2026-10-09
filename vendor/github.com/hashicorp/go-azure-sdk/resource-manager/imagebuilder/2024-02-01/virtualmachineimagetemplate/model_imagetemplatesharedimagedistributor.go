package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateDistributor = ImageTemplateSharedImageDistributor{}

type ImageTemplateSharedImageDistributor struct {
	ExcludeFromLatest  *bool                          `json:"excludeFromLatest,omitempty"`
	GalleryImageId     string                         `json:"galleryImageId"`
	ReplicationRegions *[]string                      `json:"replicationRegions,omitempty"`
	StorageAccountType *SharedImageStorageAccountType `json:"storageAccountType,omitempty"`
	TargetRegions      *[]TargetRegion                `json:"targetRegions,omitempty"`
	Versioning         DistributeVersioner            `json:"versioning"`

	// Fields inherited from ImageTemplateDistributor

	ArtifactTags  *map[string]string `json:"artifactTags,omitempty"`
	RunOutputName string             `json:"runOutputName"`
	Type          string             `json:"type"`
}

func (s ImageTemplateSharedImageDistributor) ImageTemplateDistributor() BaseImageTemplateDistributorImpl {
	return BaseImageTemplateDistributorImpl{
		ArtifactTags:  s.ArtifactTags,
		RunOutputName: s.RunOutputName,
		Type:          s.Type,
	}
}

var _ json.Marshaler = ImageTemplateSharedImageDistributor{}

func (s ImageTemplateSharedImageDistributor) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateSharedImageDistributor
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateSharedImageDistributor: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateSharedImageDistributor: %+v", err)
	}

	decoded["type"] = "SharedImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateSharedImageDistributor: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ImageTemplateSharedImageDistributor{}

func (s *ImageTemplateSharedImageDistributor) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ExcludeFromLatest  *bool                          `json:"excludeFromLatest,omitempty"`
		GalleryImageId     string                         `json:"galleryImageId"`
		ReplicationRegions *[]string                      `json:"replicationRegions,omitempty"`
		StorageAccountType *SharedImageStorageAccountType `json:"storageAccountType,omitempty"`
		TargetRegions      *[]TargetRegion                `json:"targetRegions,omitempty"`
		ArtifactTags       *map[string]string             `json:"artifactTags,omitempty"`
		RunOutputName      string                         `json:"runOutputName"`
		Type               string                         `json:"type"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ExcludeFromLatest = decoded.ExcludeFromLatest
	s.GalleryImageId = decoded.GalleryImageId
	s.ReplicationRegions = decoded.ReplicationRegions
	s.StorageAccountType = decoded.StorageAccountType
	s.TargetRegions = decoded.TargetRegions
	s.ArtifactTags = decoded.ArtifactTags
	s.RunOutputName = decoded.RunOutputName
	s.Type = decoded.Type

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageTemplateSharedImageDistributor into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["versioning"]; ok {
		impl, err := UnmarshalDistributeVersionerImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Versioning' for 'ImageTemplateSharedImageDistributor': %+v", err)
		}
		s.Versioning = impl
	}

	return nil
}
