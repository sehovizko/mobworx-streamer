package signals

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
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
			Expected: generateTestDataGeneralShape(
				"6d2325da-b11f-11ed-afa1-0242ac120002",
				"932ac3aa-b11f-11ed-afa1-0242ac120002",
				"a3e4e680-b11f-11ed-afa1-0242ac120002",
				"d02288ec-b11f-11ed-afa1-0242ac120002",
				"dc5daa10-b11f-11ed-afa1-0242ac120002",
				"a8652304-b120-11ed-afa1-0242ac120002",
				"c9258c1e-b120-11ed-afa1-0242ac120002",
				"d9c836d4-b120-11ed-afa1-0242ac120002",
			),
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

func TestNewDataMessage(t *testing.T) {
	cases := []struct {
		Value    string
		Encoded  bool
		Expected *DataGeneralShape
	}{
		{
			Value:   `ew0KCSJhY3Rpb24iOiAidXBkYXRlVmFyaWFudCIsDQoJInZlcnNpb24iOiAxLA0KCSJpZCI6ICI2ZDIzMjVkYS1iMTFmLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLA0KCSJ0aW1lc3RhbXAiOiAxNjc2ODk4NDMzLA0KCSJudW1ieXRlcyI6IDIwNDgsDQoJInBheWxvYWQiOiB7DQoJCSJwbGF5bGlzdCI6IHsNCgkJCSJpZCI6ICI5MzJhYzNhYS1iMTFmLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLA0KCQkJInZlcnNpb24iOiAxDQoJCX0sDQoJCSJ2YXJpYW50Ijogew0KCQkJImlkIjogImEzZTRlNjgwLWIxMWYtMTFlZC1hZmExLTAyNDJhYzEyMDAwMiIsDQoJCQkiY29kZWNzIjogImF2YzEuNGRjMDBkLCBtcDRhLjQwLjIiLA0KCQkJImJhbmR3aWR0aCI6IDIwNDgsDQoJCQkiYXVkaW8iOiAiZXhhbXBsZS12YXJpYW50LWF1ZGlvIiwNCgkJCSJ2ZXJzaW9uIjogMSwNCgkJCSJ0YXJnZXREdXJhdGlvbiI6IDEyLA0KCQkJInRhcmdldFBhcnREdXJhdGlvbiI6IDEyLjIzMjMyMw0KCQl9LA0KCQkicmVuZGl0aW9uIjogew0KCQkJImlkIjogImQwMjI4OGVjLWIxMWYtMTFlZC1hZmExLTAyNDJhYzEyMDAwMiIsDQoJCQkidHlwZSI6ICJBVURJTyIsDQoJCQkiZ3JvdXBJZCI6ICJkYzVkYWExMC1iMTFmLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLA0KCQkJIm5hbWUiOiAiYXVkaW8tZW4iLA0KCQkJImxhbmd1YWdlIjogImVuIiwNCgkJCSJpc0RlZmF1bHQiOiB0cnVlLA0KCQkJImF1dG9TZWxlY3QiOiB0cnVlLA0KCQkJInRhcmdldER1cmF0aW9uIjogMTIzLA0KCQkJInRhcmdldFBhcnREdXJhdGlvbiI6IDEyMy4xMjMxMjMxMjMNCgkJfSwNCgkJInNlZ21lbnQiOiB7DQoJCQkiaWQiOiAiYTg2NTIzMDQtYjEyMC0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyIiwNCgkJCSJzZXF1ZW5jZSI6IDAsDQoJCQkiZHVyYXRpb24iOiAxMjMuMTIzMTIzMTIzLA0KCQkJImRpc2NvbnRpbnVpdHkiOiB0cnVlLA0KCQkJInByb2dyYW1EYXRlVGltZSI6IDE2NzY4OTg0MzMsDQoJCQkibWFwIjogew0KCQkJCSJpZCI6ICJjOTI1OGMxZS1iMTIwLTExZWQtYWZhMS0wMjQyYWMxMjAwMDIiLA0KCQkJCSJkYXRhIjogInFld3F3ZXF3ZXF3ZSINCgkJCX0sDQoJCQkiZGF0YSI6ICJxd2VkcXdlcXdlcSINCgkJfSwNCgkJInBhcnQiOiB7DQoJCQkiaWQiOiAiZDljODM2ZDQtYjEyMC0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyIiwNCgkJCSJzZXF1ZW5jZSI6IDAsDQoJCQkiZHVyYXRpb24iOiAxMjMuMTIzLA0KCQkJImluZGVwZW5kZW50IjogdHJ1ZSwNCgkJCSJnYXAiOiB0cnVlLA0KCQkJImRhdGEiOiAicXdlcXdlcXdlIg0KCQl9DQoJfQ0KfQ==`,
			Encoded: true,
			Expected: generateTestDataGeneralShape(
				"6d2325da-b11f-11ed-afa1-0242ac120002",
				"932ac3aa-b11f-11ed-afa1-0242ac120002",
				"a3e4e680-b11f-11ed-afa1-0242ac120002",
				"d02288ec-b11f-11ed-afa1-0242ac120002",
				"dc5daa10-b11f-11ed-afa1-0242ac120002",
				"a8652304-b120-11ed-afa1-0242ac120002",
				"c9258c1e-b120-11ed-afa1-0242ac120002",
				"d9c836d4-b120-11ed-afa1-0242ac120002",
			),
		},
	}

	for i, c := range cases {
		t.Run("Case/"+strconv.Itoa(i+1), func(t *testing.T) {
			got, err := NewDataMessage(c.Value, c.Encoded)
			require.NoError(t, err)
			assert.Equal(t, got, c.Expected)
		})
	}
}

func TestDataGeneralShape_UploadLatencyFromNow(t *testing.T) {
	dgs := generateTestDataGeneralShape(
		"6d2325da-b11f-11ed-afa1-0242ac120002",
		"932ac3aa-b11f-11ed-afa1-0242ac120002",
		"d02288ec-b11f-11ed-afa1-0242ac120002",
		"d02288ec-b11f-11ed-afa1-0242ac120002",
		"dc5daa10-b11f-11ed-afa1-0242ac120002",
		"a8652304-b120-11ed-afa1-0242ac120002",
		"c9258c1e-b120-11ed-afa1-0242ac120002",
		"d9c836d4-b120-11ed-afa1-0242ac120002",
	)
	now, err := dgs.UploadLatencyFromNow()
	require.NoError(t, err)
	assert.Greater(t, now, time.Duration(10206304))
}

func TestNewDataMessage_FromResources(t *testing.T) {
	cases := []struct {
		Filename       string
		ExpectedAction DataAction
	}{
		{
			Filename:       "data_update_variant_init.mp4.json",
			ExpectedAction: DataActionUpdateVariant,
		},
		{
			Filename:       "data_update_part_index_0.m4s.json",
			ExpectedAction: DataActionUpdatePart,
		},
		{
			Filename:       "data_update_part_index_1.m4s.json",
			ExpectedAction: DataActionUpdatePart,
		},
	}

	for i, c := range cases {
		t.Run("Case/"+strconv.Itoa(i+1), func(t *testing.T) {
			open, err := os.Open(c.Filename)
			require.NoError(t, err)
			defer open.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, open)
			require.NoError(t, err)

			dsg, err := NewDataMessageFromBuffer(buf.Bytes())
			require.NoError(t, err)

			assert.Equal(t, dsg.Action, c.ExpectedAction)
		})
	}
}
