# DHT
DHT represents a client for a distributed hash table.
```
type DHT struct {
	// Has unexported fields.
}

func New(node string, receivePort uint16, bits uint64) (*DHT, error)
func (dht *DHT) Get(k string) (string, error)
func (dht *DHT) Put(k string, v string) error
func (dht *DHT) Start()
```

The test for DHT can be found [here](../test/dht/dht_main.go)