package data

import (
	"github.com/sehovizko/mobworx-streamer/src/internal/errspec"
)

type Playlist struct {
	*Data
	*PlaylistProps
}

type PlaylistProps struct {
	IsMasterPlaylist    *bool  `json:"isMasterPlaylist"`
	Uri                 string `json:"uri"`
	Version             int    `json:"version"`
	IndependentSegments bool   `json:"independentSegments"`
	Start               string `json:"start"`
	Source              string `json:"source"`
}

func (p Playlist) Validate() error {
	if p.IsMasterPlaylist == nil {
		return errspec.ParameterIsUndefined("isMasterPlaylist")
	}
	return nil
}

func NewPlaylist(props *PlaylistProps) (*Playlist, error) {
	d, err := NewData(TypePlaylist)
	if err != nil {
		return nil, err
	}

	p := &Playlist{d, props}
	err = p.Validate()
	if err != nil {
		return nil, err
	}

	return p, nil
}
