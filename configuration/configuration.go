package stats

import (
    "fmt"

    "github.com/jmespath/go-jmespath"
    "gopkg.in/yaml.v2"

    "github.com/georgepsarakis/go-stats/helpers"
)

const ConfigurationData = `
sources:
  test_1:
    ttl: 60
    load:
      pubsub:
        channel: test_1_golang
    extraction:
      module: jmespath
      parameters:
        field_expression: "a.b.c"
        value_expression: "a.b.d"
sinks:
  url: redis://
  frequency: 30
`

type Configuration struct {
    settings interface{}
}

func (configuration *Configuration) Load() map[string]interface{} {
    var settings interface{}
    parse_error := yaml.Unmarshal([]byte(ConfigurationData), &settings)
    helpers.PanicOnError(parse_error)
    return helpers.MapWithStringKeys(settings)
}

func (configuration *Configuration) FetchNestedSettings(expression string) interface{} {
    settings := configuration.Load()
    result, search_error := jmespath.Search(expression, settings)
    fmt.Println("expression", expression, "->", result)
    helpers.PanicOnError(search_error)
    return result
}

func (configuration *Configuration) GetSettingsBySource(name string) map[string]interface{} {
    return configuration.FetchNestedSettings("sources." + name).(map[string]interface{})
}

func (configuration *Configuration) GetSources() map[string]interface{} {
    return configuration.FetchNestedSettings("sources").(map[string]interface{})
}

func (configuration *Configuration) GetSourceNames() []string {
    source_names := make([]string, 0)
    for name, _ := range configuration.GetSources() {
        fmt.Println("Found source with name ->", name)
        source_names = append(source_names, name)
    }
    return source_names
}


