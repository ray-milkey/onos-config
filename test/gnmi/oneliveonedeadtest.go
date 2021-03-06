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
//

package gnmi

import (
	"github.com/onosproject/onos-config/test/utils"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"testing"
)

// TestOneLiveOneDeadDevice tests GNMI operations to an offline device followed by operations to a connected device
func (s *TestSuite) TestOneLiveOneDeadDevice(t *testing.T) {

	const (
		modPath  = "/system/clock/config/timezone-name"
		modValue = "Europe/Rome"
	)

	// Make a GNMI client to use for requests
	gnmiClient := getGNMIClientOrFail(t)

	// Set a value using gNMI client to the offline device
	offlineDevicePath := getDevicePathWithValue("offline-device", modPath, modValue, StringVal)

	// Set the value - should return a pending change
	setGNMIValueOrFail(t, gnmiClient, offlineDevicePath, noPaths, getSimulatorExtensions())

	// Check that the value was set correctly in the cache
	checkGNMIValue(t, gnmiClient, offlineDevicePath, modValue, 0, "Query after set returned the wrong value")

	// Create an online device
	onlineSimulator := env.NewSimulator().AddOrDie()

	// Set a value to the online device
	onlineDevicePath := getDevicePathWithValue(onlineSimulator.Name(), modPath, modValue, StringVal)
	nid := setGNMIValueOrFail(t, gnmiClient, onlineDevicePath, noPaths, noExtensions)
	utils.WaitForNetworkChangeComplete(t, nid)

	// Check that the value was set correctly in the cache
	checkGNMIValue(t, gnmiClient, onlineDevicePath, modValue, 0, "Query after set returned the wrong value")

	// Check that the value was set correctly on the device
	deviceGnmiClient := getDeviceGNMIClientOrFail(t, onlineSimulator)
	checkDeviceValue(t, deviceGnmiClient, onlineDevicePath, modValue)
}
