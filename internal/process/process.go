package process

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/eduardoschulz/memory-sim/internal/constants"
	"github.com/eduardoschulz/memory-sim/internal/types"
)

type LightProcess struct {
	ID           int
	size         int
	virtualStart int
}

func NewLightProcess(id, size, virtualStart int) *LightProcess {
	return &LightProcess{
		ID:           id,
		size:         size,
		virtualStart: virtualStart,
	}
}

/*
Execute executa um processo leve, gerando acessos aleatórios ao end. virtual e enviando requisições para a MMU via go channel.
*/
func (p *LightProcess) Execute(ctx context.Context, wg *sync.WaitGroup, reqCh chan<- types.AccessRequest) {
	defer wg.Done()

	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(p.ID)))

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		virtualAddr := p.virtualStart + rng.Intn(p.size)
		vPage := virtualAddr / constants.PageSize
		offset := virtualAddr % constants.PageSize

		/*
			   channels é mecanismo de comunicação entre goroutines; uma goroutines envia ch <-, e a outra recebe <-ch e bloqueia até o
				 dados chegar, sem precisar usar mutex
		*/
		respCh := make(chan types.AccessResponse, 1)
		req := types.AccessRequest{
			ProcessID:      p.ID,
			VirtualAddress: virtualAddr,
			ResponseCh:     respCh,
		}

		fmt.Printf("[Processo %d] acessando o end. virtual %d → página %d, offset %d\n",
			p.ID, virtualAddr, vPage, offset)

		// envia a requisição para a MMU via reqCh ou desiste se o contexto foi cancelado.
		select {
		case reqCh <- req:
		case <-ctx.Done():
			return
		}

		select {
		case <-respCh:
		case <-ctx.Done():
			return
		}

		time.Sleep(80 * time.Millisecond)
	}
}
