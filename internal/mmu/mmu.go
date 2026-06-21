package mmu

import (
	"fmt"

	"github.com/eduardoschulz/memory-sim/internal/constants"
	"github.com/eduardoschulz/memory-sim/internal/memory"
	"github.com/eduardoschulz/memory-sim/internal/pagetable"
	"github.com/eduardoschulz/memory-sim/internal/types"
)

type MMU struct {
	PhysicalMemory *memory.PhysicalMemory
	VirtualMemory  []byte
	PageTable      *pagetable.PageTable
}

func NewMMU(pm *memory.PhysicalMemory, vMem []byte, table *pagetable.PageTable) *MMU {
	return &MMU{
		PhysicalMemory: pm,
		VirtualMemory:  vMem,
		PageTable:      table,
	}
}

/*
ProcessRequest é o metodo principal de acesso ao MMU
*/
func (m *MMU) ProcessRequest(req types.AccessRequest) {
	vPage := req.VirtualAddress / constants.PageSize
	offset := req.VirtualAddress % constants.PageSize

	entry := m.PageTable.Translate(vPage)

	if entry.Present {
		req.ResponseCh <- m.handleHit(req, vPage, offset, entry)
	} else {
		req.ResponseCh <- m.handlePageFault(req, vPage, offset)
	}
}

/*
handleHit processa um acesso cuja a pag. já esta carregada na memoria fisica.
*/
func (m *MMU) handleHit(_ types.AccessRequest, vPage, offset int, entry *pagetable.PageEntry) types.AccessResponse {
	m.PageTable.RecordAccess(vPage)
	physAddr := entry.PhysicalFrame*constants.PageSize + offset
	content := m.PhysicalMemory.ReadByte(entry.PhysicalFrame, offset)

	fmt.Printf("[MMU]   Tradução: Página virtual %d → Frame %d → End físico %d\n",
		vPage, entry.PhysicalFrame, physAddr)
	fmt.Printf("[MMU]   HIT! Conteúdo no endereço %d: 0x%02X\n", physAddr, content)
	fmt.Println()

	return types.AccessResponse{
		PhysicalAddress: physAddr,
		Content:         content,
		PageFault:       false,
		VirtualPage:     vPage,
		PhysicalFrame:   entry.PhysicalFrame,
		Offset:          offset,
		VictimEvicted:   -1,
	}
}

/*
handlePageFault da com page fault: tenta alocar frame livre; se não houver, aplica o algoritmo de substituição (FIFO/LRU) conforme configurado.
*/
func (m *MMU) handlePageFault(_ types.AccessRequest, vPage, offset int) types.AccessResponse {
	fmt.Printf("[MMU]   Tradução: Página virtual %d → Falta de página!\n", vPage)

	frame := m.PhysicalMemory.AllocateFrame()
	victim := -1

	if frame == -1 {
		victim = m.PageTable.FindVictim()
		victimEntry := m.PageTable.Translate(victim)
		frame = victimEntry.PhysicalFrame
		m.PageTable.Remove(victim)

		fmt.Printf("[MMU]   Sem frames livres! Substituindo página %d (Frame %d).\n",
			victim, frame)
	} else {
		fmt.Printf("[MMU]   Frame %d livre. Carregando página %d no frame %d.\n",
			frame, vPage, frame)
	}

	m.loadPage(vPage, frame)
	m.PageTable.Update(vPage, frame)

	fmt.Printf("[MMU]   Tabela atualizada: Página virtual %d → Frame físico %d\n",
		vPage, frame)

	physAddr := frame*constants.PageSize + offset
	content := m.PhysicalMemory.ReadByte(frame, offset)

	fmt.Printf("[MMU]   End físico %d, conteúdo: 0x%02X\n", physAddr, content)
	fmt.Println()

	return types.AccessResponse{
		PhysicalAddress: physAddr,
		Content:         content,
		PageFault:       true,
		VirtualPage:     vPage,
		PhysicalFrame:   frame,
		Offset:          offset,
		VictimEvicted:   victim,
	}
}

/*
loadPage copia os 8kb da pág. virtual para o frame físico indicado.
*/
func (m *MMU) loadPage(vPage int, frame int) {
	start := vPage * constants.PageSize
	end := start + constants.PageSize
	data := m.VirtualMemory[start:end]
	m.PhysicalMemory.WriteToFrame(frame, data)
}
