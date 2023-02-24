package eventsources

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceKind string

const (
	EventSourceKindMicrosoftPointEventHub EventSourceKind = "Microsoft.EventHub"
	EventSourceKindMicrosoftPointIoTHub   EventSourceKind = "Microsoft.IoTHub"
)

func PossibleValuesForEventSourceKind() []string {
	return []string{
		string(EventSourceKindMicrosoftPointEventHub),
		string(EventSourceKindMicrosoftPointIoTHub),
	}
}

func parseEventSourceKind(input string) (*EventSourceKind, error) {
	vals := map[string]EventSourceKind{
		"microsoft.eventhub": EventSourceKindMicrosoftPointEventHub,
		"microsoft.iothub":   EventSourceKindMicrosoftPointIoTHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EventSourceKind(input)
	return &out, nil
}

type IngressStartAtType string

const (
	IngressStartAtTypeCustomEnqueuedTime      IngressStartAtType = "CustomEnqueuedTime"
	IngressStartAtTypeEarliestAvailable       IngressStartAtType = "EarliestAvailable"
	IngressStartAtTypeEventSourceCreationTime IngressStartAtType = "EventSourceCreationTime"
)

func PossibleValuesForIngressStartAtType() []string {
	return []string{
		string(IngressStartAtTypeCustomEnqueuedTime),
		string(IngressStartAtTypeEarliestAvailable),
		string(IngressStartAtTypeEventSourceCreationTime),
	}
}

func parseIngressStartAtType(input string) (*IngressStartAtType, error) {
	vals := map[string]IngressStartAtType{
		"customenqueuedtime":      IngressStartAtTypeCustomEnqueuedTime,
		"earliestavailable":       IngressStartAtTypeEarliestAvailable,
		"eventsourcecreationtime": IngressStartAtTypeEventSourceCreationTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressStartAtType(input)
	return &out, nil
}

type Kind string

const (
	KindMicrosoftPointEventHub Kind = "Microsoft.EventHub"
	KindMicrosoftPointIoTHub   Kind = "Microsoft.IoTHub"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindMicrosoftPointEventHub),
		string(KindMicrosoftPointIoTHub),
	}
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"microsoft.eventhub": KindMicrosoftPointEventHub,
		"microsoft.iothub":   KindMicrosoftPointIoTHub,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type LocalTimestampFormat string

const (
	LocalTimestampFormatEmbedded LocalTimestampFormat = "Embedded"
)

func PossibleValuesForLocalTimestampFormat() []string {
	return []string{
		string(LocalTimestampFormatEmbedded),
	}
}

func parseLocalTimestampFormat(input string) (*LocalTimestampFormat, error) {
	vals := map[string]LocalTimestampFormat{
		"embedded": LocalTimestampFormatEmbedded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalTimestampFormat(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
