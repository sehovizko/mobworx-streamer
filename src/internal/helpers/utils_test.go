package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestGetAckSignal(t *testing.T) {
	part := &Part{Data: []byte("part data")}
	videoPart := &VideoPart{Data: nil}
	audioPart := &AudioPart{Data: nil}
	mapping := &Mapping{Data: []byte("map Data")}
	videoSegment := &VideoSegment{
		Mapping: mapping,
		Data:    []byte("VideoSegment Data"),
	}
	audioSegment := &AudioSegment{
		Mapping: mapping,
		Data:    []byte("AudioSegment Data"),
	}
	segment := &Segment{
		Mapping: mapping,
		Data:    []byte("Segment Data"),
	}
	payload := &Payload{
		Segment:      segment,
		VideoSegment: videoSegment,
		AudioSegment: audioSegment,
		Part:         part,
		VideoPart:    videoPart,
		AudioPart:    audioPart,
	}
	message := &Message{
		Payload:   payload,
		TimeStamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Action:    UpdatePart,
	}

	mySession := session.Must(session.NewSession())
	utils := NewUtils(mySession, "domain", "stage")

	messageByte, _ := json.Marshal(*message)

	request := &events.APIGatewayProxyRequest{Body: string(messageByte), IsBase64Encoded: false}

	testCases := []struct {
		event   *events.APIGatewayProxyRequest
		latency int
		resp    *Response
	}{
		{
			event:   request,
			latency: 10,
			resp: &Response{
				Action: string(AckPart),
			},
		},
	}

	for i, c := range testCases {
		t.Run(fmt.Sprintf("SubTest %d", i+1), func(t *testing.T) {
			resp := utils.GetAckSignal(c.event, c.latency)
			assert.Equal(t, c.resp.Action, resp.Action)
		})
	}
}

func TestAck(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	part := &Part{Data: []byte("part data")}
	videoPart := &VideoPart{Data: nil}
	audioPart := &AudioPart{Data: nil}
	mapping := &Mapping{Data: []byte("map Data")}
	videoSegment := &VideoSegment{
		Mapping: mapping,
		Data:    []byte("VideoSegment Data"),
	}
	audioSegment := &AudioSegment{
		Mapping: mapping,
		Data:    []byte("AudioSegment Data"),
	}
	segment := &Segment{
		Mapping: mapping,
		Data:    []byte("Segment Data"),
	}
	payload := &Payload{
		Segment:      segment,
		VideoSegment: videoSegment,
		AudioSegment: audioSegment,
		Part:         part,
		VideoPart:    videoPart,
		AudioPart:    audioPart,
	}
	message := &Message{
		Payload:   payload,
		TimeStamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Action:    UpdatePart,
	}

	mySession := session.Must(session.NewSession())
	utils := NewUtils(mySession, "domain", "stage")

	messageByte, _ := json.Marshal(*message)
	request := &events.APIGatewayWebsocketProxyRequest{
		Body:            string(messageByte),
		IsBase64Encoded: false,
		RequestContext: events.APIGatewayWebsocketProxyRequestContext{
			Stage:        "stg",
			ConnectionID: "555",
			DomainName:   "dmn",
		},
	}

	err := utils.Ack(request, 10)
	require.NoError(t, err)
	//Todo: Should write assertions to here.
}

func TestDumpToS3(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	mySession := session.Must(session.NewSession())
	utils := NewUtils(mySession, "domain", "stage")
	response, err := utils.DumpToS3("key", []byte("data"))
	require.NoError(t, err)
	fmt.Println(response)
	//Todo: Should write assertions to here.
}
