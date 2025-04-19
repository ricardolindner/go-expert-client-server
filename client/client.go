package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Resposta não OK: %d", resp.StatusCode)
	}

	var cotacao Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		log.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	err = salvarCotacao(cotacao.Bid)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação: %v", err)
	}

	fmt.Println("Cotação salva com sucesso!")
}

func salvarCotacao(bid string) error {
	conteudo := fmt.Sprintf("Dólar: %s\n", bid)
	return os.WriteFile("cotacao.txt", []byte(conteudo), os.ModePerm)
}
