package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eduardoschulz/memory-sim/internal/pagetable"
	"github.com/eduardoschulz/memory-sim/internal/simulator"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "simulator",
	Short: "Simulador de paginação virtual com MMU, processos leves e algoritmos de substituição",
	Long: `Simulador de memória virtual implementado em Go.

Aplica conceitos de gerenciamento de memória:
  - Memória física: 64 KB (8 frames de 8 KB)
  - Memória virtual: 1 MB (128 páginas de 8 KB)
  - MMU com tabela de páginas e tradução de endereços
  - Processos leves concorrentes via goroutines (produtor/consumidor)
  - Algoritmos de substituição: FIFO (padrão) e LRU`,
	Run: run,
}

func init() {
	rootCmd.Flags().StringP("sizes", "s", "16384,8192",
		"Tamanhos dos processos em bytes (separados por vírgula)")
	rootCmd.Flags().IntP("instructions", "i", 20,
		"Número total de instruções de acesso à memória")
	rootCmd.Flags().StringP("algorithm", "a", "fifo",
		"Opções de Algoritmo: fifo ou lru")
}

func run(cmd *cobra.Command, args []string) {
	sizesStr, _ := cmd.Flags().GetString("sizes")
	instructions, _ := cmd.Flags().GetInt("instructions")
	algorithmStr, _ := cmd.Flags().GetString("algorithm")

	sizes, err := parseSizes(sizesStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao processar --sizes: %v\n", err)
		os.Exit(1)
	}

	var alg pagetable.Algorithm
	switch strings.ToLower(algorithmStr) {
	case "lru":
		alg = pagetable.LRU
	case "fifo":
		alg = pagetable.FIFO
	default:
		fmt.Fprintf(os.Stderr, "Algoritmo inválido: %q. Use 'fifo' ou 'lru'.\n", algorithmStr)
		os.Exit(1)
	}

	if instructions < 1 {
		fmt.Fprintln(os.Stderr, "Número de instruções deve ser >= 1")
		os.Exit(1)
	}

	cfg := simulator.Config{
		Sizes:        sizes,
		Instructions: instructions,
		Algorithm:    alg,
	}

	if err := simulator.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na simulação: %v\n", err)
		os.Exit(1)
	}
}

/*
parseSizes converte uma string "16384,8192,4096" em []int.
*/
func parseSizes(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	sizes := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("valor inválido %q: %w", p, err)
		}
		if v < 1 {
			return nil, fmt.Errorf("tamanho %d inválido (mínimo 1 byte)", v)
		}
		sizes = append(sizes, v)
	}
	return sizes, nil
}
