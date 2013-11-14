package simplejson

import (
	"encoding/json"
	"errors"
	"log"
	"bytes"
	"strings"
	"io/ioutil"
)

// returns the current implementation version
func Version() string {
	return "0.4.5"
}

type Json struct {
	data interface{}
}

// NewJson returns a pointer to a new `Json` object
// after unmarshaling `body` bytes
func NewJson(body []byte) (*Json, error) {
	j := new(Json)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// NewJsonFromFile return a pointer to a new `Json` object
// after unmarshalling a json file
// forked from github.com/polaris1119/autogo/src/simplejson
func NewJsonFromFile(filename string) (*Json, error) {
    stream, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    content := string(stream)
    var builder bytes.Buffer
    lines := strings.Split(content, "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        builder.WriteString(line)
    }
    return NewJson(builder.Bytes())
}

// Set Data of Json
func (j *Json) SetData(v interface{}) *Json {
  j.data = v
  return j
}

// Get Data of Json
func (j *Json) GetData() interface{} {
  return j.data
}

// Get Num of Data
func (j *Json) Count() (int, error) {
  m, err := j.Map()
  if err != nil {
    return -1, err
  }
  return len(m), nil
}

// Convert data of Json to string
func (j *Json) ToString() (string, error) {
  if j.IsNull() {
    return "", errors.New("data is null")
  }

  b, err := j.MarshalJSON()
  if err == nil {
    return string(b), nil
  }
  return "", err
}

// Check whether its data is nil
func (j *Json) IsNull() bool {
  if j.data == nil {
    return true
  }

  // check whether it's a null map
  if m, ok := j.Map(); ok == nil {
    return len(m) == 0 
  }
  return false
}

// Encode returns its marshaled data as `[]byte`
func (j *Json) Encode() ([]byte, error) {
	return j.MarshalJSON()
}

// Implements the json.Unmarshaler interface.
func (j *Json) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, &j.data)
}

// Implements the json.Marshaler interface.
func (j *Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// Set modifies `Json` map by `key` and `value`
// Useful for changing single key/value in a `Json` object easily.
func (j *Json) Set(key string, val interface{}) {
	m, err := j.Map()
	if err != nil {
		return
	}

  switch val.(type) {
    case *Json:
      val = val.(*Json).data
  }

	m[key] = val
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *Json) Get(key string) *Json {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}
		}
	}
	return &Json{nil}
}

// Delete a value for a specified key
func (j *Json) Delete(key string) *Json {
  m, err := j.Map()
  if err == nil {
    if _, ok := m[key]; ok {
      delete(m, key)
    }
  }
  return j
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (j *Json) GetPath(branch ...string) *Json {
	jin := j
	for i := range branch {
		m, err := jin.Map()
		if err != nil {
			return &Json{nil}
		}
		if val, ok := m[branch[i]]; ok {
			jin = &Json{val}
		} else {
			return &Json{nil}
		}
	}
	return jin
}

// GetIndex resturns a pointer to a new `Json` object
// for `index` in its `array` representation
//
// this is the analog to Get when accessing elements of
// a json array instead of a json object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (j *Json) GetIndex(index int) *Json {
	a, err := j.Array()
	if err == nil {
		if len(a) > index {
			return &Json{a[index]}
		}
	}
	return &Json{nil}
}

// CheckGet returns a pointer to a new `Json` object and
// a `bool` identifying success or failure
//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
//        log.Println(data)
//    }
func (j *Json) CheckGet(key string) (*Json, bool) {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}, true
		}
	}
	return nil, false
}

