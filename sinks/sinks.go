package datasinks

import (
    "net/url"

    "github.com/go-redis/redis"

    "github.com/georgepsarakis/go-stats/helpers"
)

type Sink struct {
    url string
    parsed_url *url.URL
    client *redis.Client
}

func (sink *Sink) connect() {
    parsed_url, parse_error := url.Parse(sink.url)
    helpers.PanicOnError(parse_error)
    sink.parsed_url = parsed_url
    sink.client = redis.NewClient(&redis.Options{})
}

func (sink *Sink) Set(name string, value uint64) {
    sink.client.Set(name, value, 0)
}
