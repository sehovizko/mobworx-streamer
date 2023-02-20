package helpers

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Utils interface {
	GetAckSignal(event *events.APIGatewayProxyRequest, uploadLatency int) *Response
	Ack(event *events.APIGatewayWebsocketProxyRequest, uploadLatency int) error
	DumpToS3(key string, data []byte) (*s3.PutObjectOutput, error)
}

type utils struct {
	S3Session          *s3.S3
	S3UserBucket       string
	StreamerConnection *StreamerConnection
}

func NewUtils(session *session.Session, domain, stage string) Utils {
	return &utils{
		S3Session:          s3.New(session, aws.NewConfig().WithRegion("us-west-2")),
		S3UserBucket:       os.Getenv("S3_USER_BUCKET"),
		StreamerConnection: NewStreamerConnection(session, domain, stage),
	}
}

type Message struct {
	Payload   *Payload
	TimeStamp string
	Action    SignalAction
}

type Payload struct {
	Segment      *Segment
	VideoSegment *VideoSegment
	AudioSegment *AudioSegment
	Part         *Part
	VideoPart    *VideoPart
	AudioPart    *AudioPart
}

type Segment struct {
	Mapping *Mapping
	Data    []byte
}

type VideoSegment struct {
	Mapping *Mapping
	Data    []byte
}

type AudioSegment struct {
	Mapping *Mapping
	Data    []byte
}

type Part struct {
	Data []byte
}

type VideoPart struct {
	Data []byte
}

type AudioPart struct {
	Data []byte
}

type Mapping struct {
	Data []byte
}

type Response struct {
	Action    string
	Version   int
	Id        string
	Timestamp string
	Size      int
	Latency   int
	Payload   *Payload
}

func (u *utils) GetAckSignal(event *events.APIGatewayProxyRequest, uploadLatency int) *Response {
	msg := new(Message)
	body := []byte(event.Body)

	var err error
	if event.IsBase64Encoded {
		//Todo:This converting is possibly wrong.
		data := b64.StdEncoding.EncodeToString(body)
		err = json.Unmarshal([]byte(data), &msg)
	} else {
		err = json.Unmarshal(body, &msg)
	}
	if err != nil {
		panic(err)
	}

	payload := msg.Payload
	ackAction := "unknown"
	switch msg.Action {
	case UpdateVariant:
		ackAction = string(AckVariant)
		payload.Segment.Mapping = nil
		break
	case UpdateRendition:
		ackAction = string(AckRendition)
		payload.Segment.Mapping = nil
		break
	case UpdateSegment:
		ackAction = string(AckSegment)
		payload.Segment.Data = nil
		if payload.Segment.Mapping != nil {
			payload.Segment.Mapping.Data = nil
		}
	case UpdatePart:
		ackAction = string(AckPart)
		payload.Part.Data = nil
		if payload.Segment != nil && payload.Segment.Mapping != nil {
			payload.Segment.Mapping.Data = nil
		}
		break
	case UpdateDemuxSegment:
		ackAction = string(AckDemuxSegment)
		payload.VideoSegment.Data = nil
		if payload.VideoSegment.Mapping != nil {
			payload.VideoSegment.Mapping.Data = nil
		}
		payload.AudioSegment.Data = nil
		if payload.AudioSegment.Mapping != nil {
			payload.AudioSegment.Mapping.Data = nil
		}
		break
	case UpdateDemuxPart:
		ackAction = string(AckDemuxPart)
		payload.VideoPart.Data = nil
		if payload.VideoSegment != nil && payload.VideoSegment.Mapping != nil {
			payload.VideoSegment.Mapping.Data = nil
		}
		payload.AudioPart.Data = nil
		if payload.AudioSegment != nil && payload.AudioSegment.Mapping != nil {
			payload.AudioSegment.Mapping.Data = nil
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
	//Latency := uploadLatency || ((new Date).getTime() - Number(msg.Timestamp));
	return &Response{
		Action:    ackAction,
		Version:   1,
		Id:        uuid.New().String(),
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Size:      size,
		Latency:   uploadLatency, //TODO: Should be Latency variable.
		Payload:   payload,
	}
}

func (u *utils) Ack(event *events.APIGatewayWebsocketProxyRequest, uploadLatency int) error {
	tempData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	proxyRequest := new(events.APIGatewayProxyRequest)
	err = json.Unmarshal(tempData, proxyRequest)
	if err != nil {
		return err
	}

	data := u.GetAckSignal(proxyRequest, uploadLatency)
	connectionId := event.RequestContext.ConnectionID
	connection := NewStreamerConnection(u.StreamerConnection.Session, event.RequestContext.DomainName, event.RequestContext.Stage)
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	startTime := time.Now().UnixMilli()
	_, err = connection.PostData(dataByte, connectionId)
	if err != nil {
		return err
	}
	endTime := time.Now().UnixMilli()

	log.Printf("%+v\n", *data)
	log.Println("ack postData totalTime in ms: ", endTime-startTime)
	return nil
}

func (u *utils) DumpToS3(key string, data []byte) (*s3.PutObjectOutput, error) {
	putObject := &s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   aws.ReadSeekCloser(strings.NewReader(string(data))),
		Bucket: aws.String(u.S3UserBucket),
		Key:    aws.String(key),
	}
	return u.S3Session.PutObject(putObject)
}
