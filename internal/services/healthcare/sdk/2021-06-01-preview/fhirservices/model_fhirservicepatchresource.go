package fhirservices

type FhirServicePatchResource struct {
	Identity *ServiceManagedIdentityIdentity `json:"identity,omitempty"`
	Tags     *map[string]string              `json:"tags,omitempty"`
}
