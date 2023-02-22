package data

import (
	"github.com/sehovizko/mobworx-streamer/src/internal/errspec"
	"strconv"
)

type Rendition struct {
	*RenditionProps
}

type RenditionProps struct {
	RenditionType   `json:"type"`
	Uri             string   `json:"uri"`
	GroupId         string   `json:"groupId"`
	Language        string   `json:"language"`
	AssocLanguage   string   `json:"assocLanguage"`
	Name            string   `json:"name"`
	IsDefault       bool     `json:"isDefault"`
	AutoSelect      bool     `json:"autoselect"`
	Forced          bool     `json:"forced"`
	InStreamId      string   `json:"instreamId"`
	Characteristics []string `json:"characteristics"`
	Channels        []string `json:"channels"`
}

type RenditionType string

const (
	RenditionTypeAudio          RenditionType = "AUDIO"
	RenditionTypeSubtitles      RenditionType = "SUBTITLES"
	RenditionTypeClosedCaptions RenditionType = "CLOSED-CAPTIONS"
)

func (r Rendition) Validate() error {
	if r.RenditionType == "" {
		return errspec.ParameterIsUndefined("Rendition.RenditionType")
	}

	if r.GroupId == "" {
		return errspec.ParameterIsUndefined("Rendition.GroupId")
	}

	if r.Name == "" {
		return errspec.ParameterIsUndefined("Rendition.Name")
	}

	if r.RenditionType == RenditionTypeSubtitles && r.Uri == "" {
		return errspec.ParameterIsUndefined("Rendition.Uri")
	}

	if r.RenditionType == RenditionTypeClosedCaptions && r.InStreamId == "" {
		return errspec.ParameterIsUndefined("Rendition.InStreamId")
	}

	if r.RenditionType == RenditionTypeClosedCaptions && r.Uri != "" {
		return errspec.ParameterShouldBeNull("Rendition.Uri", r.Uri)
	}

	if r.RenditionType == RenditionTypeClosedCaptions && r.Forced {
		return errspec.InvalidParameter("Rendition.Forced", strconv.FormatBool(r.Forced))
	}

	return nil
}

func NewRendition(props *RenditionProps) (*Rendition, error) {
	r := &Rendition{props}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return r, nil
}
