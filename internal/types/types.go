package types

type AccessRequest struct {
	ProcessID      int
	VirtualAddress int
	ResponseCh     chan AccessResponse
}

type AccessResponse struct {
	PhysicalAddress int
	Content         byte
	PageFault       bool
	VirtualPage     int
	PhysicalFrame   int
	Offset          int
	VictimEvicted   int
}
