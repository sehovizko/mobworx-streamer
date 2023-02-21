package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type MediaInitializationSection struct {
	*MediaInitializationSectionProps
}

type MediaInitializationSectionProps struct {
	Uri       string `json:"uri"`
	MimeType  string `json:"mimeType"`
	ByteRange string `json:"byteRange"`
	Prefix    string `json:"prefix"`
}

func (m MediaInitializationSection) Validate() error {
	if m.Uri == "" {
		return errspec.ParameterIsUndefined("MediaInitializationSection.uri")
	}
	return nil
}

func NewMediaInitializationSection(props *MediaInitializationSectionProps) (*MediaInitializationSection, error) {
	mis := &MediaInitializationSection{props}
	if err := mis.Validate(); err != nil {
		return nil, err
	}
	return mis, nil
}
