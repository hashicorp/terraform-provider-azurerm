package datasources

type DataSourceOperationPredicate struct {
	Etag       *string
	Id         *string
	Name       *string
	Properties *interface{}
	Type       *string
}

func (p DataSourceOperationPredicate) Matches(input DataSource) bool {

	if p.Etag != nil && (input.Etag == nil && *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Properties != nil && *p.Properties != input.Properties {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}
