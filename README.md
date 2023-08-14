# Go Redis Websocket
 
This repository contain a test to listen to a single server websocket channel and receive message from multiple server using socket io implementation.

### How to Run

1. Download source code
```
git clone https://github.com/Girihanbudi/go-redis-websocket.git
```

2. Install dependencies
```
go mod download
```

3. Install service environment (redis as adapter for server communication)
```
make env
```

4. Run applications
```
make server1
make server2
```