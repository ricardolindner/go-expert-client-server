package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, data_hora DATETIME)`)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		cotacao, err := buscarCotacao(ctx)
		if err != nil {
			http.Error(w, "Erro ao buscar cotação: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Erro ao buscar cotação: %v", err)
			return
		}

		ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
		defer cancelDB()

		err = salvarCotacao(ctxDB, db, cotacao.Bid)
		if err != nil {
			log.Printf("Erro ao salvar cotação no banco: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cotacao)
	})

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status não OK: %d", resp.StatusCode)
	}

	var dados map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return nil, err
	}

	return &Cotacao{Bid: dados["USDBRL"]["bid"]}, nil
}

func salvarCotacao(ctx context.Context, db *sql.DB, bid string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := db.ExecContext(ctx, `INSERT INTO cotacoes (bid, data_hora) VALUES (?, ?)`, bid, time.Now())
		return err
	}
}
