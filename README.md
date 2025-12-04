# mcp-rag-vector

## Configure MCP with Claude Desktop

1. Generate code install in GOBIN:
```bash
make generate-code-design
make install-bin
```

2. Configure the following in Claude Desktop settings:
```json
{
  "mcpServers": {
    "greetingmcp": {
      "command": "$GOPATH/bin/mcp-rag-vector",
      "args": ["mcp"]
    }
  }
}
```

## Local usage with ollama

1. Install ollama and run `ollama serve`

2. Start the stack

```bash
docker compose up -d --build
```

3. Wait for the model to be downloaded (2 GB)

4. Query the API

```bash
curl -X POST http://localhost:8000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama3.2:3b",
    "messages": [
      { "role": "user", "content": "Use the greet tool with my name thomas, return what it says" }
    ],
    "stream": false
  }'
```

## Query the MCP via cURL

1. Start the API

`docker-commpose up` or `make run-server`

2. Use `make http-call-mcp`