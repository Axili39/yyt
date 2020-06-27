package godict

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Dict : Dictionnary of anything type
type Dict map[interface{}]interface{}

// Array : Array of anything type
type Array []interface{}

func merge(a interface{}, b interface{}) (interface{}, error) {
	switch x := a.(type) {
	case Dict:
		return mergeDict(x, b)
	case Array:
		return mergeArray(x, b)
	default:
		return a, nil
	}
}

func mergeDict(a Dict, b ...interface{}) (Dict, error) {
	ret := Dict{}
	for k, v := range a {
		ret[k] = v
	}

	for _, bi := range b {
		if bi == nil {
			continue
		}

		db, ok := bi.(Dict)
		if !ok {
			return nil, fmt.Errorf("Dict element expected during merging Dict")
		}
		var err error
		for k, v := range db {
			if ret[k] != nil {
				ret[k], err = merge(ret[k], v)
				if err != nil {
					return nil, err
				}
			} else {
				ret[k] = v
			}
		}
	}
	return ret, nil
}

func mergeArray(a Array, b ...interface{}) (Array, error) {
	retlen := len(a)

	// Pre-process b to compute len
	for _, bi := range b {
		ab := bi.(Array)
		if ab == nil {
			return nil, fmt.Errorf("Array element expected during merging Array")
		}
		retlen += len(ab)
	}
	// create Array
	ret := make([]interface{}, retlen)
	for i, v := range a {
		ret[i] = v
	}
	offset := len(a)

	for _, bi := range b {
		ab := bi.(Array)
		for i, v := range ab {
			ret[offset+i] = v
		}
		offset += len(ab)

	}
	return ret, nil
}

// Merge : Merge 2 dictionnaries
func Merge(a Dict, b Dict) (Dict, error) {
	if a == nil {
		return b, nil
	}

	if b == nil {
		return a, nil
	}

	return mergeDict(a, b)
}

// LoadFromYamlFiles : Load dictionnary from YAML encoded files
func LoadFromYamlFiles(filenames []string) (Dict, error) {
	var data Dict
	for _, f := range filenames {
		yamlFile, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("Error loading %s : %v", f, err)
		}

		var df Dict
		err = df.FromYamlData(yamlFile)
		if err != nil {
			return nil, fmt.Errorf("Error Unmarshalling %s : %v", f, err)
		}
		data, _ = Merge(data, df)
	}
	return data, nil
}

// FromYamlData : unMarshal data into Dict
func (d *Dict) FromYamlData(data []byte) error {
	return yaml.Unmarshal(data, d)
}

// ToYamlData : Marshal Dict into data
func (d *Dict) ToYamlData() ([]byte, error) {
	return yaml.Marshal(d)
}

// Path Management
func tokenize(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// ExtractFromXPath : Try to get a node from path
func (d Dict) ExtractFromXPath(path string) (interface{}, error) {
	// get Token
	head, tail := tokenize(path)
	var node interface{} = d
	for {
		switch x := node.(type) {
		case Dict:
			node = x[head]
		case Array:
			index, _ := strconv.Atoi(head)
			node = x[index]
		default:
			return nil, fmt.Errorf("can't interprete token %s at type: %v", head, x)
		}

		if tail != "/" {
			head, tail = tokenize(tail)
		} else {
			break
		}
	}
	return node, nil
}
