package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
)

func main() {
	godotenv.Load()

	router := gin.New()

	server := socketio.NewServer(&engineio.Options{})

	_, err := server.Adapter(&socketio.RedisAdapterOptions{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PWD"),
	})
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		var projectId string
		params := strings.Split(s.URL().RawQuery, "&")
		for _, param := range params {
			data := strings.Split(param, "=")
			if len(data) == 2 && data[0] == "project_id" {
				projectId = data[1]
			}
		}

		if projectId == "" {
			s.Close()
			return nil
		}

		log.Println("Joining project with id:", projectId)
		server.JoinRoom("/", projectId, s)
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("Error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Closed:", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socketio listen error: %s\n", err)
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		log.Println("Start broadcasting")
		counter := 1
		for {
			server.BroadcastToRoom("/", "48288ebe-99be-48b8-9c06-755057ee4890", "msg", fmt.Sprintf("send data %d from server 2", counter))
			time.Sleep(1 * time.Second)
			counter++
		}
	}()

	defer server.Close()

	router.GET("/websocket/*any", gin.WrapH(server))
	router.POST("/websocket/*any", gin.WrapH(server))

	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed run app: ", err)
	}
}
