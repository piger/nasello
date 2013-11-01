// Configuration handling
//
// The configuration file is a JSON file with a simple structure; the following
// configuration specify 3 forwarders: *.example.com and 10.1.2.* will be
// resolved by OpenDNS and a catch-all for resolving with Google DNS.
// 
// {
//		"filters": [
// 				{
// 						"pattern": "example.com.",
// 						"addresses": [
// 								"208.67.222.222",
// 								"208.67.220.220"
// 						]
// 				},
// 				{
// 						"pattern": "2.1.10.in-addr.arpa.",
// 						"addresses": [
// 								"208.67.222.222",
// 								"208.67.220.220"
// 						]
// 				},
// 				{
// 						"pattern": ".",
// 						"addresses": [
// 								"8.8.8.8"
// 						]
// 				}
// 		]
// }
//
package nasello

import (
	"encoding/json"
	"log"
	"io/ioutil"
)


type Configuration struct {
	Filters []ConfigFilter
}

type ConfigFilter struct {
	Pattern string
	Addresses []string	
}

func ReadConfig(filename string) Configuration {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't open config file: %s\n", err.Error())
	}

	var jsonConfig Configuration
	json.Unmarshal(file, &jsonConfig)

	return jsonConfig
}
