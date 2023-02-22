package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type Key struct {
	*KeyProps
}

type KeyProps struct {
	Method        KeyMethod `json:"method"`
	Uri           string    `json:"uri"`
	IV            string    `json:"iv"`
	Format        string    `json:"format"`
	FormatVersion string    `json:"formatVersion"`
}

type KeyMethod string

const (
	KeyMethodNone KeyMethod = "NONE"
)

func (k Key) Validate() error {
	if k.Method != KeyMethodNone && k.Uri != "" {
		return errspec.ParameterIsUndefined("Key.Uri")
	}
	// TODO: investigate other cases when Key will be invalid
	return nil
}

func NewKey(props *KeyProps) (*Key, error) {
	key := &Key{props}
	if err := key.Validate(); err != nil {
		return nil, err
	}
	return key, nil
}
