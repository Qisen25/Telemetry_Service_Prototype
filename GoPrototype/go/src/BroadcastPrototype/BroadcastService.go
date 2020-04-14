package main

import (
	"log"
	"net/http"
	"fmt"

    "github.com/go-redis/redis"
	"myWebsocket"
)

var redisCl *redis.Client

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
    fmt.Println("WebSocket Endpoint Hit")
    conn, err := websocket.Upgrade(w, r)
    if err != nil {
        fmt.Fprintf(w, "%+v\n", err)
    }

    client := &websocket.Client{
        Conn: conn,
        Pool: pool,
        AlertsMade:    make(map[string]websocket.Alert),
    }

    pool.Register <- client
    
    //run tasks concurrently
    go client.Read()
    go client.StreamVulnerableUsers(redisCl)//better to move this to pool
    go client.RemoveOfflinePos(redisCl)//Better to move this to the pool
}

func setupRoutes() {
    pool := websocket.NewPool()
    go pool.Start()

    //listen for requests
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(pool, w, r)//answer request
    })
}

func CreateRedisClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })

    pong, err := client.Ping().Result()
    fmt.Println(pong, err)

    return client
}

func main() {

	fmt.Println("Program has started?")

	//establish redis connection
    redisCl = CreateRedisClient()

	setupRoutes()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Listening to PORT 31415")
	err := http.ListenAndServeTLS(":31415", "https-server.crt", "https-server.key", nil)
    if err != nil {
        log.Fatal(err)
    }
}