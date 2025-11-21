package design

import (
	. "goa.design/goa/v3/dsl"
)

var HealthCheckResponse = Type("HealthCheckResponse", func() {
	Description("Response containing health check status")

	Attribute("healthy", Boolean, "health or not", func() {
		Example(true)
	})

	Attribute("message", String, "health message", func() {
		Example("MCP RAG Vector API is healthy")
	})

	Required("healthy")
})
