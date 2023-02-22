package signals

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sehovizko/mobworx-streamer/src/internal/errspec"
	"github.com/sehovizko/mobworx-streamer/src/internal/helpers"
	"time"
)

type DataGeneralShape struct {
	Action    DataAction               `json:"action"`
	Version   int                      `json:"version"`
	Id        uuid.UUID                `json:"id"`
	Timestamp helpers.Timestamp        `json:"timestamp"`
	NumBytes  int                      `json:"numbytes"`
	UserId    uuid.UUID                `json:"UserId"`
	Payload   *DataGeneralShapePayload `json:"payload"`
}

type DataGeneralShapePayload struct {
	Playlist  *DataGeneralShapePayloadPlaylist  `json:"playlist,omitempty"`
	Variant   *DataGeneralShapePayloadVariant   `json:"variant,omitempty"`
	Rendition *DataGeneralShapePayloadRendition `json:"rendition,omitempty"`
	Segment   *DataGeneralShapePayloadSegment   `json:"segment"`
	Part      *DataGeneralShapePayloadPart      `json:"part,omitempty"`
}

type DataGeneralShapePayloadPlaylist struct {
	Id      uuid.UUID `json:"id"`
	Version int       `json:"version,omitempty"`
}

type DataGeneralShapePayloadVariant struct {
	Id                 string  `json:"id"`
	Codecs             string  `json:"codecs,omitempty"`
	Bandwidth          int64   `json:"bandwidth,omitempty"`
	Audio              string  `json:"audio,omitempty"`
	Version            int     `json:"version,omitempty"`
	TargetDuration     int     `json:"targetDuration,omitempty"`
	TargetPartDuration float64 `json:"targetPartDuration,omitempty"`
	CacheKey           string  `json:"cacheKey,omitempty"`
	InitCacheKey       string  `json:"initCacheKey,omitempty"`
}

type DataGeneralShapePayloadRendition struct {
	Id                 uuid.UUID         `json:"id"`
	Type               DataRenditionType `json:"type"`
	GroupId            uuid.UUID         `json:"groupId"`
	Name               string            `json:"name"`
	Language           string            `json:"language"`
	IsDefault          bool              `json:"isDefault"`
	AutoSelect         bool              `json:"autoSelect"`
	TargetDuration     int               `json:"targetDuration"`
	TargetPartDuration float64           `json:"targetPartDuration"`
	CacheKey           string            `json:"cacheKey,omitempty"`
	InitCacheKey       string            `json:"initCacheKey,omitempty"`
}

type DataGeneralShapePayloadSegment struct {
	Id              string                      `json:"id,omitempty"`
	Sequence        int                         `json:"sequence,omitempty"`
	Duration        float64                     `json:"duration,omitempty"`
	Discontinuity   bool                        `json:"discontinuity,omitempty"`
	ProgramDateTime helpers.Timestamp           `json:"programDateTime,omitempty"`
	Map             *MediaInitializationSection `json:"map,omitempty"`
	Data            string                      `json:"data,omitempty"`
	CacheKey        string                      `json:"cacheKey,omitempty"`
}

type MediaInitializationSection struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type DataGeneralShapePayloadPart struct {
	Id          string  `json:"id"`
	Sequence    int     `json:"sequence"`
	Duration    float64 `json:"duration"`
	Independent bool    `json:"independent,omitempty"`
	Gap         bool    `json:"gap,omitempty"`
	Data        string  `json:"data"`
	CacheKey    string  `json:"cacheKey,omitempty"`
}

type DataAction string

const (
	DataActionUpdatePart      DataAction = "updatePart"
	DataActionUpdateRendition DataAction = "updateRendition"
	DataActionUpdateSegment   DataAction = "updateSegment"
	DataActionUpdateVariant   DataAction = "updateVariant"
)

type DataRenditionType string

const (
	DataRenditionTypeAudio          DataRenditionType = "AUDIO"
	DataRenditionTypeSubtitles      DataRenditionType = "SUBTITLES"
	DataRenditionTypeClosedCaptions DataRenditionType = "CLOSED-CAPTIONS"
)

var (
	ErrNoTimestampFound = fmt.Errorf("%d: no timestamp found", 400)
)

func NewDataMessage(message string, encoded bool) (*DataGeneralShape, error) {
	if encoded {
		decodedMessage, err := base64.StdEncoding.DecodeString(message)
		if err != nil {
			return nil, err
		}
		return NewDataMessageFromBuffer(decodedMessage)
	}
	return NewDataMessageFromBuffer([]byte(message))
}

