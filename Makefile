BACKEND_DIR := /root/bitrixChatgpt/

all: backend-first restart-golang

backend-first:
	@echo "Building backend..."
	go build ./

restart-golang:
	@echo "Restart golang..."
	@pm2 restart golang


