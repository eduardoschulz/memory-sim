package pagetable

import (
	"time"
)

type Algorithm int

/*
FIFO First in First out
LRU Least Recently Used
*/
const (
	FIFO Algorithm = iota
	LRU
)

type PageEntry struct {
	PhysicalFrame int
	Present       bool
	FIFOOrder     int64
	LastAccess    int64
}

type PageTable struct {
	Entries     map[int]*PageEntry
	Algorithm   Algorithm
	FIFOOrder   []int
	FIFOCounter int64
}

/*
NewPageTable cria uma uma tabela de páginas vazias usando o algoritmo passado por argumento
*/
func NewPageTable(alg Algorithm) *PageTable {
	return &PageTable{
		Entries:   make(map[int]*PageEntry),
		Algorithm: alg,
	}
}

/*
Translate retorna a entry da tabela para uma página virtual
se a pág. ainda não foi mapeada, retorna entry com presente = false
*/

func (pt *PageTable) Translate(vPg int) *PageEntry {
	entry, exist := pt.Entries[vPg]
	if exist {
		return entry
	}
	return &PageEntry{PhysicalFrame: -1, Present: false}
}

/*
Update insere ou atualiza os mapas da página no frame da tabela
para o FIFOm registra a ordem de insercao na fila
*/
func (pt *PageTable) Update(vPg int, frame int) {
	entry, existe := pt.Entries[vPg]
	if !existe {
		entry = &PageEntry{}
		pt.Entries[vPg] = entry
	}
	entry.PhysicalFrame = frame
	entry.Present = true
	entry.LastAccess = time.Now().UnixNano()

	if pt.Algorithm == FIFO {
		entry.FIFOOrder = pt.FIFOCounter
		pt.FIFOCounter++
		pt.FIFOOrder = append(pt.FIFOOrder, vPg)
	}
}

/*
Remove tira o mapeamento de uma pág. virtual da tabela e marca como não presente
*/
func (pt *PageTable) Remove(vPg int) {
	entry, exists := pt.Entries[vPg]
	if !exists {
		return
	}
	entry.Present = false
	entry.PhysicalFrame = -1
}

/*
RecordAccess atualiza o timestamp do ultimo acesso da página, utiliza o LRY para rastrear a ordem de uso.
*/
func (pt *PageTable) RecordAccess(vPg int) {
	if entry, exists := pt.Entries[vPg]; exists {
		entry.LastAccess = time.Now().UnixNano()
	}
}

/*
FindVictim seleciona uma página a ser removida da mem. física usand o algoritmo
configurado, retorna -1 se não existir nenhuma página
*/
func (pt *PageTable) FindVictim() int {
	switch pt.Algorithm {
	case LRU:
		return pt.findVictimLRU()
	default:
		return pt.findVictimFIFO()
	}
}

/*
findVictimFIFO remove e retorna a pag. mais antiga da fila de insercao
*/
func (pt *PageTable) findVictimFIFO() int {
	if len(pt.FIFOOrder) == 0 {
		return -1
	}
	victim := pt.FIFOOrder[0]
	pt.FIFOOrder = pt.FIFOOrder[1:]
	return victim
}

/*
TotalPresentEntries retorna quantas pag. estao atualmente mapeadas na mem. fisica (present = true)
*/
func (pt *PageTable) TotalPresentEntries() int {
	count := 0
	for _, e := range pt.Entries {
		if e.Present {
			count++
		}
	}
	return count
}

/*
findVictimLRU percorre todas as entradas presentes e retorna a pag. com o timestamp de ultime acesso mais antigo
*/
func (pt *PageTable) findVictimLRU() int {
	var victim int = -1
	var oldest int64 = 1<<63 - 1
	for vPg, entry := range pt.Entries {
		if entry.Present && entry.LastAccess < oldest {
			oldest = entry.LastAccess
			victim = vPg
		}
	}
	return victim
}
