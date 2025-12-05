package qdrant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmbedderRepository interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

type Repository struct {
	client   pb.PointsClient
	embedder EmbedderRepository
}

func NewRepository(addr string, embedder EmbedderRepository) *Repository {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("could not connect to Qdrant: %w", err))
	}

	client := pb.NewPointsClient(conn)

	// Validate the connection with a health check
	qdrantClient := pb.NewQdrantClient(conn)
	ctx := context.Background()
	_, err = qdrantClient.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		panic(fmt.Errorf("qdrant health check failed: %w", err))
	}

	return &Repository{
		client:   client,
		embedder: embedder,
	}
}

func (r *Repository) Upsert(ctx context.Context, collection string, content string, metadata map[string]any) error {
	embedding, err := r.embedder.GenerateEmbedding(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	payload := make(map[string]*pb.Value)
	for k, v := range metadata {
		payload[k] = toPbValue(v)
	}
	if _, ok := payload["content"]; !ok {
		payload["content"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: content}}
	}

	id := uuid.New().String()

	_, err = r.client.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: collection,
		Points: []*pb.PointStruct{
			{
				Id: &pb.PointId{
					PointIdOptions: &pb.PointId_Uuid{Uuid: id},
				},
				Vectors: &pb.Vectors{
					VectorsOptions: &pb.Vectors_Vector{
						Vector: &pb.Vector{
							Data: embedding,
						},
					},
				},
				Payload: payload,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upsert to qdrant: %w", err)
	}

	return nil
}

func toPbValue(v any) *pb.Value {
	switch val := v.(type) {
	case string:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: val}}
	case int:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: int64(val)}}
	case int64:
		return &pb.Value{Kind: &pb.Value_IntegerValue{IntegerValue: val}}
	case float64:
		return &pb.Value{Kind: &pb.Value_DoubleValue{DoubleValue: val}}
	case float32:
		return &pb.Value{Kind: &pb.Value_DoubleValue{DoubleValue: float64(val)}}
	case bool:
		return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: val}}
	default:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: fmt.Sprint(val)}}
	}
}
