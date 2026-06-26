# Memory Sim

Simulador de gerenciamento de memória virtual paginada para a disciplina de Sistemas Operacionais II.

```bash
# clonando o repositório
git clone https://github.com/eduardoschulz/memory-sim
cd memory-sim

# baixando as dependências
go mod download

# compilando o programa
go build -o simulator cmd/simulator/main.go

# executando o programa
./simulator -s 16384,8192 -i 20 -a fifo
```

### Flags

| Flag | Descrição | Padrão |
|------|-----------|--------|
| `-s`, `--sizes` | Tamanhos dos processos em bytes (separados por vírgula) | `16384,8192` |
| `-i`, `--instructions` | Número total de instruções de acesso à memória | `20` |
| `-a`, `--algorithm` | Algoritmo de substituição de páginas (`fifo` ou `lru`) | `fifo` |
