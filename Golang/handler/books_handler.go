package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type VolumeInfoSimplified struct {
	Title         string `json:"title"`
	Publisher     string `json:"publisher"`
	PublishedDate string `json:"publishedDate"`
}

type SimplifiedBookResult struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Publisher     string `json:"publisher,omitempty"`
	PublishedDate string `json:"publishedDate,omitempty"`
}

type GoogleBooksResponse struct {
	Items []struct {
		ID string `json:"id"`
		// CORRIGIDO (Linha 20): Verifique a sintaxe da tag. Deve ser entre crases.
		VolumeInfo VolumeInfoSimplified `json:"volumeInfo"`
	} `json:"items"`
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "O parâmetro 'q' (busca) é obrigatório.", http.StatusBadRequest)
		return
	}

	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=5", encodedQuery)

	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Erro ao conectar com a API do Google Books.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API do Google Books retornou status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta da API.", http.StatusInternalServerError)
		return
	}

	var booksResponse GoogleBooksResponse
	if err := json.Unmarshal(body, &booksResponse); err != nil {
		http.Error(w, "Erro ao decodificar a resposta JSON da API.", http.StatusInternalServerError)
		return
	}

	var results []SimplifiedBookResult
	for _, item := range booksResponse.Items {
		results = append(results, SimplifiedBookResult{
			ID:            item.ID,
			Title:         item.VolumeInfo.Title,
			Publisher:     item.VolumeInfo.Publisher,
			PublishedDate: item.VolumeInfo.PublishedDate,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Erro ao encodar a resposta para o cliente.", http.StatusInternalServerError)
		return
	}
}
