package dto

import (
	"database/sql/driver"
	"encoding/json"
)

// Props is a map of key:value
type Props map[string]interface{}

// Value converts props into a database value
func (p Props) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Merge current props with given props
func (p Props) Merge(props Props) Props {
	new := Props{}
	for k, v := range p {
		new[k] = v
	}
	for k, v := range props {
		new[k] = v
	}
	return new
}

// Append add the given props to current props
func (p Props) Append(props Props) {
	for k, v := range props {
		p[k] = v
	}
}
