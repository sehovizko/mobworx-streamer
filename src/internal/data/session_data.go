package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type SessionData struct {
	*SessionDataProps
}

type SessionDataProps struct {
	Id       string `json:"id"`
	Value    string `json:"value"`
	Uri      string `json:"uri"`
	Language string `json:"language"`
}

func (d SessionData) Validate() error {
	if d.Id == "" {
		return errspec.ParameterIsUndefined("SessionData.Id")
	}

	if d.Value != "" && d.Uri != "" {
		return errspec.InvalidParameter("both SessionData.Value and SessionData.Uri", d.Value+"/"+d.Uri)
	}

	if !(d.Value != "" || d.Uri != "") {
		return errspec.ParameterIsUndefined("both SessionData.Value and SessionData.Uri is undefined")
	}

	return nil
}

func NewSessionData(props *SessionDataProps) (*SessionData, error) {
	sd := &SessionData{props}
	if err := sd.Validate(); err != nil {
		return nil, err
	}
	return sd, nil
}
