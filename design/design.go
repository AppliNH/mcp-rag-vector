package design

import (
	. "goa.design/goa/v3/dsl"
)

// API definition
var _ = API("mcp-rag-vector", func() {
	Title("MCP server to serve as a RAG pipeline by allowing an LLM to write and read in a vector DB.")
	Description("A simple API for MCP server to serve as a RAG pipeline by allowing an LLM to write and read in a vector DB.")
	Version("1.0")

	Server("mcp-rag-vector", func() {
		Description("MCP server to serve as a RAG pipeline by allowing an LLM to write and read in a vector DB.")
	})
})

var _ = Service("health", func() {
	Description("Health check service")

	HTTP(func() {
		Path("/checks/health")
	})

	Method("getHealth", func() {
		Description("Returns OK if the service is healthy")
		Payload(func() {
			// No payload
		})
		Result(HealthCheckResponse)
		HTTP(func() {
			GET("/")

			Response(StatusOK)
		})
	})
})

// Greeting service
var _ = Service("greeting", func() {
	Description("Greeting service")

	Method("greet", func() {
		Description("Returns a greeting for the given name")
		Payload(func() {
			Field(1, "name", String, "Name to greet")
			Required("name")
		})
		Result(String)
	})
})

// Knowledge base service
var _ = Service("knowledge-base", func() {
	Description("Knowledge base operations")

	Method("upsert", func() {
		Description("Upserts text into a Qdrant collection")
		Payload(func() {
			Field(1, "collection", String, "Name of the collection")
			Field(2, "content", String, "Text content to embed and store")
			Field(3, "metadata", MapOf(String, Any), "Optional metadata")
			Required("collection", "content")
		})
		Result(Boolean)
	})
})
