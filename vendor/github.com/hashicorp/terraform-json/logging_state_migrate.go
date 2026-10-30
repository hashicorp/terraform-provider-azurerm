// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package tfjson

const (
	// Provider installation messages
	MessageProviderInstallationStart           LogMessageType = "provider_installation_start"
	MessageStateStoreProviderInstallationStart LogMessageType = "state_store_provider_installation_start"
	MessageProviderQueryUsePreviousVersion     LogMessageType = "provider_query_use_previous_version"
	MessageProviderQueryUsePreviousConstraints LogMessageType = "provider_query_use_constraints"
	MessageProviderQueryUseLatest              LogMessageType = "provider_query_use_latest"
	MessageProviderVersionAlreadyInstalled     LogMessageType = "provider_version_already_installed"
	MessageProviderVersionFoundInCacheDir      LogMessageType = "provider_version_found_in_cache_dir"
	MessageProviderVersionInstallationStart    LogMessageType = "provider_version_installation_start"
	MessageProviderVersionInstallationComplete LogMessageType = "provider_version_installation_complete"
	MessageBuiltInProviderAvailable            LogMessageType = "built_in_provider_available"
	MessageThirdPartyProvidersInstalled        LogMessageType = "third_party_providers_installed"

	// Dependency lock file messages
	MessageProviderLockfileCreated LogMessageType = "provider_lockfile_created"
	MessageProviderLockfileUpdated LogMessageType = "provider_lockfile_updated"

	// Provider trust-related messages
	MessageProviderInteractiveApproval  LogMessageType = "provider_interactive_approval"
	MessageProviderInteractiveRejection LogMessageType = "provider_interactive_rejection"
	MessageProviderAutomaticApproval    LogMessageType = "provider_automatic_approval"

	// State migration-related messages
	MessageMigrationStart                             LogMessageType = "migration_start"
	MessageMigrationComplete                          LogMessageType = "migration_complete"
	MessageMigrationErrored                           LogMessageType = "migration_errored"
	MessageMigrationFinalized                         LogMessageType = "migration_finalized"
	MessageMigrationSourceInitializationStart         LogMessageType = "migration_source_initialization_start"
	MessageMigrationSourceInitializationComplete      LogMessageType = "migration_source_initialization_complete"
	MessageMigrationDestinationInitializationStart    LogMessageType = "migration_destination_initialization_start"
	MessageMigrationDestinationInitializationComplete LogMessageType = "migration_destination_initialization_complete"
)

type ProviderInstallationStartMessage struct {
	baseLogMessage
}

type StateStoreProviderInstallationStartMessage struct {
	baseLogMessage
}

type ProviderQueryUsePreviousVersionMessage struct {
	baseLogMessage
}

type ProviderQueryUsePreviousConstraintsMessage struct {
	baseLogMessage
}

type ProviderQueryUseLatestMessage struct {
	baseLogMessage
}

type ProviderVersionAlreadyInstalledMessage struct {
	baseLogMessage
}

type ProviderVersionFoundInCacheDirMessage struct {
	baseLogMessage
}

type ProviderVersionInstallationStartMessage struct {
	baseLogMessage
}

type ProviderVersionInstallationCompleteMessage struct {
	baseLogMessage
}

type BuiltInProviderAvailableMessage struct {
	baseLogMessage
}

type ThirdPartyProvidersInstalledMessage struct {
	baseLogMessage
}

type ProviderLockfileCreatedMessage struct {
	baseLogMessage
}

type ProviderLockfileUpdatedMessage struct {
	baseLogMessage
}

type ProviderInteractiveApprovalMessage struct {
	baseLogMessage
}

type ProviderInteractiveRejectionMessage struct {
	baseLogMessage
}

type ProviderAutomaticApprovalMessage struct {
	baseLogMessage
}

type MigrationStartMessage struct {
	baseLogMessage
}

type MigrationCompleteMessage struct {
	baseLogMessage
}

type MigrationErroredMessage struct {
	baseLogMessage
}

type MigrationFinalizedMessage struct {
	baseLogMessage
}

type MigrationSourceInitializationStartMessage struct {
	baseLogMessage
}

type MigrationSourceInitializationCompleteMessage struct {
	baseLogMessage
}

type MigrationDestinationInitializationStartMessage struct {
	baseLogMessage
}

type MigrationDestinationInitializationCompleteMessage struct {
	baseLogMessage
}
