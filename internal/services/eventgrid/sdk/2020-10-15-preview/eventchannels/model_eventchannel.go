package eventchannels

type EventChannel struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *EventChannelProperties `json:"properties,omitempty"`
	SystemData *SystemData             `json:"systemData,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
