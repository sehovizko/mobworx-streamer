package helpers

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Utils struct {
	StreamerConnection *StreamerConnection
	S3UserBucked       string
}

func CreateUtils(streamerConnection *StreamerConnection) *Utils {
	var utils Utils
	utils.StreamerConnection = new(StreamerConnection)
	utils.StreamerConnection = streamerConnection
	utils.S3UserBucked = os.Getenv("S3_USER_BUCKET")
	return &utils

}

type message struct {
	payload   *Payload
	timeStamp string
	action    SignalAction
}

type Payload struct {
	segment      *segment
	videoSegment *videoSegment
	audioSegment *audioSegment
	part         *part
	videoPart    *videoPart
	audioPart    *audioPart
}
type segment struct {
	mapping *mapping
	data    []byte
}
type videoSegment struct {
	mapping *mapping
	data    []byte
}
type audioSegment struct {
	mapping *mapping
	data    []byte
}
type part struct {
	data []byte
}
type videoPart struct {
	data []byte
}
type audioPart struct {
	data []byte
}
type mapping struct {
	data []byte
}

type Response struct {
	action    string
	version   int
	id        string
	timestamp string
	size      int
	latency   int
	payload   Payload
}

func (u *Utils) GetAckSignal(event *events.APIGatewayProxyRequest, uploadLatency int) *Response {
	msg := new(message)
	var err error
	body := []byte(event.Body)
	if event.IsBase64Encoded {
		data := b64.StdEncoding.EncodeToString(body)
		err = json.Unmarshal([]byte(data), msg)

	} else {
		err = json.Unmarshal(body, &msg)
	}

	if err != nil {
		panic(err)
	}

	payload := msg.payload
	ackAction := "unknown"
	switch msg.action {
	case UpdateVariant:
		ackAction = string(AckVariant)
		payload.segment.mapping = nil
		break
	case UpdateRendition:
		ackAction = string(AckRendition)
		payload.segment.mapping = nil
		break
	case UpdateSegment:
		ackAction = string(AckSegment)
		payload.segment.data = nil
		if payload.segment.mapping != nil {
			payload.segment.mapping.data = nil
		}
	case UpdatePart:
		ackAction = string(AckPart)
		payload.part.data = nil
		if payload.segment != nil && payload.segment.mapping != nil {
			payload.segment.mapping.data = nil
		}
		break
	case UpdateDemuxSegment:
		ackAction = string(AckDemuxSegment)
		payload.videoSegment.data = nil
		if payload.videoSegment.mapping != nil {
			payload.videoSegment.mapping.data = nil
		}
		payload.audioSegment.data = nil
		if payload.audioSegment.mapping != nil {
			payload.audioSegment.mapping.data = nil
		}
		break
	case UpdateDemuxPart:
		ackAction = string(AckDemuxPart)
		payload.videoPart.data = nil
		if payload.videoSegment != nil && payload.videoSegment.mapping != nil {
			payload.videoSegment.mapping.data = nil
		}
		payload.audioPart.data = nil
		if payload.audioSegment != nil && payload.audioSegment.mapping != nil {
			payload.audioSegment.mapping.data = nil
		}
		break
	case Ping:
		ackAction = string(Pong)
		break
	case Abort:
		ackAction = string(Aborted)
		break
	case Terminate:
		ackAction = string(Terminated)
		break
	default:
		break
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	size := len(jsonPayload)
	//latency := uploadLatency || ((new Date).getTime() - Number(msg.timestamp));
	return &Response{
		action:    ackAction,
		version:   1,
		id:        uuid.New().String(),
		timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		size:      size,
		latency:   uploadLatency, //TODO: Should be latency variable.
		payload:   Payload{},
	}
}

func (u *Utils) Ack(event *events.APIGatewayWebsocketProxyRequest, uploadLatency int) {
	tempData, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	proxyRequest := new(events.APIGatewayProxyRequest)
	err = json.Unmarshal(tempData, proxyRequest)
	if err != nil {
		panic(err)
	}
	data := u.GetAckSignal(proxyRequest, uploadLatency)
	connectionId := event.RequestContext.ConnectionID
	connection := NewStreamerConnection(u.StreamerConnection.Session, event.RequestContext.DomainName, event.RequestContext.Stage)
	dataByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	startTime := time.Now().UnixMilli()
	_, err = connection.PostData(dataByte, connectionId)
	endTime := time.Now().UnixMilli()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n", *data)
	log.Println("ack postData totalTime in ms: ", endTime-startTime)

}
func (u *Utils) DumpToS3(key string, data []byte) (*s3.PutObjectOutput, error) {
	svc := s3.New(u.StreamerConnection.Session, aws.NewConfig().WithRegion("us-west-2"))
	putObject := &s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   aws.ReadSeekCloser(strings.NewReader(string(data))),
		Bucket: aws.String(u.S3UserBucked),
	}
	return svc.PutObject(putObject)
}
