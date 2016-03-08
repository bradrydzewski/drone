package stream

type PacketType string

const (
	PacketLine   = "L"
	PacketHeader = "H"
	PacketFooter = "F"
)

type Packet struct {
	Type PacketType `json:"type,omitempty"`
	Step int        `json:"step,omitempty"`
	Line int        `json:"line,omitempty"`
	Time int        `json:"time,omitempty"`
	Data string     `json:"data,omitempty"`
}
