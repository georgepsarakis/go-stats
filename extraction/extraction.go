package extraction

import (
    "encoding/json"
    "errors"
    "fmt"

    "github.com/jmespath/go-jmespath"

    datasources "github.com/georgepsarakis/go-stats/sources"
    "github.com/georgepsarakis/go-stats/metrics"
    "github.com/georgepsarakis/go-stats/helpers"
)


func ExtractFromMessage(expression string, raw_message string) (interface{}, error) {
    var parsed_json_message interface{}
    json.Unmarshal([]byte(raw_message), &parsed_json_message)
    return jmespath.Search(expression, parsed_json_message)
}


func FindSourceForChannel(sources []datasources.Redis, channel chan string) (datasources.Redis, error) {
    var source datasources.Redis

    for _, source := range sources {
        if source.UsesChannel(channel) {
            return source, nil
        }
    }

    return source, errors.New("No source could be found")
}


func Collect(sources []datasources.Redis, combined_channel chan metrics.Metric, channels ...chan string) {
    fetch_next_available_value := helpers.DynamicChannelSelect(channels...)
    for {
        selected_channel_index, message, _ := fetch_next_available_value()
        source, search_error := FindSourceForChannel(sources, channels[selected_channel_index])
        helpers.PanicOnError(search_error)
        field_name, _ := ExtractFromMessage(source.NameExtractionExpression, message)
        value, _ := ExtractFromMessage(source.ValueExtractionExpression, message)
fmt.Println(value)
        transformed_value, conversion_error := helpers.TransformToNumber(value)
        helpers.PanicOnError(conversion_error)
        combined_channel <- metrics.Metric{field_name.(string), transformed_value}
    }
}


