package helpers

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
)

type Utils struct {
	signalAction       *SignalAction
	streamerConnection *StreamerConnection
	s3UserBucked       string
}

func GetUtils(streamerConnection *StreamerConnection) *Utils {
	return &Utils{
		signalAction:       GetSignalAction(),
		streamerConnection: streamerConnection,
		s3UserBucked:       os.Getenv("S3_USER_BUCKET"),
	}

}

type message struct {
	payload   []byte
	timeStamp string
	action    string
}

func (u *Utils) GetAckSignal(event events.ALBTargetGroupRequest, uploadLatency string) {
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
	//TODO: Check the payload from js object and create same here
	//TODO: Then continue from here to develop.
	switch msg.action {
	case u.signalAction.UpdateVariant:
		ackAction = u.signalAction.AckVariant

	}

}
func (u *Utils) DumpToS3(key string, data []byte) (*s3.PutObjectOutput, error) {
	mySession := session.Must(session.NewSession())
	svc := s3.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
	putObject := &s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   aws.ReadSeekCloser(strings.NewReader(string(data))),
		Bucket: aws.String(u.s3UserBucked),
	}
	return svc.PutObject(putObject)
}
