package iotconnectors

type IotFhirDestinationProperties struct {
	FhirMapping                    IotMappingProperties      `json:"fhirMapping"`
	FhirServiceResourceId          string                    `json:"fhirServiceResourceId"`
	ProvisioningState              *ProvisioningState        `json:"provisioningState,omitempty"`
	ResourceIdentityResolutionType IotIdentityResolutionType `json:"resourceIdentityResolutionType"`
}
