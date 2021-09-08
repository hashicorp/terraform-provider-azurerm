package tf

import "fmt"

// todo this should be moved to internal somewhere?
func ImportAsExistsError(resourceName, id string) error {
	msg := "A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return fmt.Errorf(msg, id, resourceName)
}
