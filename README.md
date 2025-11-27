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

## Test MCP via HTTP with SSE

1. In one terminal, start the server:
```bash
curl -N -H "Accept: text/event-stream" http://localhost:3000/mcp/sse
```

2. Copy the `sessionId` from the response:

```
event: endpoint
data: /mcp/message?sessionId=XXX
```

3. In another terminal, invoke the greeting tool:
```bash
curl -i -X POST "http://localhost:3000/mcp/message?sessionId=<SESSION_ID>" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "0",
    
      "method": "tools/call",
      "params": {
        "name": "greet",
        "arguments": { "name": "Alice" }
      }
    
  }'
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
      { "role": "user", "content": "Use the greeting tool to greet me. My name is Thomas." }
    ],
    "stream": false
  }'
```