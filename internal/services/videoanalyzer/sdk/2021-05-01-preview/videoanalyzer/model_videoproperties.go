package videoanalyzer

type VideoProperties struct {
	Description *string         `json:"description,omitempty"`
	Flags       *VideoFlags     `json:"flags,omitempty"`
	MediaInfo   *VideoMediaInfo `json:"mediaInfo,omitempty"`
	Streaming   *VideoStreaming `json:"streaming,omitempty"`
	Title       *string         `json:"title,omitempty"`
	Type        *VideoType      `json:"type,omitempty"`
}
