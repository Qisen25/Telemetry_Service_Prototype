package websocket

//pool for handling multiple clients
import (
    "fmt"
    _ "strconv"
)

type Pool struct {
    Register   chan *Client
    Unregister chan *Client
    Clients    map[*Client]bool
    Broadcast  chan Message
    GeoStream  chan GeoPack
    ClientAlert chan Alert
}

func NewPool() *Pool {
    return &Pool{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan Message),
        GeoStream:  make(chan GeoPack),
        ClientAlert:  make(chan Alert),
    }
}   

func (pool *Pool) Start() {
    for {
        select {
        case client := <-pool.Register:
            pool.Clients[client] = true
            client.ConnectStatus = true
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            fmt.Println("User has joined")
            for client, _ := range pool.Clients {
                client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
            }
            break
        case client := <-pool.Unregister:
            removeID := client.Details.ID
            fmt.Println("leaver user Id: " + removeID)
            client.ConnectStatus = false
            
            delete(pool.Clients, client)
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                //client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
                client.Conn.WriteJSON(GeoPack{Type: "RemoveLoc", ID: removeID})
            }
            break
        case geopos := <-pool.GeoStream:
            //fmt.Println("Update geo locations to all clients in Pool")
            for client, _ := range pool.Clients {
                if err := client.Conn.WriteJSON(geopos); err != nil {
                    fmt.Println(err)
                    return
                }
            }
            break
        case clAlert := <-pool.ClientAlert:
            fmt.Println("Alert! " + clAlert.ClientAddr.Details.ID)
            client := clAlert.ClientAddr

            //create struct to write as JSON
            alert := struct {
                Type string `json:"type"`
                Reason string `json:"reason"`
            } {
                clAlert.Type,
                clAlert.Reason,
            }
            if err := client.Conn.WriteJSON(alert); err != nil {
                    fmt.Println(err)
                    return
            }
            break
        case message := <-pool.Broadcast:
            //fmt.Println("Sending message to all clients in Pool")
            for client, _ := range pool.Clients {
                if err := client.Conn.WriteJSON(message); err != nil {
                    fmt.Println(err)
                    return
                }

                //client.Conn.WriteJSON(GeoPack{ID: "1", longitude: "-34.122", latitude: "115.123"});
            }

        }
    }
}