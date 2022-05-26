package operationsinacluster

type OperationStatusResultOperationPredicate struct {
	Id     *string
	Name   *string
	Status *string
}

func (p OperationStatusResultOperationPredicate) Matches(input OperationStatusResult) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Status != nil && *p.Status != input.Status {
		return false
	}

	return true
}
