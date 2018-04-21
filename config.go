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
	`fmt`
	`net/http`
	`log`
	`html/template`
	`os`
	`path/filepath`
)

const configFile = `cmdbcws.json`

var conf = new(Config)

type Config struct {
	Hostname string
	Server *http.Server
	Include *Include
	Template *template.Template
	Templates []string
	Resources []string
}

type Include struct {
	VendorID	map[string]bool
	ProductID	map[string]map[string]bool
	Default		bool
}

func init() {

	log.SetFlags(log.Lshortfile)

	appDir := filepath.Dir(os.Args[0])

	if err := load(conf, filepath.Join(appDir, configFile)); err != nil {
		log.Fatal(err)
	}

	if hn, err := os.Hostname(); err != nil {
		log.Fatal(err)
	} else {
		conf.Hostname = hn
	}

	for index, file := range conf.Templates {
		conf.Templates[index] = filepath.Join(appDir, file)
	}

	if tmpl, err := template.ParseFiles(conf.Templates...); err != nil {
		log.Fatal(err)
	} else {
		conf.Template = tmpl
	}

	for _, dir := range conf.Resources {
		fs := http.FileServer(http.Dir(dir))
		path := fmt.Sprintf(`/%s/`, dir)
		http.Handle(path, http.StripPrefix(path, fs))
	}
}
