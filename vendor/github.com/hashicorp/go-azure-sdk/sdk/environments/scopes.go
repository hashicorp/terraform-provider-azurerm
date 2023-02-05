package environments

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func Resource(endpoint Api) (*string, error) {
	resource, ok := endpoint.ResourceIdentifier()
	if !ok {
		return nil, fmt.Errorf("the endpoint %q doesn't define a resource identifier", endpoint.Name())
	}

	return resource, nil
}

func Scope(endpoint Api) (*string, error) {
	e, ok := endpoint.ResourceIdentifier()
	if !ok {
		return nil, fmt.Errorf("the endpoint %q is not supported in this Azure Environment", endpoint.Name())
	}
	out := fmt.Sprintf("%s/.default", *e)
	return pointer.To(out), nil
}
