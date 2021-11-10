package videoanalyzer

type AccessPolicyEccAlgo string

const (
	AccessPolicyEccAlgoESFiveOneTwo     AccessPolicyEccAlgo = "ES512"
	AccessPolicyEccAlgoESThreeEightFour AccessPolicyEccAlgo = "ES384"
	AccessPolicyEccAlgoESTwoFiveSix     AccessPolicyEccAlgo = "ES256"
)

type AccessPolicyRole string

const (
	AccessPolicyRoleReader AccessPolicyRole = "Reader"
)

type AccessPolicyRsaAlgo string

const (
	AccessPolicyRsaAlgoRSFiveOneTwo     AccessPolicyRsaAlgo = "RS512"
	AccessPolicyRsaAlgoRSThreeEightFour AccessPolicyRsaAlgo = "RS384"
	AccessPolicyRsaAlgoRSTwoFiveSix     AccessPolicyRsaAlgo = "RS256"
)

type AccountEncryptionKeyType string

const (
	AccountEncryptionKeyTypeCustomerKey AccountEncryptionKeyType = "CustomerKey"
	AccountEncryptionKeyTypeSystemKey   AccountEncryptionKeyType = "SystemKey"
)

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

type VideoAnalyzerEndpointType string

const (
	VideoAnalyzerEndpointTypeClientApi VideoAnalyzerEndpointType = "ClientApi"
)

type VideoType string

const (
	VideoTypeArchive VideoType = "Archive"
)
