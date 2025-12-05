package knowledgebaseapi

import (
	"context"
	"fmt"

	knowledgebase "github.com/applinh/mcp-rag-vector/gen/knowledge_base"
	"github.com/applinh/mcp-rag-vector/internal/infra/logger"
)

type VectorRepository interface {
	Upsert(ctx context.Context, collection string, content string, metadata map[string]any) error
}

type Service struct {
	logger logger.LoggerInterface
	repo   VectorRepository
}

func NewService(logger logger.LoggerInterface, repo VectorRepository) Service {
	return Service{logger: logger, repo: repo}
}

func (s *Service) Upsert(ctx context.Context, p *knowledgebase.UpsertPayload) (bool, error) {
	s.logger.Info(ctx, fmt.Sprintf("Upserting to collection %s", p.Collection))
	err := s.repo.Upsert(ctx, p.Collection, p.Content, p.Metadata)
	if err != nil {
		s.logger.Error(ctx, fmt.Sprintf("Failed to upsert: %v", err))
		return false, err
	}
	return true, nil
}
