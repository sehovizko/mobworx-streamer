package streaming

func GetMimeType(renditionType string) string {
	switch renditionType {
	case "video":
		return "video/mp4"
	case "audio":
		return "audio/mp4"
	default:
		return "application/mp4"

	}

}
