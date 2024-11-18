FRONTEND_DIR := /root/bitrixChatgpt/frontend/vite-project
BACKEND_DIR := /root/bitrixChatgpt/backend

all: backend-first backend-second frontend-first restart-golang restart-vuejs

backend-first:
	@echo "Cd backend..."
	@cd $(BACKEND_DIR)

backend-second:
	@echo "Building backend..."
	go build ./

frontend-first:
	@echo "Cd frontend..."
	@cd $(FRONTEND_DIR) && npm run build

restart-golang:
	@echo "Restart golang..."
	@pm2 restart golang

restart-vuejs:
	@echo "Restart vuejs..."
	@pm2 restart vuejs

