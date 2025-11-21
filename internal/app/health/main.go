package healthapi

import (
	"context"

	"github.com/applinh/mcp-rag-vector/gen/health"
	"github.com/applinh/mcp-rag-vector/internal/infra/logger"
)

type healthService struct {
	logger logger.LoggerInterface
}

var healthMessage string = "MCP RAG Vector API is healthy"

func (g *healthService) GetHealth(ctx context.Context) (res *health.HealthCheckResponse, err error) {
	return &health.HealthCheckResponse{
		Healthy: true,
		Message: &healthMessage,
	}, nil
}

func NewHealthService(logger logger.LoggerInterface) health.Service {
	return &healthService{logger: logger}
}
