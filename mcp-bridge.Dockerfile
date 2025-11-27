FROM python:3.12-slim

WORKDIR /app

RUN apt update 
RUN apt install -y curl

COPY mcp-bridge.requirements.txt .
RUN pip install --no-cache-dir -r mcp-bridge.requirements.txt
