package websocket

import (
    "fmt"
    "log"
    _ "sync"
    "time"
    "encoding/json"

    "github.com/go-redis/redis"
    "github.com/gorilla/websocket"
)

const (
    // Time allowed to write a message to the peer.
    writeWait = 500 * time.Millisecond

    //notification wait
    riskWait = 1 * time.Minute

    // Time allowed to read the next pong message from the peer.
    pongWait = 1* time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

type Client struct {
    // ID   string
    Conn *websocket.Conn
    Pool *Pool
    ConnectStatus bool
    // Longitude float64   `json:"longitude"`
    // Latitude float64    `json:"latitude"`
    Details ClientDetails
    AlertsMade map[string]Alert
}

type ClientDetails struct {
    ID   string      `json:"uid"`
    Longitude float64   `json:"longitude"`
    Latitude float64    `json:"latitude"`
}

type Message struct {
    Type int    `json:"type"`
    Body string `json:"body"`
}

type Alert struct {
    Type string    `json:"type"`
    Reason string `json:"reason"`
    ClientAddr *Client `json:"client"`
    VulnerableID string `json:"vulnerable"`
    waitRisk time.Time `json:"waitrisk"`
}

type GeoPack struct {
    Type string `json:"type"`
    ID string    `json:"ID"`
    Longitude float64   `json:"longitude"`
    Latitude float64    `json:"latitude"`
}

//read messages from client app
func (c *Client) Read() {
    defer func() {
        c.Pool.Unregister <- c
        c.Conn.Close()
    }()

    
    for {
        messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        //convert json string to struct "getting client location
        var info ClientDetails
        in := []byte(p)
        err = json.Unmarshal(in, &info)
        if err != nil {
            panic(err)
        }

        //assign struct
        c.Details = info

        message := Message{Type: messageType, Body: string(p)}
        c.Pool.Broadcast <- message
        //fmt.Printf("Message Received: %+v\n", message)
        //fmt.Println("Message Received: " + info.ID)
        clientGeo := GeoPack{Type: "Location", ID: c.Details.ID, Longitude: c.Details.Longitude, Latitude: c.Details.Latitude}
        c.Pool.GeoStream <- clientGeo

    }
}

func (c *Client) StreamVulnerableUsers(redisClient *redis.Client) {
    ticker := time.NewTicker(pingPeriod)
    notifyTicker := time.NewTicker(riskWait)

    defer func() {
        ticker.Stop()
        notifyTicker.Stop()
        c.Conn.Close()
    }()

    for {
        //fmt.Println("Sending")
        select{
        case <- ticker.C:
            //c.Conn.SetWriteDeadline(time.Now().Add(1000))
            //fmt.Println("write")

            //check if client disconnects then stop them from streaming data
            if c.ConnectStatus == false {
                log.Println(c.Details.ID + " Stop streaming vulnerables")
                return
            }

            locations, err := redisClient.GeoRadius("publishers", c.Details.Longitude, c.Details.Latitude,  &redis.GeoRadiusQuery{
                Radius:      99999999999999,
                Unit:        "m",
                //WithGeoHash: true,
                WithCoord:   true,
                WithDist:    true,
            }).Result()
            if err != nil {
                log.Fatal(err)
            }
            //fmt.Println(c.Details.Longitude)
            for _, coord := range locations{
                geo := GeoPack{Type: "Location", ID: coord.Name, Longitude: coord.Longitude, Latitude: coord.Latitude}
                c.Pool.GeoStream <- geo

                //fmt.Println(coord.Dist)
                if coord.Dist <= 3.0  {//check if cyclist nearby
                    fmt.Println(coord.Dist)

                    newAlert := Alert{Type: "Alert", Reason: "Careful now, cyclist very close", 
                                      ClientAddr: c, VulnerableID: coord.Name, waitRisk: time.Now()}
                    userAlert, found := c.AlertsMade[coord.Name]

                    fmt.Println(time.Since(userAlert.waitRisk))

                    if !found {
                        c.AlertsMade[coord.Name] = newAlert
                        c.Pool.ClientAlert <- c.AlertsMade[coord.Name]  //broadcast alert to specific client
                        //c.Conn.WriteJSON(newAlert)//send to client side
                        fmt.Println(c.AlertsMade[coord.Name])
                    } else {
                        if time.Since(userAlert.waitRisk).Round(time.Second) >= riskWait{
                            fmt.Println(userAlert)  
                            //c.Conn.WriteJSON(userAlert)//send to client side
                            c.Pool.ClientAlert <- c.AlertsMade[coord.Name] 
                            c.AlertsMade[coord.Name] = newAlert
                        }  
                    }   
                }
            } 
        }
    }      
}

func (c *Client) RemoveOfflinePos(redisClient *redis.Client) {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {

        // if c.ConnectStatus == false {
        //         return
        // }

        select{
        case <- ticker.C:
            //c.Conn.SetWriteDeadline(time.Now().Add(1000))
            //fmt.Println("rm")

            offlinePos, err := redisClient.SPop("offline").Result()

            if err == nil {
                removePkg := GeoPack{Type: "RemoveLoc", ID: offlinePos}
                c.Pool.GeoStream <- removePkg //broadcast position status to all clients in pool
                fmt.Println("cleaning up offline vulnerables")  
            } 
        }
    }     
}

// func (c *Client) AlertClient(redisClient *redis.Client) {
//     ticker := time.NewTicker(pingPeriod)
//     defer func() {
//         ticker.Stop()
//         c.Conn.Close()
//     }()

//     for {
//         select{
//         case <- ticker.C:
//             //c.Conn.SetWriteDeadline(time.Now().Add(1000))
//             //fmt.Println("rm")

//             offlinePos, err := redisClient.SPop("offline").Result()

//             if err == nil {
//                 removePkg := GeoPack{Type: "RemoveLoc", ID: offlinePos}
//                 c.Pool.GeoStream <- removePkg //broadcast position status to all clients in pool
//                 fmt.Println("cleaning up offline users")  
//             } 
//         }
//     }     
// }