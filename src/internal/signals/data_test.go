package signals

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hls.streaming.com/src/internal/helpers"
	"strconv"
	"testing"
	"time"
)

func TestNewDataMessageFromBuffer(t *testing.T) {
	cases := []struct {
		Value    string
		Expected *DataGeneralShape
	}{
		{
			Value: `
{
	"action": "updateVariant",
	"version": 1,
	"id": "6d2325da-b11f-11ed-afa1-0242ac120002",
	"timestamp": 1676898433,
	"numbytes": 2048,
	"payload": {
		"playlist": {
			"id": "932ac3aa-b11f-11ed-afa1-0242ac120002",
			"version": 1
		},
		"variant": {
			"id": "a3e4e680-b11f-11ed-afa1-0242ac120002",
			"codecs": "avc1.4dc00d, mp4a.40.2",
			"bandwidth": 2048,
			"audio": "example-variant-audio",
			"version": 1,
			"targetDuration": 12,
			"targetPartDuration": 12.232323
		},
		"rendition": {
			"id": "d02288ec-b11f-11ed-afa1-0242ac120002",
			"type": "AUDIO",
			"groupId": "dc5daa10-b11f-11ed-afa1-0242ac120002",
			"name": "audio-en",
			"language": "en",
			"isDefault": true,
			"autoSelect": true,
			"targetDuration": 123,
			"targetPartDuration": 123.123123123
		},
		"segment": {
			"id": "a8652304-b120-11ed-afa1-0242ac120002",
			"sequence": 0,
			"duration": 123.123123123,
			"discontinuity": true,
			"programDateTime": 1676898433,
			"map": {
				"id": "c9258c1e-b120-11ed-afa1-0242ac120002",
				"data": "qewqweqweqwe"
			},
			"data": "qwedqweqweq"
		},
		"part": {
			"id": "d9c836d4-b120-11ed-afa1-0242ac120002",
			"sequence": 0,
			"duration": 123.123,
			"independent": true,
			"gap": true,
			"data": "qweqweqwe"
		}
	}
}`,
			Expected: &DataGeneralShape{
				Action:  DataActionUpdateVariant,
				Version: 1,
				Id:      uuid.MustParse("6d2325da-b11f-11ed-afa1-0242ac120002"),
				Timestamp: helpers.Timestamp{
					Time: time.Unix(1676898433, 0),
				},
				NumBytes: 2048,
				Payload: &DataGeneralShapePayload{
					Playlist: &DataGeneralShapePayloadPlaylist{
						Id:      uuid.MustParse("932ac3aa-b11f-11ed-afa1-0242ac120002"),
						Version: 1,
					},
					Variant: &DataGeneralShapePayloadVariant{
						Id:                 uuid.MustParse("a3e4e680-b11f-11ed-afa1-0242ac120002"),
						Codecs:             "avc1.4dc00d, mp4a.40.2",
						Bandwidth:          2048,
						Audio:              "example-variant-audio",
						Version:            1,
						TargetDuration:     12,
						TargetPartDuration: 12.232323,
						CacheKey:           "932ac3aa-b11f-11ed-afa1-0242ac120002/a3e4e680-b11f-11ed-afa1-0242ac120002",
						InitCacheKey:       "932ac3aa-b11f-11ed-afa1-0242ac120002/c9258c1e-b120-11ed-afa1-0242ac120002",
					},
					Rendition: &DataGeneralShapePayloadRendition{
						Id:                 uuid.MustParse("d02288ec-b11f-11ed-afa1-0242ac120002"),
						Type:               DataRenditionTypeAudio,
						GroupId:            uuid.MustParse("dc5daa10-b11f-11ed-afa1-0242ac120002"),
						Name:               "audio-en",
						Language:           "en",
						IsDefault:          true,
						AutoSelect:         true,
						TargetDuration:     123,
						TargetPartDuration: 123.123123123,
						CacheKey:           "932ac3aa-b11f-11ed-afa1-0242ac120002/d02288ec-b11f-11ed-afa1-0242ac120002",
						InitCacheKey:       "932ac3aa-b11f-11ed-afa1-0242ac120002/c9258c1e-b120-11ed-afa1-0242ac120002",
					},
					Segment: &DataGeneralShapePayloadSegment{
						Id:              uuid.MustParse("a8652304-b120-11ed-afa1-0242ac120002"),
						Sequence:        0,
						Duration:        123.123123123,
						Discontinuity:   true,
						ProgramDateTime: helpers.Timestamp{Time: time.Unix(1676898433, 0)},
						Map: &MediaInitializationSection{
							Id:   uuid.MustParse("c9258c1e-b120-11ed-afa1-0242ac120002"),
							Data: "qewqweqweqwe",
						},
						Data:     "qwedqweqweq",
						CacheKey: "932ac3aa-b11f-11ed-afa1-0242ac120002/a8652304-b120-11ed-afa1-0242ac120002",
					},
					Part: &DataGeneralShapePayloadPart{
						Id:          uuid.MustParse("d9c836d4-b120-11ed-afa1-0242ac120002"),
						Sequence:    0,
						Duration:    123.123,
						Independent: true,
						Gap:         true,
						Data:        "qweqweqwe",
						CacheKey:    "932ac3aa-b11f-11ed-afa1-0242ac120002/d9c836d4-b120-11ed-afa1-0242ac120002",
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run("Case/"+strconv.Itoa(i+1), func(t *testing.T) {
			got, err := NewDataMessageFromBuffer([]byte(c.Value))
			require.NoError(t, err)
			assert.Equal(t, got, c.Expected)
		})
	}
}
