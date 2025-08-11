// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

// MetadataRequest represents a request for the Resource to return metadata,
// such as its type name. An instance of this request struct is supplied as
// an argument to the Resource type Metadata method.
type MetadataRequest struct {
	// ProviderTypeName is the string returned from
	// [provider.MetadataResponse.TypeName], if the Provider type implements
	// the Metadata method. This string should prefix the Resource type name
	// with an underscore in the response.
	ProviderTypeName string
}

// MetadataResponse represents a response to a MetadataRequest. An
// instance of this response struct is supplied as an argument to the
// Resource type Metadata method.
type MetadataResponse struct {
	// TypeName should be the full resource type, including the provider
	// type prefix and an underscore. For example, examplecloud_thing.
	TypeName string

	// ResourceBehavior is used to control framework-specific logic when
	// interacting with this resource.
	ResourceBehavior ResourceBehavior
}

// ResourceBehavior controls framework-specific logic when interacting
// with a resource.
type ResourceBehavior struct {
	// ProviderDeferred enables provider-defined logic to be executed
	// in the case of an automatic deferred response from provider configure.
	//
	// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
	// to change or break without warning. It is not protected by version compatibility guarantees.
	ProviderDeferred ProviderDeferredBehavior

	// MutableIdentity indicates that the managed resource supports an identity that can change during the
	// resource's lifecycle. Setting this flag to true will disable the SDK validation that ensures identity
	// data doesn't change during RPC calls.
	MutableIdentity bool
}

// ProviderDeferredBehavior enables provider-defined logic to be executed
// in the case of a deferred response from provider configuration.
//
// NOTE: This functionality is related to deferred action support, which is currently experimental and is subject
// to change or break without warning. It is not protected by version compatibility guarantees.
type ProviderDeferredBehavior struct {
	// When EnablePlanModification is true, framework will still execute
	// provider-defined resource plan modification logic if
	// provider.Configure defers. Framework will then automatically return a
	// deferred response along with the modified plan.
	EnablePlanModification bool
}
