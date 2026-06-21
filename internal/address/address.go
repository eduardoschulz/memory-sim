package address


type VirtualAddress struct {
	Value uint32
}

type PhysicalAddress struct {
	Frame uint8
	Offset uint16
}

func NewVirtualAddress(r uint32) *VirtualAddress {
	return &VirtualAddress{
		Value: r,
	}
}

func NewPhysicalAddress(frame uint8, offset uint16) *PhysicalAddress {
	return &PhysicalAddress{
		Frame: frame,
		Offset: offset,
	}
}

