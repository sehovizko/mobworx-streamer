package data

import "github.com/aws/aws-sdk-go/aws"

type MasterPlaylist struct {
	*Playlist
	*MasterPlaylistProps
}

type MasterPlaylistProps struct {
	*PlaylistProps
	Variants        []*Variant `json:"variants"`
	CurrentVariant  *Variant   `json:"currentVariant,omitempty"`
	SessionDataList []string   `json:"sessionDataList"`
	SessionKeyList  []string   `json:"sessionKeyList"`
}

func (p MasterPlaylist) Validate() error {
	return nil
}

func NewMasterPlaylist(props *MasterPlaylistProps) (*MasterPlaylist, error) {
	props.IsMasterPlaylist = aws.Bool(true)
	p, err := NewPlaylist(props.PlaylistProps)
	if err != nil {
		return nil, err
	}

	mp := &MasterPlaylist{p, props}
	if err := mp.Validate(); err != nil {
		return nil, err
	}

	return mp, nil
}
