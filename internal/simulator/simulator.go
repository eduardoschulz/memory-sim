package simulator

import (
	"context"
	"fmt"
	"sync"

	"github.com/eduardoschulz/memory-sim/internal/constants"
	"github.com/eduardoschulz/memory-sim/internal/memory"
	"github.com/eduardoschulz/memory-sim/internal/mmu"
	"github.com/eduardoschulz/memory-sim/internal/pagetable"
	"github.com/eduardoschulz/memory-sim/internal/process"
	"github.com/eduardoschulz/memory-sim/internal/types"
)

/*
Config agrupa os parâmetros de entrada da simulação.
*/
type Config struct {
	Sizes        []int
	Instructions int
	Algorithm    pagetable.Algorithm
}

/*
Run orquestra toda a simulação:
  - Cria e preenche o espaço virtual de 1 MB
  - Instancia a memória física (64 KB, 8 frames)
  - Cria a tabela de páginas e a MMU
  - Lança os processos leves como goroutines produtoras
  - Inicia a MMU como goroutine consumidora
  - Aguarda o número configurado de instruções e encerra
*/
func Run(cfg Config) error {
	if err := constants.ValidateProcessSize(cfg.Sizes); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	fmt.Println("=== Simulador de Memória Virtual ===")
	fmt.Printf("Memória física: %d KB (%d frames de %d KB)\n",
		constants.PhysicalMem/1024,
		constants.NumFrames,
		constants.PageSize/1024,
	)
	fmt.Printf("Memória virtual: %d MB (%d páginas de %d KB)\n",
		constants.VirtualMem/(1024*1024),
		constants.NumPages,
		constants.PageSize/1024,
	)
	fmt.Printf("Processos: %d | Instruções: %d | Algoritmo: ",
		len(cfg.Sizes), cfg.Instructions)
	if cfg.Algorithm == pagetable.LRU {
		fmt.Println("LRU")
	} else {
		fmt.Println("FIFO")
	}
	fmt.Println("=====================================")
	fmt.Println()

	// Inicialização do espaço virtual (backing store de 1 MB).
	virtualMem := make([]byte, constants.VirtualMem)
	globalOffset := 0
	processes := make([]*process.LightProcess, len(cfg.Sizes))

	for i, size := range cfg.Sizes {
		proc := process.NewLightProcess(i+1, size, globalOffset)
		processes[i] = proc

		for j := 0; j < size; j++ {
			virtualMem[globalOffset+j] = byte(((i+1)*73 + j) % 256)
		}

		fmt.Printf("[Init] Processo %d: %d bytes (endereços virtuais %d–%d), %d páginas\n",
			proc.ID, size, globalOffset, globalOffset+size-1,
			(size+constants.PageSize-1)/constants.PageSize,
		)
		globalOffset += size
	}
	fmt.Println()

	// Componentes da MMU.
	physMem := memory.NovaPhysicalMemory()
	table := pagetable.NewPageTable(cfg.Algorithm)
	mmuInstance := mmu.NewMMU(physMem, virtualMem, table)

	// Canal produtor → consumidor.
	reqCh := make(chan types.AccessRequest)

	// Contexto para cancelamento.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Goroutines produtoras (processos leves).
	for _, proc := range processes {
		wg.Add(1)
		go proc.Execute(ctx, &wg, reqCh)
	}

	// Goroutine consumidora (MMU):
	// lê requisições do canal, processa e devolve a resposta.
	// Cancela o contexto ao atingir o número de instruções.
	remaining := cfg.Instructions
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-reqCh:
				resp := mmuInstance.ProcessRequest(req)
				req.ResponseCh <- resp
				remaining--
				if remaining <= 0 {
					cancel()
					close(done)
					return
				}
			}
		}
	}()

	// Aguarda o término (N instruções processadas).
	<-done

	// Drena requisições pendentes e aguarda processos encerrarem.
	wg.Wait()
	close(reqCh)

	// Relatório final.
	fmt.Println("=====================================")
	fmt.Println("Simulação encerrada.")
	fmt.Printf("Instruções processadas: %d\n", cfg.Instructions)
	fmt.Printf("Páginas mapeadas na tabela: %d\n", table.TotalPresentEntries())
	fmt.Println("=====================================")

	return nil
}
