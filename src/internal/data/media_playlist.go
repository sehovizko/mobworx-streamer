package data

import (
	"github.com/aws/aws-sdk-go/aws"
)

type MediaPlaylist struct {
	*Playlist
	*MediaPlaylistProps
}

type MediaPlaylistProps struct {
	*PlaylistProps
	MediaSequenceBase         int                         `json:"mediaSequenceBase"`
	DiscontinuitySequenceBase int                         `json:"discontinuitySequenceBase"`
	EndList                   bool                        `json:"endlist"`
	PlaylistType              MediaPlaylistType           `json:"playlistType"`
	IsIFrame                  bool                        `json:"isIFrame"`
	Segments                  []string                    `json:"segments"`
	Hash                      string                      `json:"hash"`
	Lock                      bool                        `json:"lock"`
	Map                       *MediaInitializationSection `json:"map"`
	UpdatedAt                 string                      `json:"updatedAt"`
	HoldingSegments           []string                    `json:"holdingSegments"`
	DiscardedSegments         []string                    `json:"discardedSegments"`
	ServerControl             *ServerControl              `json:"serverControl"`
	TargetDuration            int                         `json:"targetDuration"`
	TargetPartDuration        float64                     `json:"targetPartDuration"`
	RenditionReports          []string                    `json:"renditionReports"`
}

type MediaPlaylistType string

const (
	MediaPlaylistTypeVOD  MediaPlaylistType = "VOD"
	MediaPlaylistTypeLive MediaPlaylistType = "LIVE"
)

func (p MediaPlaylist) Validate() error {
	return nil
}

func NewMediaPlaylist(props *MediaPlaylistProps) (*MediaPlaylist, error) {
	props.IsMasterPlaylist = aws.Bool(false)
	p, err := NewPlaylist(props.PlaylistProps)
	if err != nil {
		return nil, err
	}

	mp := &MediaPlaylist{p, props}
	if err := mp.Validate(); err != nil {
		return nil, err
	}

	return mp, err
}
