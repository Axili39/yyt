package godict

import (
	"fmt"
	"testing"
)

const data = `---
version: "1.0"
object:
  member1: A1
  member2: B1
  array:
  - item1
  - item2`

const data2 = `---
version: "2.0"
object:
  member1: A2
  member3: C2
  array:
  - item1
  - item3`

const unmergable = `---
version: "3.0"
object:
  array: not an array
`

func TestUnmarshal(t *testing.T) {
	var d Dict
	err := d.FromYamlData([]byte(data))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}

	// version should exist
	if d["version"] == nil {
		t.Errorf("missing version node %v", d)
	}
	if d["version"] != "1.0" {
		t.Errorf("bad version node %v", d)
	}

	// object should exist
	if d["object"] == nil {
		t.Errorf("missing object node %v", d)
	}

	if d["object"].(Dict)["member1"] == nil {
		t.Errorf("missing object/member1 node %v", d)
	}

	// array should exists
	if d["object"].(Dict)["array"] == nil {
		t.Errorf("missing array node %v", d)
	}
	if len(d["object"].(Dict)["array"].([]interface{})) != 2 {
		t.Errorf("bad array count %v", d)
	}
}

func TestExtractFromXPath(t *testing.T) {
	var d Dict
	err := d.FromYamlData([]byte(data))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}

	// nominal object node extraction
	node, err := d.ExtractFromXPath("object/member1")
	if err != nil {
		t.Errorf("Extract from XPath failed")
	}
	if node != "A1" {
		t.Errorf("Unexpected Extract Value")
	}

	// nominal array element extraction
	node, err = d.ExtractFromXPath("object/array/1")
	if err != nil {
		t.Errorf("Extract from XPath failed")
	}
	if node != "item2" {
		t.Errorf("Unexpected Extract Value")
	}

	// unexisting node extraction
	node, err = d.ExtractFromXPath("object/member3")
	if err == nil {
		t.Errorf("Unexisting node extraction must return an error")
	}

	// out of bounds array element extraction
	node, err = d.ExtractFromXPath("object/array/56")
	if err == nil {
		t.Errorf("Unexisting node extraction from array must return an error")
	}

	// bad array element index
	node, err = d.ExtractFromXPath("object/array/xyz")
	if err == nil {
		t.Errorf("Bad node extraction from array must return an error")
	}
}

func TestMerge(t *testing.T) {
	var d1, d2 Dict
	err := d1.FromYamlData([]byte(data))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}
	err = d2.FromYamlData([]byte(data2))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}

	d, err := Merge(d1, d2)
	if err != nil {
		fmt.Errorf("Unexpected error during nominal merge : %v", err)
	}

	// check d content

	// version should exist
	if d["version"] == nil {
		t.Errorf("missing version node %v", d)
	}
	if d["version"] != "1.0" {
		t.Errorf("bad version node %v", d)
	}

	// object members
	if d["object"] == nil || d["object"].(Dict)["member1"] != "A1" || d["object"].(Dict)["member2"] != "B1" || d["object"].(Dict)["member3"] != "C2" {
		t.Errorf("object merge failed %v", d)
	}

	// array
	if d["object"] == nil || d["object"].(Dict)["array"] == nil {
		t.Errorf("missing array in merged result")
	}
	a := d["object"].(Dict)["array"].([]interface{})
	if a[0] != "item1" || a[1] != "item2" || a[2] != "item1" || a[3] != "item3" {
		t.Errorf("unexpected array merged result")
	}
}

func TestUnMergable(t *testing.T) {
	var d1, d2 Dict
	err := d1.FromYamlData([]byte(data))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}
	err = d2.FromYamlData([]byte(unmergable))
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
	}

	_, err = Merge(d1, d2)
	if err == nil {
		fmt.Errorf("Missing error when merging uncomptatible datas")
	}
}

func TestLoad(t *testing.T) {
	filenames := []string{"tests/net1.yml", "tests/net2.yml", "tests/net3.yml"}
	_, err := LoadFromYamlFiles(filenames)
	if err != nil {
		fmt.Errorf("Unexpected error during nominal load : %v", err)
	}

}
