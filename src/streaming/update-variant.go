package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/redis/go-redis/v9"
	"github.com/sehovizko/mobworx-streamer/src/internal/data"
	"github.com/sehovizko/mobworx-streamer/src/internal/signals"
	"log"
	"os"
	"time"
)

var (
	redisClient *redis.Client
	s3Uploader  *s3manager.Uploader
)

func HandleUpdateVariant(ctx aws.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	signal, err := signals.NewDataMessage(event.Body, event.IsBase64Encoded)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("data signal update variant signal: %+v", signal)

	now, err := signal.UploadLatencyFromNow()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("upload latency: %dms", now)

	masterPlaylistId, err := signal.GetMasterPlaylistId()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistId, err := signal.GetMediaPlaylistId()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	masterPlaylistCached, err := redisClient.Get(ctx, masterPlaylistId).Result()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var masterPlaylist *data.MasterPlaylist
	if masterPlaylistCached != "" {
		masterPlaylist, err = data.NewMasterPlaylistFromString(masterPlaylistCached)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		if masterPlaylist.Variants != nil && len(masterPlaylist.Variants) > 0 {
			if masterPlaylist.Variants[0].Uri == mediaPlaylistId {
				log.Printf("Master found: uri=%s", masterPlaylistId)
				return events.APIGatewayProxyResponse{StatusCode: 200}, nil
			}
		}
	} else {
		masterPlaylist, err = data.NewMasterPlaylist(&data.MasterPlaylistProps{
			PlaylistProps: &data.PlaylistProps{
				Uri:                 "playlist.m3u8",
				Version:             signal.Payload.Playlist.Version,
				IndependentSegments: true,
			},
			Variants: make([]*data.Variant, 0),
		})
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	}

	mediaPlaylistCacheKey, err := signal.GetMediaPlaylistCacheKey()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistInitCacheKey, err := signal.GetMediaPlaylistInitCacheKey()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	storageKey, err := signal.GetStorageKey()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaInitializationSection, err := data.NewMediaInitializationSection(&data.MediaInitializationSectionProps{
		Uri:    signal.Payload.Segment.Map.Id,
		Prefix: mediaPlaylistId,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	_, err = redisClient.Set(ctx, mediaPlaylistInitCacheKey, signal.Payload.Segment.Map.Data, time.Hour*10).Result()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistTargetDuration, err := signal.GetMediaPlaylistTargetDuration()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistTargetPartDuration, err := signal.GetMediaPlaylistTargetPartDuration()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	serverControl, err := data.NewServerControl(&data.ServerControlProps{
		CanBlockReload: aws.Bool(true),
		CanSkipUntil:   6 * mediaPlaylistTargetDuration,
		PartHoldBack:   3 * mediaPlaylistTargetPartDuration,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistVersion, err := signal.GetMediaPlaylistVersion()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylist, err := data.NewMediaPlaylist(&data.MediaPlaylistProps{
		PlaylistProps: &data.PlaylistProps{
			Uri:     mediaPlaylistId,
			Version: mediaPlaylistVersion,
		},
		TargetDuration:     mediaPlaylistTargetDuration,
		TargetPartDuration: mediaPlaylistTargetPartDuration,
		MediaSequenceBase:  0,
		PlaylistType:       data.MediaPlaylistTypeLive,
		Map:                mediaInitializationSection,
		UpdatedAt:          time.Now().String(),
		ServerControl:      serverControl,
		Segments:           make([]*data.Segment, 0),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistSerialized, err := json.Marshal(mediaPlaylist)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	_, err = redisClient.Set(ctx, mediaPlaylistCacheKey, mediaPlaylistSerialized, time.Hour*10).Result()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylist.PlaylistType = data.MediaPlaylistTypeVOD
	mediaPlaylist.EndList = true

	mediaPlaylistSerializedChanged, err := json.Marshal(mediaPlaylist)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	_, err = redisClient.Set(ctx, storageKey, mediaPlaylistSerializedChanged, time.Hour*10).Result()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mediaPlaylistAudio, err := signal.GetMediaPlaylistAudio()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	renditions := make([]*data.Rendition, 0)
	if mediaPlaylistAudio != "" {
		rendition, err := data.NewRendition(&data.RenditionProps{
			RenditionType: data.RenditionTypeAudio,
			Name:          "audio-en",
			GroupId:       mediaPlaylistAudio,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		renditions = append(renditions, rendition)
	}

	mediaPlaylistBandwidth, err := signal.GetMediaPlaylistBandwidth()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	variant, err := data.NewVariant(&data.VariantProps{
		Uri:       mediaPlaylistId,
		Codecs:    "avc1.4dc00d, mp4a.40.2",
		Bandwidth: mediaPlaylistBandwidth,
		Audio:     renditions,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// TODO: walk around media playlist averageBandwidth and frameRate props

	masterPlaylist.Variants = append(masterPlaylist.Variants, variant)

	masterPlaylistSerialized, err := json.Marshal(masterPlaylist)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	redisClient.Set(ctx, masterPlaylistId, masterPlaylistSerialized, time.Hour*10)

	initSectionStorageKey, err := signal.GetInitSectionStorageKey(mediaInitializationSection.Uri)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	base64Decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(signal.Payload.Segment.Map.Data))
	_, err = s3Uploader.Upload(&s3manager.UploadInput{
		Body:   base64Decoder,
		ACL:    aws.String("public-read"),
		Bucket: aws.String(os.Getenv("S3_USER_BUCKET")),
		Key:    aws.String(initSectionStorageKey),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
	sess := session.Must(session.NewSession())
	s3Uploader = s3manager.NewUploader(sess)
	lambda.Start(HandleUpdateVariant)
}
