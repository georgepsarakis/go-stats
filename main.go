package main

import (
    "fmt"

    "github.com/jmespath/go-jmespath"
    log "github.com/sirupsen/logrus"

    config "github.com/georgepsarakis/go-stats/configuration"
    datasources "github.com/georgepsarakis/go-stats/sources"
    "github.com/georgepsarakis/go-stats/metrics"
    "github.com/georgepsarakis/go-stats/extraction"
    "github.com/georgepsarakis/go-stats/transports"
)


func ConstructSource(name string, configuration *config.Configuration) (source datasources.Redis) {
    source = datasources.Redis{
        name,
        make(chan string, 1000),
        new(transports.Redis),
        "",
        "",
    }
    // source.transport = new(transports.Redis)
    source.Transport.Connect("localhost:6379", 0) 
    loader_settings := configuration.FetchNestedSettings(fmt.Sprintf("sources.%s.load", name)).(map[string]interface{})
    pubsub_channel, _ := jmespath.Search("pubsub.channel", loader_settings)

    go source.Transport.Listen(pubsub_channel.(string), source.Messages)

    extraction_settings := configuration.FetchNestedSettings(fmt.Sprintf("sources.%s.extraction.parameters", name)).(map[string]interface{})
    name_extraction_expression := extraction_settings["field_expression"].(string)
    source.NameExtractionExpression = name_extraction_expression
    value_extraction_expression := extraction_settings["value_expression"].(string)
    source.ValueExtractionExpression = value_extraction_expression
    return
}


func main() {
    configuration := new(config.Configuration)

    sources := make([]datasources.Redis, 0)
    channels := make([]chan string, 0)

    source_names := configuration.GetSourceNames()
    log.Info(source_names)

    for _, name := range source_names {
        source := ConstructSource(name, configuration)
        sources = append(sources, source)
        channels = append(channels, source.Messages)
    }

    combined_channel := make(chan metrics.Metric, 10000)
    go extraction.Collect(sources, combined_channel, channels...)

    for metric := range combined_channel {
        log.Info(fmt.Sprintf("Received metric %s = %d", metric.Name, metric.Value))
    }
}
