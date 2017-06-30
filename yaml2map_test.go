package yaml2map

import (
	"testing"
	"reflect"
	"fmt"
)

const data1 = `
Colors:
  - red: red
  - pink:
      mix:
        - white
        - red
      main: false
  - Brown:
      - mix:
          - red
          - black
      - used:
          inside: false
          outside: true
  - blue: blue
  - white: white
nonColor1: black
nonColor2: white
`

const data2 = `
---
collectd::graphite_url: 'h1-int-gpr.ovp.bskyb.com'
collectd::graphite_prefix: 'graphite-stg1-h1.'
collectd::repo_source: 'http://h1rep01-v00.devops.int.ovp.bskyb.com'
common::icinga::icinga_host_tag: 'h1_devops_int'
common::icinga::custom_pool: 'int'

logstash::filebeat::indexers:
  - "h1-int-lsi.ovp.bskyb.com:8080"

kubernetes::devops_docker_registry: 'h1drg01-v01.devops.int.ovp.bskyb.com'
`
func TestYaml2Map(t *testing.T) {
    
  //Test1
  data1_map := map[string]interface{}{
  	"Colors/red": "red",
  	"Colors/pink/main": "false",
  	"Colors/pink/mix": "white,red",
  	"Colors/Brown/mix": "red,black",
  	"Colors/Brown/used/inside": "false",
  	"Colors/Brown/used/outside": "true",
  	"Colors/blue": "blue",
  	"Colors/white": "white",
    "nonColor1":"black",
    "nonColor2":"white",
  }

	ret_data1_map := Yaml2Map([]byte(data1))
	
	eq := reflect.DeepEqual(ret_data1_map, data1_map)
	if eq {
    	fmt.Println("Test1 passed, both maps are equal.")
	} else {
      fmt.Println("Expecting : ", data1_map)
      fmt.Println("GOT: ", ret_data1_map)
      t.Fatalf("Maps are unequal.")
	}

  //Test2
  //Test1
  data2_map := map[string]interface{}{
    "collectd::graphite_url": "h1-int-gpr.ovp.bskyb.com",
    "collectd::graphite_prefix": "graphite-stg1-h1.",
    "collectd::repo_source": "http://h1rep01-v00.devops.int.ovp.bskyb.com",
    "common::icinga::icinga_host_tag": "h1_devops_int",
    "common::icinga::custom_pool": "int",
    "logstash::filebeat::indexers": "h1-int-lsi.ovp.bskyb.com:8080",
    "kubernetes::devops_docker_registry": "h1drg01-v01.devops.int.ovp.bskyb.com",
  }

  ret_data2_map := Yaml2Map([]byte(data2))
  
  eq2 := reflect.DeepEqual(ret_data2_map, data2_map)
  if eq2 {
      fmt.Println("Test1 passed, both maps are equal.")
  } else {
      fmt.Println("Expecting : ", data2_map)
      fmt.Println("GOT: ", ret_data2_map)
      t.Fatalf("Maps are unequal.")
  }



}