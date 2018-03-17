package helpers

import (
    "reflect"
    "strconv"
)


func IsNil(variable interface{}) bool {
  return variable == nil
}


func PanicOnError(possible_error error) {
    if !IsNil(possible_error) {
        panic(possible_error)
    }
}


func MapWithStringKeys(map_instance interface{}) map[string]interface{} {
    new_map := make(map[string]interface{})

    for key, value := range map_instance.(map[interface{}]interface{}) {
        value_type := reflect.ValueOf(value)
        if value_type.Kind() == reflect.Map {
            new_map[key.(string)] = MapWithStringKeys(value)
        } else {
            new_map[key.(string)] = value
        }
    }
    return new_map
}


func TransformToNumber(value interface{}) (uint64, error) {
    switch value.(type) {
    case string:
        return strconv.ParseUint(value.(string), 10, 64)
    case float64:
        return uint64(value.(float64)), nil
    default:
        return value.(uint64), nil
    }
}