func NewDataMessageFromBuffer(buffer []byte) (*DataGeneralShape, error) {
	dgs := &DataGeneralShape{}
	err := json.Unmarshal(buffer, dgs)
	if err != nil {
		return nil, err
	}

	masterPlaylistId := dgs.Payload.Playlist.Id.String()

	if dgs.Payload.Variant != nil {
		dgs.Payload.Variant.CacheKey = masterPlaylistId + "/" + dgs.Payload.Variant.Id
		if dgs.Payload.Segment.Map != nil {
			dgs.Payload.Variant.InitCacheKey = masterPlaylistId + "/" + dgs.Payload.Segment.Map.Id
		}
	}

	if dgs.Payload.Rendition != nil {
		dgs.Payload.Rendition.CacheKey = masterPlaylistId + "/" + dgs.Payload.Rendition.Id.String()
		if dgs.Payload.Segment.Map != nil {
			dgs.Payload.Rendition.InitCacheKey = masterPlaylistId + "/" + dgs.Payload.Segment.Map.Id
		}
	}

	if dgs.Payload.Segment != nil {
		dgs.Payload.Segment.CacheKey = masterPlaylistId + "/" + dgs.Payload.Segment.Id
	}

	if dgs.Payload.Part != nil {
		dgs.Payload.Part.CacheKey = masterPlaylistId + "/" + dgs.Payload.Part.Id
	}

	return dgs, nil
}

func (s DataGeneralShape) UploadLatencyFromNow() (int64, error) {
	if !s.Timestamp.IsZero() {
		return time.Now().Sub(s.Timestamp.Time).Milliseconds(), nil
	}
	return 0, ErrNoTimestampFound
}

func (s DataGeneralShape) GetMediaPlaylistId() (string, error) {
	if s.Payload == nil {
		return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.Id, nil
	}

	if s.Payload.Rendition != nil {
		return s.Payload.Rendition.Id.String(), nil
	}

	return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant and GeneralShape.Payload.Rendition")
}

func (s DataGeneralShape) GetMediaPlaylistCacheKey() (string, error) {
	if s.Payload == nil {
		return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.CacheKey, nil
	}

	if s.Payload.Rendition != nil {
		return s.Payload.Rendition.CacheKey, nil
	}

	return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant and GeneralShape.Payload.Rendition")
}

func (s DataGeneralShape) GetMediaPlaylistInitCacheKey() (string, error) {
	if s.Payload == nil {
		return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.InitCacheKey, nil
	}

	if s.Payload.Rendition != nil {
		return s.Payload.Rendition.InitCacheKey, nil
	}

	return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant and GeneralShape.Payload.Rendition")
}

func (s DataGeneralShape) GetMediaPlaylistTargetDuration() (int, error) {
	if s.Payload == nil {
		return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.TargetDuration, nil
	}

	if s.Payload.Rendition != nil {
		return s.Payload.Rendition.TargetDuration, nil
	}

	return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant and GeneralShape.Payload.Rendition")
}

func (s DataGeneralShape) GetMediaPlaylistTargetPartDuration() (float64, error) {
	if s.Payload == nil {
		return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.TargetPartDuration, nil
	}

	if s.Payload.Rendition != nil {
		return s.Payload.Rendition.TargetPartDuration, nil
	}

	return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant and GeneralShape.Payload.Rendition")
}

func (s DataGeneralShape) GetMediaPlaylistBandwidth() (int64, error) {
	if s.Payload == nil {
		return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.Bandwidth, nil
	}

	return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant")
}

func (s DataGeneralShape) GetMediaPlaylistAudio() (string, error) {
	if s.Payload == nil {
		return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.Audio, nil
	}

	return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant")
}

func (s DataGeneralShape) GetMediaPlaylistVersion() (int, error) {
	if s.Payload == nil {
		return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Variant != nil {
		return s.Payload.Variant.Version, nil
	}

	return 0, errspec.ParameterIsUndefined("DataGeneralShape.Payload.Variant")
}

func (s DataGeneralShape) GetMasterPlaylistId() (string, error) {
	if s.Payload == nil {
		return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload")
	}

	if s.Payload.Playlist != nil {
		return s.Payload.Playlist.Id.String(), nil
	}

	return "", errspec.ParameterIsUndefined("DataGeneralShape.Payload.Id")
}

func (s DataGeneralShape) GetStorageKey() (string, error) {
	masterPlaylistId, err := s.GetMasterPlaylistId()
	if err != nil {
		return "", err
	}

	mediaPlaylistId, err := s.GetMediaPlaylistId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/storage", masterPlaylistId, mediaPlaylistId), nil
}

func (s DataGeneralShape) GetS3Prefix() (string, error) {
	masterPlaylistId, err := s.GetMasterPlaylistId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.UserId, masterPlaylistId), nil
}

func (s DataGeneralShape) GetInitSectionStorageKey(mediaInitializationSectionUri string) (string, error) {
	s3Prefix, err := s.GetS3Prefix()
	if err != nil {
		return "", err
	}

	mediaPlaylistId, err := s.GetMediaPlaylistId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s3Prefix, mediaPlaylistId, mediaInitializationSectionUri), nil
}
