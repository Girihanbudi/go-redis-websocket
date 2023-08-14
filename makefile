APP_NAME="go-redis-websocket"

# Deploy service dev environments
env:
	cd ./dev && docker-compose -p $(APP_NAME) up -d

# Run the service
server1:
	go run cmd/server/main.go
	
server2:
	go run cmd/server2/main.go