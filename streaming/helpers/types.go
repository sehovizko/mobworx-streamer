package helpers

type SignalAction struct {
	UpdateVariant      string
	UpdateRendition    string
	UpdateSegment      string
	UpdatePart         string
	UpdateDemuxSegment string
	UpdateDemuxPart    string
	Ping               string
	Abort              string
	Terminate          string
	AckVariant         string
	AckRendition       string
	AckSegment         string
	AckPart            string
	AckDemuxSegment    string
	AckDemuxPart       string
	Pong               string
	Aborted            string
	Terminated         string
}

func GetSignalAction() *SignalAction {
	return &SignalAction{
		UpdateVariant:      "updateVariant",
		UpdateRendition:    "updateRendition",
		UpdateSegment:      "updateSegment",
		UpdatePart:         "updatePart",
		UpdateDemuxSegment: "updateDemuxSegment",
		UpdateDemuxPart:    "updateDemuxPart",
		Ping:               "ping",
		Abort:              "abort",
		Terminate:          "terminate",
		AckVariant:         "ackVariant",
		AckRendition:       "ackRendition",
		AckSegment:         "ackSegment",
		AckPart:            "ackPart",
		AckDemuxSegment:    "ackDemuxSegment",
		AckDemuxPart:       "ackDemuxPart",
		Pong:               "pong",
		Aborted:            "aborted",
		Terminated:         "terminated",
	}
}
