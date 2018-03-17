package datasources

import (
    transports "github.com/georgepsarakis/go-stats/transports"
)

/*
TODO: use Source interface
*/

type Redis struct {
    Name string
    /*
    url string
    parsed_url *url.URL
    */
    Messages chan string
    Transport *transports.Redis
    // TODO: create Extractor struct
    NameExtractionExpression string
    ValueExtractionExpression string
}


func (source *Redis) IsNamed(name string) bool {
    return source.Name == name
}


func (source *Redis) UsesChannel(channel chan string) bool {
    return channel == source.Messages
}
