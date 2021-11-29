package managedcluster

type ApplicationTypeVersionsCleanupPolicy struct {
	MaxUnusedVersionsToKeep int64 `json:"maxUnusedVersionsToKeep"`
}