// Map type asserts to `map`
func (j *Json) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// Array type asserts to an `array`
func (j *Json) Array() ([]interface{}, error) {
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// Bool type asserts to `bool`
func (j *Json) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (j *Json) String() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Float64 type asserts to `float64`
func (j *Json) Float64() (float64, error) {
	if i, ok := (j.data).(float64); ok {
		return i, nil
	}
	return -1, errors.New("type assertion to float64 failed")
}

// Int type asserts to `float64` then converts to `int`
func (j *Json) Int() (int, error) {
  switch f := (j.data).(type) {
  case float64:
    return int(f), nil
  case int:
    return f, nil
  case int64:
    return int(f), nil
  }

	return -1, errors.New("type assertion to int failed")
}

// Int type asserts to `float64` then converts to `int64`
func (j *Json) Int64() (int64, error) {
  switch f := (j.data).(type) {
  case float64:
    return int64(f), nil
  case int:
    return int64(f), nil
  case int64:
    return f, nil
  }

	return -1, errors.New("type assertion to int64 failed")
}

// Bytes type asserts to `[]byte`
func (j *Json) Bytes() ([]byte, error) {
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

// StringArray type asserts to an `array` of `string`
func (j *Json) StringArray() ([]string, error) {
	switch arr := j.data.(type) {
  case []interface{}:
	  retArr := make([]string, 0, len(arr))
    for _, a := range arr {
      s, ok := a.(string)
      if !ok {
        return nil, errors.New("type assertion to string failed")
      }
      retArr = append(retArr, s)
	  }
	  return retArr, nil
	case []string:
	  retArr := make([]string, 0, len(arr))
    for _, a := range arr {
      retArr = append(retArr, a) 
    }
    return retArr, nil
  }
	return nil, errors.New("type assertion to []string failed")
}

// IntArray type asserts to an `array` of `int64`
func (j *Json) Int64Array() ([]int64, error) {
  switch arr := j.data.(type) {
  case []interface{}:
	  retArr := make([]int64, 0, len(arr))
    for _, a := range arr {
      s, ok := a.(float64)
      if !ok {
        return nil, errors.New("type assertion to float64 failed")
      }
      retArr = append(retArr, int64(s))
	  }
	  return retArr, nil
	case []int64:
	  retArr := make([]int64, 0, len(arr))
    for _, a := range arr {
      retArr = append(retArr, a) 
    }
    return retArr, nil
  case []int:
    retArr := make([]int64, 0, len(arr))
    for _, a := range arr {
      retArr = append(retArr, int64(a))
    }
    return retArr, nil
  }
	return nil, errors.New("type assertion to []int64 failed")
}

// IntArray type asserts to an `array` of `int`
func (j *Json) IntArray() ([]int, error) {
  switch arr := j.data.(type) {
  case []interface{}:
	  retArr := make([]int, 0, len(arr))
    for _, a := range arr {
      s, ok := a.(float64)
      if !ok {
        return nil, errors.New("type assertion to float64 failed")
      }
      retArr = append(retArr, int(s))
	  }
	  return retArr, nil
	case []int64:
	  retArr := make([]int, 0, len(arr))
    for _, a := range arr {
      retArr = append(retArr, int(a)) 
    }
    return retArr, nil
  case []int:
    retArr := make([]int, 0, len(arr))
    for _, a := range arr {
      retArr = append(retArr, a)
    }
    return retArr, nil
  }
	return nil, errors.New("type assertion to []int failed")
}

// MustArray guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, v := range js.Get("results").MustArray() {
//			fmt.Println(i, v)
//		}
func (j *Json) MustArray(args ...[]interface{}) []interface{} {
	var def []interface{}
	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustArray() received too many arguments %d", len(args))
	}

	a, err := j.Array()
	if err == nil {
		return a
	}

	return def
}

// MustMap guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").MustMap() {
//			fmt.Println(k, v)
//		}
func (j *Json) MustMap(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}
	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustMap() received too many arguments %d", len(args))
	}

	a, err := j.Map()
	if err == nil {
		return a
	}

	return def
}

// MustString guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").MustString(), js.Get("optional_param").MustString("my_default"))
func (j *Json) MustString(args ...string) string {
	var def string

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustString() received too many arguments %d", len(args))
	}

	s, err := j.String()
	if err == nil {
		return s
	}

	return def
}

// MustInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (j *Json) MustInt(args ...int) int {
	var def int

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}

	i, err := j.Int()
	if err == nil {
		return i
	}

	return def
}

// MustInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (j *Json) MustInt64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}

	i, err := j.Int64()
	if err == nil {
		return i
	}

	return def
}

// MustFloat64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.Get("param1").MustFloat64(), js.Get("optional_param").MustFloat64(5.150))
func (j *Json) MustFloat64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustFloat64() received too many arguments %d", len(args))
	}

	i, err := j.Float64()
	if err == nil {
		return i
	}

	return def
}
