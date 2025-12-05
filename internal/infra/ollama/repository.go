package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

type Repository struct {
	client *api.Client
	model  string
}

func NewRepository(host, model string) *Repository {
	parsedURL, err := url.Parse(host)
	if err != nil {
		// Fallback to default if parsing fails
		panic(fmt.Errorf("could not parse Ollama host: %w", err))
	}

	client := api.NewClient(parsedURL, http.DefaultClient)
	if err := client.Heartbeat(context.Background()); err != nil {
		panic(fmt.Errorf("could not heartbeat Ollama: %w", err))
	}

	return &Repository{
		client: client,
		model:  model,
	}
}

func (r *Repository) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	req := &api.EmbedRequest{
		Model: r.model,
		Input: text,
	}

	resp, err := r.client.Embed(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not generate embedding: %w", err)
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Convert []float64 to []float32
	embedding := make([]float32, len(resp.Embeddings[0]))
	for i, v := range resp.Embeddings[0] {
		embedding[i] = float32(v)
	}

	return embedding, nil
}
