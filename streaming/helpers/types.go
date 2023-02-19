package helpers

type SignalAction string

const (
	UpdateVariant      SignalAction = "updateVariant"
	UpdateRendition    SignalAction = "updateRendition"
	UpdateSegment      SignalAction = "updateSegment"
	UpdatePart         SignalAction = "updatePart"
	UpdateDemuxSegment SignalAction = "updateDemuxSegment"
	UpdateDemuxPart    SignalAction = "updateDemuxPart"
	Ping               SignalAction = "ping"
	Abort              SignalAction = "abort"
	Terminate          SignalAction = "terminate"
	AckVariant         SignalAction = "ackVariant"
	AckRendition       SignalAction = "ackRendition"
	AckSegment         SignalAction = "ackSegment"
	AckPart            SignalAction = "ackPart"
	AckDemuxSegment    SignalAction = "ackDemuxSegment"
	AckDemuxPart       SignalAction = "ackDemuxPart"
	Pong               SignalAction = "pong"
	Aborted            SignalAction = "aborted"
	Terminated         SignalAction = "terminated"
)
