package data

import "github.com/sehovizko/mobworx-streamer/src/internal/errspec"

type Variant struct {
	*VariantProps
}

type VariantProps struct {
	Uri               string   `json:"uri"`
	IsIFrameOnly      bool     `json:"isIFrameOnly"`
	Bandwidth         int64    `json:"bandwidth"`
	AverageBandwidth  int64    `json:"averageBandwidth"`
	Codecs            string   `json:"codecs"`
	Resolution        string   `json:"resolution"`
	FrameRate         string   `json:"frameRate"`
	HdcpLevel         int      `json:"hdcpLevel"`
	Audio             []string `json:"audio"`
	Video             []string `json:"video"`
	Subtitles         []string `json:"subtitles"`
	ClosedCaptions    []string `json:"closedCaptions"`
	CurrentRenditions `json:"currentRenditions"`
}

type CurrentRenditions struct {
	Audio          int `json:"audio"`
	Video          int `json:"video"`
	Subtitles      int `json:"subtitles"`
	ClosedCaptions int `json:"closedCaptions"`
}

func (v Variant) Validate() error {
	if v.Uri == "" {
		return errspec.ParameterIsUndefined("Variant.Uri")
	}

	if v.Bandwidth == 0 {
		return errspec.ParameterIsUndefined("Variant.Bandwidth")
	}

	return nil
}

func NewVariant(props *VariantProps) (*Variant, error) {
	v := &Variant{props}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return v, nil
}
