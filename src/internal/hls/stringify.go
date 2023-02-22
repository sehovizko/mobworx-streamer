package hls

import (
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sehovizko/mobworx-streamer/src/internal/data"
	"github.com/sehovizko/mobworx-streamer/src/internal/errspec"
	"math"
	"strings"
)

type Stringify []string

var (
	AllowRedundancy = []string{
		"#EXTINF",
		"#EXT-X-BYTERANGE",
		"#EXT-X-DISCONTINUITY",
		"#EXT-X-MAP",
		"#EXT-X-STREAM-INF",
		// low lat
		"#EXT-X-PART",
		"#EXT-X-RENDITION-REPORT",
		"#EXT-X-GAP",
	}
	SkipIfRedundant = []string{
		"#EXT-X-KEY",
		"#EXT-X-MEDIA",
	}
)

func (s *Stringify) Push(elems ...string) error {
push:
	for _, elem := range elems {
		for _, ar := range AllowRedundancy {
			if strings.HasPrefix(elem, ar) {
				*s = append(*s, elem)
				continue push
			}
		}
		for _, exists := range *s {
			if exists == elem {
				for _, sir := range SkipIfRedundant {
					if strings.HasPrefix(elem, sir) {
						continue push
					}
				}
				return errspec.InvalidParameter("Stringify", elem)
			}
		}
		*s = append(*s, elem)
	}
	return nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func buildDecimalFloatingNumber(num float64, fixed int) float64 {
	rounded := math.Round(num*1000) / 1000
	if fixed > 0 {
		return toFixed(rounded, fixed)
	}
	return rounded
}

func (s *Stringify) buildRendition(rendition *data.Rendition) error {
	attrs := []string{
		fmt.Sprintf("TYPE=%s", rendition.RenditionType),
		fmt.Sprintf("GROUP-ID=%s", rendition.GroupId),
		fmt.Sprintf("NAME=%s", rendition.Name),
	}
	if rendition.IsDefault != nil {
		if aws.BoolValue(rendition.IsDefault) {
			attrs = append(attrs, "DEFAULT=YES")
		} else {
			attrs = append(attrs, "DEFAULT=NO")
		}
	}
	if rendition.AutoSelect != nil {
		if aws.BoolValue(rendition.AutoSelect) {
			attrs = append(attrs, "AUTOSELECT=YES")
		} else {
			attrs = append(attrs, "AUTOSELECT=NO")
		}
	}
	if rendition.Forced != nil {
		if aws.BoolValue(rendition.Forced) {
			attrs = append(attrs, "FORCED=YES")
		} else {
			attrs = append(attrs, "FORCED=NO")
		}
	}
	if rendition.Language != "" {
		attrs = append(attrs, fmt.Sprintf("LANGUAGE=\"%s\"", rendition.Language))
	}
	if rendition.AssocLanguage != "" {
		attrs = append(attrs, fmt.Sprintf("ASSOC-LANGUAGE=\"%s\"", rendition.AssocLanguage))
	}
	if rendition.InStreamId != "" {
		attrs = append(attrs, fmt.Sprintf("INSTREAM-ID=\"%s\"", rendition.InStreamId))
	}
	if rendition.Characteristics != "" {
		attrs = append(attrs, fmt.Sprintf("CHARACTERISTICS=\"%s\"", rendition.InStreamId))
	}
	if rendition.Channels != "" {
		attrs = append(attrs, fmt.Sprintf("CHANNELS=\"%s\"", rendition.Channels))
	}
	if rendition.Uri != "" {
		attrs = append(attrs, fmt.Sprintf("URI=\"%s\"", rendition.Uri))
	}
	return s.Push(fmt.Sprintf("#EXT-X-MEDIA:%s", strings.Join(attrs, ",")))
}

func StringifyMasterPlaylist(playlist *data.MasterPlaylist, skipDuration bool) (string, error) {
	if playlist.DataType != data.TypePlaylist {
		return "", errspec.InvalidParameter("Playlist.DataType", string(playlist.DataType))
	}

	stringify := &Stringify{}
	if err := stringify.Push("#EXTM3U"); err != nil {
		return "", err
	}

	if playlist.Playlist.Version != 0 {
		if err := stringify.Push(fmt.Sprintf("#EXT-X-VERSION:%d", playlist.Playlist.Version)); err != nil {
			return "", err
		}
	}

	if playlist.Playlist.IndependentSegments {
		if err := stringify.Push("#EXT-X-INDEPENDENT-SEGMENTS"); err != nil {
			return "", err
		}
	}

	// TODO: walk around Start property

	for _, sessionData := range playlist.SessionDataList {
		attrs := make([]string, 0)
		attrs = append(attrs, fmt.Sprintf("DATA-ID=\"%s\"", sessionData.Id))
		if sessionData.Language != "" {
			attrs = append(attrs, fmt.Sprintf("LANGUAGE=\"%s\"", sessionData.Language))
		}
		if sessionData.Value != "" {
			attrs = append(attrs, fmt.Sprintf("VALUE=\"%s\"", sessionData.Value))
		} else if sessionData.Uri != "" {
			attrs = append(attrs, fmt.Sprintf("URI=\"%s\"", sessionData.Uri))
		}

		err := stringify.Push("#EXT-X-SESSION-DATA:" + strings.Join(attrs, ","))
		if err != nil {
			return "", err
		}
	}

	for _, sessionKey := range playlist.SessionKeyList {
		attrs := make([]string, 0)
		attrs = append(attrs, "METHOD="+string(sessionKey.Method))
		if sessionKey.Uri != "" {
			attrs = append(attrs, fmt.Sprintf("URI=\"%s\"", sessionKey.Uri))
		}
		if sessionKey.IV != "" {
			if len(sessionKey.IV) != 16 {
				return "", errspec.InvalidParameter("SessionKey.IV", sessionKey.IV)
			}
			attrs = append(attrs, fmt.Sprintf("IV=%s", hex.EncodeToString([]byte(sessionKey.IV))))
		}
		if sessionKey.Format != "" {
			attrs = append(attrs, fmt.Sprintf("KEYFORMAT=\"%s\"", sessionKey.Format))
		}
		if sessionKey.FormatVersion != "" {
			attrs = append(attrs, fmt.Sprintf("KEYFORMATVERSIONS=\"%s\"", sessionKey.FormatVersion))
		}
		err := stringify.Push(fmt.Sprintf("#EXT-X-SESSION-KEY:%s", strings.Join(attrs, ",")))
		if err != nil {
			return "", err
		}
	}

	for _, variant := range playlist.Variants {
		var name string
		if variant.IsIFrameOnly {
			name = "#EXT-X-I-FRAME-STREAM-INF"
		} else {
			name = "#EXT-X-STREAM-INF"
		}
		attrs := make([]string, 0)
		if variant.AverageBandwidth != 0 {
			attrs = append(attrs, fmt.Sprintf("AVERAGE-BANDWIDTH=%d", variant.AverageBandwidth))
		}
		if variant.IsIFrameOnly {
			attrs = append(attrs, fmt.Sprintf("URI=\"%s\"", variant.Uri))
		}
		if variant.Codecs != "" {
			attrs = append(attrs, fmt.Sprintf("CODECS=\"%s\"", variant.Codecs))
		}
		if variant.Resolution != nil {
			attrs = append(attrs, fmt.Sprintf("RESOLUTION=%dx%d", variant.Resolution.Width, variant.Resolution.Height))
		}
		if variant.FrameRate != 0 {
			attrs = append(attrs, fmt.Sprintf("FRAME-RATE=%f", buildDecimalFloatingNumber(variant.FrameRate, 3)))
		}
		if variant.HdcpLevel != 0 {
			attrs = append(attrs, fmt.Sprintf("HDCP-LEVEL=%d", variant.HdcpLevel))
		}
		if len(variant.Audio) > 0 {
			attrs = append(attrs, fmt.Sprintf("AUDIO=\"%s\"", variant.Audio[0].GroupId))
			for _, audio := range variant.Audio {
				err := stringify.buildRendition(audio)
				if err != nil {
					return "", err
				}
			}
		}
		if len(variant.Video) > 0 {
			attrs = append(attrs, fmt.Sprintf("VIDEO=\"%s\"", variant.Video[0].GroupId))
			for _, video := range variant.Video {
				err := stringify.buildRendition(video)
				if err != nil {
					return "", err
				}
			}
		}
		if len(variant.Subtitles) > 0 {
			attrs = append(attrs, fmt.Sprintf("SUBTITLES=\"%s\"", variant.Subtitles[0].GroupId))
			for _, subtitles := range variant.Subtitles {
				err := stringify.buildRendition(subtitles)
				if err != nil {
					return "", err
				}
			}
		}

		if len(variant.ClosedCaptions) > 0 {
			attrs = append(attrs, fmt.Sprintf("CLOSED-CAPTIONS=\"%s\"", variant.ClosedCaptions[0].GroupId))
			for _, closedCaptions := range variant.ClosedCaptions {
				err := stringify.buildRendition(closedCaptions)
				if err != nil {
					return "", err
				}
			}
		}
		err := stringify.Push(fmt.Sprintf("%s:%s", name, strings.Join(attrs, ",")))
		if err != nil {
			return "", err
		}
		if variant.IsIFrameOnly {
			err := stringify.Push(variant.Uri)
			if err != nil {
				return "", err
			}
		}
	}

	return strings.Join(*stringify, "\n"), nil
}
