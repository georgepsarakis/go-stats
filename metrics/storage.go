package metrics

import (
    "time"

    datasinks "github.com/georgepsarakis/go-stats/sinks"
)

type Metric struct {
    Name string
    Value uint64
}

type InternalStorage struct {
    hash map[string]uint64
    ttl int64
    last_reset int64
}


func (metrics_storage *InternalStorage) initialize(ttl int64) {
    metrics_storage.ttl = ttl
    metrics_storage.last_reset = time.Now().Unix()
    metrics_storage.hash = make(map[string]uint64)
}


func (metrics_storage *InternalStorage) maybeExpire() {
  current_timestamp := time.Now().Unix()

  if current_timestamp - metrics_storage.last_reset >= metrics_storage.ttl {
      metrics_storage.last_reset = current_timestamp
      metrics_storage.hash = make(map[string]uint64)

  }
}


func (metrics_storage *InternalStorage) Increment(name string, by uint64) {
    if _, ok := metrics_storage.hash[name]; !ok {
        metrics_storage.hash[name] = 0
    }
    metrics_storage.hash[name] += by
}


func (metrics_storage *InternalStorage) SnapshotToSink(sink *datasinks.Sink) {
    for name, metric := range metrics_storage.hash {
        sink.Set(name, metric)
    }
}


