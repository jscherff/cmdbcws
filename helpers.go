// Copyright 2017 John Scherff
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

package main

import (
	`encoding/json`
	`os`
	`github.com/google/gousb`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

// Load settings from JSON file into configuration object.
func load(t interface{}, cf string) (error) {

	if fh, err := os.Open(cf); err != nil {
		return err
	} else {
		defer fh.Close()
		return json.NewDecoder(fh).Decode(&t)
	}
}

// Scan for USB human input devices.
func scan(inc *Include) (dms []map[string]interface{}, err error) {

	ctx := gousb.NewContext()
	defer ctx.Close()

	filter := func(desc *gousb.DeviceDesc) bool {

		vid, pid := desc.Vendor.String(), desc.Product.String()

		if dev, ok := inc.ProductID[vid][pid]; ok {
			return dev
		}

		if dev, ok := inc.VendorID[vid]; ok {
			return dev
		}

		return inc.Default
	}

	devs, _ := ctx.OpenDevices(filter)

	for _, dev := range devs {

		var dm map[string]interface{}

		if pdev, err := probe(dev); err != nil {
			return nil, err
		} else if j, err := pdev.JSON(); err != nil {
			return nil, err
		} else if err := json.Unmarshal(j, &dm); err != nil {
			return nil, err
		} else {
			dms = append(dms, dm)
		}
	}

	return dms, nil
}

// Probe USB human input device for configuration settings and attributes.
func probe(dev *gousb.Device) (usb.Reporter, error) {

	switch {

	case usb.IsMagtek(dev.Desc.Vendor, dev.Desc.Product):
		return usb.NewMagtek(dev)

	case usb.IsIDTech(dev.Desc.Vendor, dev.Desc.Product):
		return usb.NewIDTech(dev)

	default:
		return usb.NewGeneric(dev)
	}
}
