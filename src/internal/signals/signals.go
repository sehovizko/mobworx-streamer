package signals

import (
	"github.com/google/uuid"
	"github.com/sehovizko/mobworx-streamer/src/internal/helpers"
	"time"
)

func generateTestDataGeneralShape(dgsId, playlistId, variantId, renditionId, renditionGroupId, segmentId, misId, partId string) *DataGeneralShape {
	return &DataGeneralShape{
		Action:  DataActionUpdateVariant,
		Version: 1,
		Id:      uuid.MustParse(dgsId),
		Timestamp: helpers.Timestamp{
			Time: time.Unix(1676898433, 0),
		},
		NumBytes: 2048,
		Payload: &DataGeneralShapePayload{
			Playlist: &DataGeneralShapePayloadPlaylist{
				Id:      uuid.MustParse(playlistId),
				Version: 1,
			},
			Variant: &DataGeneralShapePayloadVariant{
				Id:                 variantId,
				Codecs:             "avc1.4dc00d, mp4a.40.2",
				Bandwidth:          2048,
				Audio:              "example-variant-audio",
				Version:            1,
				TargetDuration:     12,
				TargetPartDuration: 12.232323,
				CacheKey:           playlistId + "/" + variantId,
				InitCacheKey:       playlistId + "/" + misId,
			},
			Rendition: &DataGeneralShapePayloadRendition{
				Id:                 uuid.MustParse(renditionId),
				Type:               DataRenditionTypeAudio,
				GroupId:            uuid.MustParse(renditionGroupId),
				Name:               "audio-en",
				Language:           "en",
				IsDefault:          true,
				AutoSelect:         true,
				TargetDuration:     123,
				TargetPartDuration: 123.123123123,
				CacheKey:           playlistId + "/" + renditionId,
				InitCacheKey:       playlistId + "/" + misId,
			},
			Segment: &DataGeneralShapePayloadSegment{
				Id:              segmentId,
				Sequence:        0,
				Duration:        123.123123123,
				Discontinuity:   true,
				ProgramDateTime: helpers.Timestamp{Time: time.Unix(1676898433, 0)},
				Map: &MediaInitializationSection{
					Id:   misId,
					Data: "qewqweqweqwe",
				},
				Data:     "qwedqweqweq",
				CacheKey: playlistId + "/" + segmentId,
			},
			Part: &DataGeneralShapePayloadPart{
				Id:          partId,
				Sequence:    0,
				Duration:    123.123,
				Independent: true,
				Gap:         true,
				Data:        "qweqweqwe",
				CacheKey:    playlistId + "/" + partId,
			},
		},
	}
}
