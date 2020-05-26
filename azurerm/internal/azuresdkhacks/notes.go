package azuresdkhacks

// There's a functional difference that exists between the Azure SDK for Go and Azure Resource Manager API
// where when performing a delta update unchanged fields are omitted from the response when they could
// also have a legitimate value of `null` (to remove/disable a sub-block).
//
// Ultimately the Azure SDK for Go has opted to serialize structs with `json:"name,omitempty"` which
// means that this value will be omitted if nil to allow for delta updates - however this means there's
// no means of removing/resetting a value of a nested object once provided since a `nil` object will be
// reset
//
// As such this set of well intentioned hacks is intended to force this behaviour where necessary.
//
// It's worth noting that these hacks are a last resort and the Swagger/API/SDK should almost always be
// fixed instead.
