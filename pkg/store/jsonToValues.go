// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"encoding/json"
	"fmt"
	devicechange "github.com/onosproject/onos-config/api/types/change/device"
)

// DecomposeTree breaks a JSON file down in to paths and values without any external
// context - a second stage is required to align these paths to the specific devicechange
// got from the Read Only paths of the Model. Also the second stage is necessary
// resolve list indices, and leaf lists. This second stage is done in the model
// registry to avoid circular dependencies here.
func DecomposeTree(genericJSON []byte) ([]*devicechange.PathValue, error) {
	var f interface{}
	err := json.Unmarshal(genericJSON, &f)
	if err != nil {
		return nil, err
	}
	values := extractValuesIntermediate(f, "")
	return values, nil
}

// extractValuesIntermediate recursively walks a JSON tree to create a flat set
// of paths and values.
// Note: it is not possible to find indices of lists and accurate devicechange directly
// from json - for that the RO Paths must be consulted
func extractValuesIntermediate(f interface{}, parentPath string) []*devicechange.PathValue {
	changes := make([]*devicechange.PathValue, 0)

	switch value := f.(type) {
	case map[string]interface{}:
		for key, v := range value {
			objs := extractValuesIntermediate(v, fmt.Sprintf("%s/%s", parentPath, key))
			changes = append(changes, objs...)
		}
	case []interface{}:
		// Iterate through to look for indexes first
		for idx, v := range value {
			objs := extractValuesIntermediate(v, fmt.Sprintf("%s[%d]", parentPath, idx))
			changes = append(changes, objs...)
		}
	case string:
		newCv := devicechange.PathValue{Path: parentPath, Value: devicechange.NewTypedValueString(value)}
		changes = append(changes, &newCv)
	case bool:
		newCv := devicechange.PathValue{Path: parentPath, Value: devicechange.NewTypedValueBool(value)}
		changes = append(changes, &newCv)
	case float64:
		newCv := devicechange.PathValue{Path: parentPath, Value: devicechange.NewTypedValueFloat(float32(value))}
		changes = append(changes, &newCv)
	default:
		fmt.Println("Unexpected type", value)
	}

	return changes
}
