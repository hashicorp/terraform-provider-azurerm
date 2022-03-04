package iscsitargets

type IscsiTargetCreate struct {
	Id                *string                     `json:"id,omitempty"`
	ManagedBy         *string                     `json:"managedBy,omitempty"`
	ManagedByExtended *[]string                   `json:"managedByExtended,omitempty"`
	Name              *string                     `json:"name,omitempty"`
	Properties        IscsiTargetCreateProperties `json:"properties"`
	Type              *string                     `json:"type,omitempty"`
}
