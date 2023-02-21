package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type ServerControl struct {
	*ServerControlProps
}

type ServerControlProps struct {
	CanBlockReload *bool `json:"canBlockReload"`
	CanSkipUntil   bool  `json:"canSkipUntil"`
	HoldBack       bool  `json:"holdBack"`
	PartHoldBack   bool  `json:"partHoldBack"`
}

func (sc ServerControl) Validate() error {
	if sc.CanBlockReload == nil {
		return errspec.ParameterIsUndefined("ServerControl.CanBlockReload")
	}
	return nil
}

func NewServerControl(props *ServerControlProps) (*ServerControl, error) {
	sc := &ServerControl{props}
	if err := sc.Validate(); err != nil {
		return nil, err
	}
	return sc, nil
}
