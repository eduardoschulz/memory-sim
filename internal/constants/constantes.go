// constants é só um pacote para conter as constantes
package constants

import "fmt"

const (
	PageSize    = 8192
	PhysicalMem = 65536
	VirtualMem  = 1048576
	NumFrames   = PhysicalMem / PageSize // 8 frames
	NumPages    = VirtualMem / PageSize
)

// ValidateProcessSize valida se o tamanho do processo tem entre 1 byte e um 1 MB
func ValidateProcessSize(sizes []int) error {
	if len(sizes) < 2 {
		return fmt.Errorf("é necessário no mínimo 2 processos (fornecidos: %d)", len(sizes))
	}

	soma := 0
	for i, t := range sizes {
		if t < 1 || t > VirtualMem {
			return fmt.Errorf(
				"processo %d: tamanho %d inválido (deve estar entre 1 e %d bytes)",
				i+1, t, VirtualMem,
			)
		}
		soma += t
	}

	if soma > VirtualMem {
		return fmt.Errorf(
			"soma dos sizes (%d bytes) excede a memória virtual de %d bytes",
			soma, VirtualMem,
		)
	}

	return nil
}
