//*****************************************************************
// Package: yaml2map
// Purpose: Read data from yaml and create map of K/V
// Author: SidhuG
//*******************************************************************

package yaml2map

import (
	"fmt"
	"github.com/alediaferia/stackgo"
	"os"
	"strconv"
	"gopkg.in/yaml.v2"
	"reflect"
    //log "github.com/sirupsen/logrus"
)

type I interface{}

//map of consul keys, whole url endpoint is one key
var m_keys map[string]interface{} = make(map[string]interface{})
//var m_keys_0 map[string]interface{} = make(map[string]interface{})
var top_key string
var current_key string
var st_keys = stackgo.NewStack()
//var st_keys_0 = stackgo.NewStack()

//buffer required for logging

func Yaml2Map( data3 []byte) map[string]interface{}{

    //TODO:Initialise logger
    
    //Initilise global variables
    m_keys_0 := make(map[string]interface{})
    st_keys_0 := stackgo.NewStack()
    m_keys = m_keys_0
    st_keys = st_keys_0

	m := make(map[string]interface{})

	err := yaml.Unmarshal(data3, &m)
	checkError(err)
	fmt.Printf("--- m:\n%v\n\n", m)

	//var valueType reflect.kind

	fmt.Println("-----Printing toplevel map-----")
	for k, v := range m {
		fmt.Printf("key[%s] value[%s]\n", k, v)
		top_key = k
		current_key = k
		st_keys.Push(current_key)
		valueType := reflect.TypeOf(v).Kind()
		fmt.Printf("ValueType is %s", valueType)
		extract(v)
		fmt.Println()
	}

	for ke, val := range m_keys {
		fmt.Println("Key: ", ke)
		fmt.Println("val: ", val)
		fmt.Println()
	}
	return m_keys
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func extract(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	extractRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func extractRecursive(copy, original reflect.Value) {
	var existingValue string
	var NewValue string
	switch original.Kind() {
	// The first cases handle nested structures and extract them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		extractRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		extractRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we extract each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			extractRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and extract each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			extractRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and extract each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		st_keys.Push(current_key)
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			
			st_keys.Push(current_key)
			current_key = fmt.Sprintf("%s/%s", current_key, key)
			fmt.Println()
			fmt.Printf(" Key: %s  -> ", current_key)

			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			extractRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
			current_key = st_keys.Pop().(string)
		}
		current_key = st_keys.Pop().(string)

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string extract it (yay finally we're doing what we came for)
	case reflect.String:
		extractdString := original.Interface().(string)
		copy.SetString(extractdString)
		if val, ok := m_keys[current_key]; ok {
			existingValue = val.(string)
			NewValue = existingValue + "," + extractdString
			m_keys[current_key] = NewValue
			fmt.Printf(" Key:  %s", current_key)
			fmt.Println(" Value:  ", NewValue)
		} else {
			m_keys[current_key] = extractdString
			fmt.Printf(" Key: %s", current_key)
			fmt.Println(" Value: ", extractdString)
		}

		// A bool type will always be a value, convert it to string before saving
	case reflect.Bool:
		var tf bool = original.Bool()
		extractdString := strconv.FormatBool(tf)
		m_keys[current_key] = extractdString
		fmt.Printf(" Key: %s", current_key)
		fmt.Println(" Value: ", extractdString)
		copy.Set(original)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}
}
