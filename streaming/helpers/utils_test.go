package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
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
	sc := NewStreamerConnection(mySession, "domain", "stage")
	utils := NewUtils(sc)

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
	sc := NewStreamerConnection(mySession, "domain", "stage")
	utils := NewUtils(sc)

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

	utils.Ack(request, 10)
	//Todo: Should write assertions to here.

}
