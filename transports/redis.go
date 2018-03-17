package transports

import (
    "fmt"
    "time"

    log "github.com/sirupsen/logrus"
    "github.com/go-redis/redis"

    helpers "github.com/georgepsarakis/go-stats/helpers"
)

type Redis struct {
    client *redis.Client
}

func (transport *Redis) Connect(address string, database int) {
    transport.client = redis.NewClient(&redis.Options{Addr: address, DB: database})
    reply, connection_error := transport.client.Ping().Result()
    helpers.PanicOnError(connection_error)
    log.Debug(fmt.Sprintf("Redis Server PING -> %s", reply))
}

func (transport *Redis) Listen(channel_name string, go_channel chan string) {
    pubsub_client := transport.client.Subscribe(channel_name)
    defer pubsub_client.Close()
    subscription, pubsub_error :=
        pubsub_client.ReceiveTimeout(time.Second)
    helpers.PanicOnError(pubsub_error)
    log.Debug(subscription)

    for condition := true; condition; {
        new_message, subscription_error := pubsub_client.ReceiveMessage()
        helpers.PanicOnError(subscription_error)
        go_channel <- new_message.Payload
        time.Sleep(time.Millisecond)
    }
    defer close(go_channel)
}

func (transport *Redis) Take(list_name string, go_channel chan string) {
    for condition := true; condition; {
        element, list_pop_error := transport.client.BLPop(0, list_name).Result()
        helpers.PanicOnError(list_pop_error)
        go_channel <- element[1]
    }
}


