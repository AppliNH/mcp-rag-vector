package design

import (
	. "goa.design/goa/v3/dsl"
)

// API definition
var _ = API("github.com/applinh/mcp-rag-vector", func() {
	Title("MCP server to serve as a RAG pipeline by allowing an LLM to write and read in a vector DB. API")
	Description("A simple API for MCP server to serve as a RAG pipeline by allowing an LLM to write and read in a vector DB.")
	Version("1.0")

	Server("github.com/applinh/mcp-rag-vector", func() {
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
