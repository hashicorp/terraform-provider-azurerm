package vaults

type AccessPolicyUpdateKind string

const (
	AccessPolicyUpdateKindAdd     AccessPolicyUpdateKind = "add"
	AccessPolicyUpdateKindRemove  AccessPolicyUpdateKind = "remove"
	AccessPolicyUpdateKindReplace AccessPolicyUpdateKind = "replace"
)

type ActionsRequired string

const (
	ActionsRequiredNone ActionsRequired = "None"
)

type CertificatePermissions string

const (
	CertificatePermissionsAll            CertificatePermissions = "all"
	CertificatePermissionsBackup         CertificatePermissions = "backup"
	CertificatePermissionsCreate         CertificatePermissions = "create"
	CertificatePermissionsDelete         CertificatePermissions = "delete"
	CertificatePermissionsDeleteissuers  CertificatePermissions = "deleteissuers"
	CertificatePermissionsGet            CertificatePermissions = "get"
	CertificatePermissionsGetissuers     CertificatePermissions = "getissuers"
	CertificatePermissionsImport         CertificatePermissions = "import"
	CertificatePermissionsList           CertificatePermissions = "list"
	CertificatePermissionsListissuers    CertificatePermissions = "listissuers"
	CertificatePermissionsManagecontacts CertificatePermissions = "managecontacts"
	CertificatePermissionsManageissuers  CertificatePermissions = "manageissuers"
	CertificatePermissionsPurge          CertificatePermissions = "purge"
	CertificatePermissionsRecover        CertificatePermissions = "recover"
	CertificatePermissionsRestore        CertificatePermissions = "restore"
	CertificatePermissionsSetissuers     CertificatePermissions = "setissuers"
	CertificatePermissionsUpdate         CertificatePermissions = "update"
)

type CreateMode string

const (
	CreateModeDefault CreateMode = "default"
	CreateModeRecover CreateMode = "recover"
)

type Filter string

const (
	FilterResourceTypeEqMicrosoftPointKeyVaultVaults Filter = "resourceType eq 'Microsoft.KeyVault/vaults'"
)

type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "Application"
	IdentityTypeKey             IdentityType = "Key"
	IdentityTypeManagedIdentity IdentityType = "ManagedIdentity"
	IdentityTypeUser            IdentityType = "User"
)

type KeyPermissions string

const (
	KeyPermissionsAll       KeyPermissions = "all"
	KeyPermissionsBackup    KeyPermissions = "backup"
	KeyPermissionsCreate    KeyPermissions = "create"
	KeyPermissionsDecrypt   KeyPermissions = "decrypt"
	KeyPermissionsDelete    KeyPermissions = "delete"
	KeyPermissionsEncrypt   KeyPermissions = "encrypt"
	KeyPermissionsGet       KeyPermissions = "get"
	KeyPermissionsImport    KeyPermissions = "import"
	KeyPermissionsList      KeyPermissions = "list"
	KeyPermissionsPurge     KeyPermissions = "purge"
	KeyPermissionsRecover   KeyPermissions = "recover"
	KeyPermissionsRelease   KeyPermissions = "release"
	KeyPermissionsRestore   KeyPermissions = "restore"
	KeyPermissionsRotate    KeyPermissions = "rotate"
	KeyPermissionsSign      KeyPermissions = "sign"
	KeyPermissionsUnwrapKey KeyPermissions = "unwrapKey"
	KeyPermissionsUpdate    KeyPermissions = "update"
	KeyPermissionsVerify    KeyPermissions = "verify"
	KeyPermissionsWrapKey   KeyPermissions = "wrapKey"
)

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "Allow"
	NetworkRuleActionDeny  NetworkRuleAction = "Deny"
)

type NetworkRuleBypassOptions string

const (
	NetworkRuleBypassOptionsAzureServices NetworkRuleBypassOptions = "AzureServices"
	NetworkRuleBypassOptionsNone          NetworkRuleBypassOptions = "None"
)

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating     PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting     PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateDisconnected PrivateEndpointConnectionProvisioningState = "Disconnected"
	PrivateEndpointConnectionProvisioningStateFailed       PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded    PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating     PrivateEndpointConnectionProvisioningState = "Updating"
)

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved     PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusDisconnected PrivateEndpointServiceConnectionStatus = "Disconnected"
	PrivateEndpointServiceConnectionStatusPending      PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected     PrivateEndpointServiceConnectionStatus = "Rejected"
)

type Reason string

const (
	ReasonAccountNameInvalid Reason = "AccountNameInvalid"
	ReasonAlreadyExists      Reason = "AlreadyExists"
)

type SecretPermissions string

const (
	SecretPermissionsAll     SecretPermissions = "all"
	SecretPermissionsBackup  SecretPermissions = "backup"
	SecretPermissionsDelete  SecretPermissions = "delete"
	SecretPermissionsGet     SecretPermissions = "get"
	SecretPermissionsList    SecretPermissions = "list"
	SecretPermissionsPurge   SecretPermissions = "purge"
	SecretPermissionsRecover SecretPermissions = "recover"
	SecretPermissionsRestore SecretPermissions = "restore"
	SecretPermissionsSet     SecretPermissions = "set"
)

type SkuFamily string

const (
	SkuFamilyA SkuFamily = "A"
)

type SkuName string

const (
	SkuNamePremium  SkuName = "premium"
	SkuNameStandard SkuName = "standard"
)

type StoragePermissions string

const (
	StoragePermissionsAll           StoragePermissions = "all"
	StoragePermissionsBackup        StoragePermissions = "backup"
	StoragePermissionsDelete        StoragePermissions = "delete"
	StoragePermissionsDeletesas     StoragePermissions = "deletesas"
	StoragePermissionsGet           StoragePermissions = "get"
	StoragePermissionsGetsas        StoragePermissions = "getsas"
	StoragePermissionsList          StoragePermissions = "list"
	StoragePermissionsListsas       StoragePermissions = "listsas"
	StoragePermissionsPurge         StoragePermissions = "purge"
	StoragePermissionsRecover       StoragePermissions = "recover"
	StoragePermissionsRegeneratekey StoragePermissions = "regeneratekey"
	StoragePermissionsRestore       StoragePermissions = "restore"
	StoragePermissionsSet           StoragePermissions = "set"
	StoragePermissionsSetsas        StoragePermissions = "setsas"
	StoragePermissionsUpdate        StoragePermissions = "update"
)

type Type string

const (
	TypeMicrosoftPointKeyVaultVaults Type = "Microsoft.KeyVault/vaults"
)

type VaultProvisioningState string

const (
	VaultProvisioningStateRegisteringDns VaultProvisioningState = "RegisteringDns"
	VaultProvisioningStateSucceeded      VaultProvisioningState = "Succeeded"
)
