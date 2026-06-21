// Package memory pacote que tem a estrutura da memória
package memory

import "github.com/eduardoschulz/memory-sim/internal/constants"

type PhysicalMemory struct {
	frames     [constants.NumFrames][]byte
	freeFrames []int
}

func NovaPhysicalMemory() *PhysicalMemory {
	pm := &PhysicalMemory{}
	for i := 0; i < constants.NumFrames; i++ {
		pm.frames[i] = make([]byte, constants.PageSize)
		pm.freeFrames = append(pm.freeFrames, i)
	}
	return pm
}

func (pm *PhysicalMemory) AllocateFrame() int {
	if len(pm.freeFrames) == 0 {
		return -1
	}
	frame := pm.freeFrames[0]
	pm.freeFrames = pm.freeFrames[1:]
	return frame
}

func (pm *PhysicalMemory) FreeFrame(frame int) {
	pm.freeFrames = append(pm.freeFrames, frame)
}

func (pm *PhysicalMemory) WriteToFrame(frame int, data []byte){
	copy(pm.frames[frame], data)
}

func (pm *PhysicalMemory) ReadByte(frame int, offset int) byte {
	return pm.frames[frame][offset]
}

func (pm *PhysicalMemory) GetFreeFrames() int {
	return len(pm.freeFrames)
}

