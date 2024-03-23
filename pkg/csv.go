package github

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
)

func SaveToCSV(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Ensure we're dealing with a slice.
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("data must be a pointer to a slice")
	}

	// Get the slice and ensure it's not empty.
	slice := v.Elem()
	if slice.Len() == 0 {
		return fmt.Errorf("slice is empty")
	}

	firstElem := slice.Index(0)
	header, err := structTagsToSlice(firstElem.Interface())
	if err != nil {
		return err
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i).Interface()
		row, err := structToStringSlice(elem)
		if err != nil {
			return err
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// structToStringSlice converts struct fields to string slice.
func structToStringSlice(data interface{}) ([]string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data must be a struct")
	}

	var row []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		row = append(row, fmt.Sprintf("%v", field.Interface()))
	}
	return row, nil
}

// structTagsToSlice extracts JSON tags from struct fields.
func structTagsToSlice(data interface{}) ([]string, error) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data must be a struct")
	}

	var header []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		header = append(header, tag)
	}
	return header, nil
}
