package tftypes

import (
	"fmt"
)

// AttributePathError represents an error associated with part of a
// tftypes.Value, indicated by the Path property.
type AttributePathError struct {
	Path *AttributePath
	err  error
}

func (a AttributePathError) Error() string {
	var path string
	if len(a.Path.Steps()) > 0 {
		path = a.Path.String() + ": "
	}
	return fmt.Sprintf("%s%s", path, a.err)
}

func (a AttributePathError) Unwrap() error {
	return a.err
}
