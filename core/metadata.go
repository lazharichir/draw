package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Metadata map[string]any

func (m Metadata) Get(key string) any {
	return m[key]
}

func (m Metadata) Set(key string, value any) {
	m[key] = value
}

func (m Metadata) Delete(key string) {
	delete(m, key)
}

func (m Metadata) Has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Metadata) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// json marshalling and unmashalling
func (m Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(m))
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	var aux map[string]any
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*m = Metadata(aux)
	return nil
}

// sql scanner and valuer interfaces
func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, m)
	case string:
		return json.Unmarshal([]byte(src), m)
	default:
		return fmt.Errorf("cannot unmarshal %T into Metadata", src)
	}
}
