package helpers

import (
    "errors"
    "reflect"
)


func DynamicChannelSelect(channels ...chan string) func() (int, string, error) {
    cases := make([]reflect.SelectCase, len(channels))
    for index, channel := range channels {
        cases[index] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)}
    }
    return func() (int, string, error) {
        var selection_channel_error error

        selected_channel_index, message, channel_closed := reflect.Select(cases)
        if channel_closed {
            selection_channel_error = errors.New("Channel closed")
        } else {
            selection_channel_error = nil
        }
        return selected_channel_index, message.String(), selection_channel_error
    }
}


