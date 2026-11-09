package bigdatapools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationType string

const (
	ConfigurationTypeArtifact ConfigurationType = "Artifact"
	ConfigurationTypeFile     ConfigurationType = "File"
)

func PossibleValuesForConfigurationType() []string {
	return []string{
		string(ConfigurationTypeArtifact),
		string(ConfigurationTypeFile),
	}
}

func (s *ConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationType(input string) (*ConfigurationType, error) {
	vals := map[string]ConfigurationType{
		"artifact": ConfigurationTypeArtifact,
		"file":     ConfigurationTypeFile,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationType(input)
	return &out, nil
}

type NodeSize string

const (
	NodeSizeLarge    NodeSize = "Large"
	NodeSizeMedium   NodeSize = "Medium"
	NodeSizeNone     NodeSize = "None"
	NodeSizeSmall    NodeSize = "Small"
	NodeSizeXLarge   NodeSize = "XLarge"
	NodeSizeXXLarge  NodeSize = "XXLarge"
	NodeSizeXXXLarge NodeSize = "XXXLarge"
)

func PossibleValuesForNodeSize() []string {
	return []string{
		string(NodeSizeLarge),
		string(NodeSizeMedium),
		string(NodeSizeNone),
		string(NodeSizeSmall),
		string(NodeSizeXLarge),
		string(NodeSizeXXLarge),
		string(NodeSizeXXXLarge),
	}
}

func (s *NodeSize) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeSize(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeSize(input string) (*NodeSize, error) {
	vals := map[string]NodeSize{
		"large":    NodeSizeLarge,
		"medium":   NodeSizeMedium,
		"none":     NodeSizeNone,
		"small":    NodeSizeSmall,
		"xlarge":   NodeSizeXLarge,
		"xxlarge":  NodeSizeXXLarge,
		"xxxlarge": NodeSizeXXXLarge,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeSize(input)
	return &out, nil
}

type NodeSizeFamily string

const (
	NodeSizeFamilyHardwareAcceleratedFPGA NodeSizeFamily = "HardwareAcceleratedFPGA"
	NodeSizeFamilyHardwareAcceleratedGPU  NodeSizeFamily = "HardwareAcceleratedGPU"
	NodeSizeFamilyMemoryOptimized         NodeSizeFamily = "MemoryOptimized"
	NodeSizeFamilyNone                    NodeSizeFamily = "None"
)

func PossibleValuesForNodeSizeFamily() []string {
	return []string{
		string(NodeSizeFamilyHardwareAcceleratedFPGA),
		string(NodeSizeFamilyHardwareAcceleratedGPU),
		string(NodeSizeFamilyMemoryOptimized),
		string(NodeSizeFamilyNone),
	}
}

func (s *NodeSizeFamily) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeSizeFamily(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeSizeFamily(input string) (*NodeSizeFamily, error) {
	vals := map[string]NodeSizeFamily{
		"hardwareacceleratedfpga": NodeSizeFamilyHardwareAcceleratedFPGA,
		"hardwareacceleratedgpu":  NodeSizeFamilyHardwareAcceleratedGPU,
		"memoryoptimized":         NodeSizeFamilyMemoryOptimized,
		"none":                    NodeSizeFamilyNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeSizeFamily(input)
	return &out, nil
}
