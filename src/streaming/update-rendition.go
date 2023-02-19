package streaming

type RenditionType string

const (
	VideoRenditionType RenditionType = "video"
	AudioRenditionType RenditionType = "audio"
)

type MimeType string

const (
	Video       MimeType = "video/mp4"
	Audio       MimeType = "audio/mp4"
	Application MimeType = "application/mp4"
)

var renditionToMimeType = map[RenditionType]MimeType{
	VideoRenditionType: Video,
	AudioRenditionType: Audio,
}

func GetMimeType(renditionType RenditionType) string {
	if value, ok := renditionToMimeType[renditionType]; ok {
		return string(value)
	}
	return string(Application)
}
