// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0
package tfjson

import (
	"bytes"
	"encoding/json"
)

type LogMessageType string

const (
	MessageTypeVersion    LogMessageType = "version"
	MessageTypeLog        LogMessageType = "log"
	MessageTypeDiagnostic LogMessageType = "diagnostic"
)

// allLogMessageTypes is a slice containing all recognised message types
// to be passed into cmp.AllowUnexported
var allLogMessageTypes = []any{
	VersionLogMessage{},
	LogMessage{},
	DiagnosticLogMessage{},
	UnknownLogMessage{},

	// init
	InitOutputMessage{},
	InitializingStateStoreStartMessage{},

	// query
	ListStartMessage{},
	ListResourceFoundMessage{},
	ListCompleteMessage{},

	// state migrate
	// provider installation
	ProviderInstallationStartMessage{},
	StateStoreProviderInstallationStartMessage{},
	ProviderQueryUsePreviousVersionMessage{},
	ProviderQueryUsePreviousConstraintsMessage{},
	ProviderQueryUseLatestMessage{},
	ProviderVersionAlreadyInstalledMessage{},
	ProviderVersionFoundInCacheDirMessage{},
	ProviderVersionInstallationStartMessage{},
	ProviderVersionInstallationCompleteMessage{},
	BuiltInProviderAvailableMessage{},
	ThirdPartyProvidersInstalledMessage{},
	// dependency lock file
	ProviderLockfileCreatedMessage{},
	ProviderLockfileUpdatedMessage{},
	// Provider trust-related message
	ProviderInteractiveApprovalMessage{},
	ProviderInteractiveRejectionMessage{},
	ProviderAutomaticApprovalMessage{},
	// state migration
	MigrationStartMessage{},
	MigrationCompleteMessage{},
	MigrationErroredMessage{},
	MigrationFinalizedMessage{},
	MigrationSourceInitializationStartMessage{},
	MigrationSourceInitializationCompleteMessage{},
	MigrationDestinationInitializationStartMessage{},
	MigrationDestinationInitializationCompleteMessage{},
}

func unmarshalByType(t LogMessageType, b []byte) (LogMsg, error) {
	d := json.NewDecoder(bytes.NewReader(b))

	// decode numbers as json.Number to avoid losing precision
	d.UseNumber()

	switch t {

	// generic
	case MessageTypeVersion:
		v := VersionLogMessage{}
		return v, d.Decode(&v)
	case MessageTypeLog:
		v := LogMessage{}
		return v, d.Decode(&v)
	case MessageTypeDiagnostic:
		v := DiagnosticLogMessage{}
		return v, d.Decode(&v)

	// init
	case InitOutput:
		v := InitOutputMessage{}
		return v, json.Unmarshal(b, &v)
	case MessageInitializingStateStoreStart:
		v := InitializingStateStoreStartMessage{}
		return v, d.Decode(&v)

	// query
	case MessageListStart:
		v := ListStartMessage{}
		return v, d.Decode(&v)
	case MessageListResourceFound:
		v := ListResourceFoundMessage{}
		return v, d.Decode(&v)
	case MessageListComplete:
		v := ListCompleteMessage{}
		return v, d.Decode(&v)

	// state migrate
	// provider installation
	case MessageProviderInstallationStart:
		v := ProviderInstallationStartMessage{}
		return v, d.Decode(&v)
	case MessageStateStoreProviderInstallationStart:
		v := StateStoreProviderInstallationStartMessage{}
		return v, d.Decode(&v)
	case MessageProviderQueryUsePreviousVersion:
		v := ProviderQueryUsePreviousVersionMessage{}
		return v, d.Decode(&v)
	case MessageProviderQueryUsePreviousConstraints:
		v := ProviderQueryUsePreviousConstraintsMessage{}
		return v, d.Decode(&v)
	case MessageProviderQueryUseLatest:
		v := ProviderQueryUseLatestMessage{}
		return v, d.Decode(&v)
	case MessageProviderVersionAlreadyInstalled:
		v := ProviderVersionAlreadyInstalledMessage{}
		return v, d.Decode(&v)
	case MessageProviderVersionFoundInCacheDir:
		v := ProviderVersionFoundInCacheDirMessage{}
		return v, d.Decode(&v)
	case MessageProviderVersionInstallationStart:
		v := ProviderVersionInstallationStartMessage{}
		return v, d.Decode(&v)
	case MessageProviderVersionInstallationComplete:
		v := ProviderVersionInstallationCompleteMessage{}
		return v, d.Decode(&v)
	case MessageBuiltInProviderAvailable:
		v := BuiltInProviderAvailableMessage{}
		return v, d.Decode(&v)
	case MessageThirdPartyProvidersInstalled:
		v := ThirdPartyProvidersInstalledMessage{}
		return v, d.Decode(&v)
	// dependency lock file
	case MessageProviderLockfileCreated:
		v := ProviderLockfileCreatedMessage{}
		return v, d.Decode(&v)
	case MessageProviderLockfileUpdated:
		v := ProviderLockfileUpdatedMessage{}
		return v, d.Decode(&v)
	// provider trust
	case MessageProviderInteractiveApproval:
		v := ProviderInteractiveApprovalMessage{}
		return v, d.Decode(&v)
	case MessageProviderInteractiveRejection:
		v := ProviderInteractiveRejectionMessage{}
		return v, d.Decode(&v)
	case MessageProviderAutomaticApproval:
		v := ProviderAutomaticApprovalMessage{}
		return v, d.Decode(&v)
	// state migration
	case MessageMigrationStart:
		v := MigrationStartMessage{}
		return v, d.Decode(&v)
	case MessageMigrationComplete:
		v := MigrationCompleteMessage{}
		return v, d.Decode(&v)
	case MessageMigrationErrored:
		v := MigrationErroredMessage{}
		return v, d.Decode(&v)
	case MessageMigrationFinalized:
		v := MigrationFinalizedMessage{}
		return v, d.Decode(&v)
	case MessageMigrationSourceInitializationStart:
		v := MigrationSourceInitializationStartMessage{}
		return v, d.Decode(&v)
	case MessageMigrationSourceInitializationComplete:
		v := MigrationSourceInitializationCompleteMessage{}
		return v, d.Decode(&v)
	case MessageMigrationDestinationInitializationStart:
		v := MigrationDestinationInitializationStartMessage{}
		return v, d.Decode(&v)
	case MessageMigrationDestinationInitializationComplete:
		v := MigrationDestinationInitializationCompleteMessage{}
		return v, d.Decode(&v)
	}

	v := UnknownLogMessage{}
	return v, d.Decode(&v)
}
