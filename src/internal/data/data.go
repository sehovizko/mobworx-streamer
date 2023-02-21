package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type Data struct {
	DataType Type `json:"type"`
}

type Type string

const (
	TypePlaylist Type = "playlist"
	TypeSegment  Type = "segment"
)

func (d Data) Validate() error {
	if d.DataType == "" {
		return errspec.ParameterIsUndefined("dataType")
	}
	if d.DataType != TypePlaylist && d.DataType != TypeSegment {
		return errspec.InvalidParameter("dataType", string(d.DataType))
	}
	return nil
}

func NewData(dataType Type) (*Data, error) {
	d := &Data{dataType}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return d, nil
}
