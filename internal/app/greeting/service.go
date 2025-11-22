package greetingapi

import (
	"context"

	"github.com/applinh/mcp-rag-vector/gen/greeting"
	"github.com/applinh/mcp-rag-vector/internal/infra/logger"
)

type greetingService struct {
	logger logger.LoggerInterface
}

// Service implements the greeting service logic.
type Service struct{}

// Greet returns a greeting for the given name.
func (s *greetingService) Greet(ctx context.Context, params *greeting.GreetPayload) (string, error) {
	return "Hello, " + params.Name + "!", nil
}

func NewGreetingService(logger logger.LoggerInterface) greeting.Service {
	return &greetingService{logger: logger}
}
