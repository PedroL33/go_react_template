package util

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// ScanRowIntoStruct is a generic function to scan a single row into a struct.
func ScanRowIntoStruct(row pgx.Row, dest interface{}) error {
	// Ensure the destination is a pointer to a struct
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	// Get the struct type and value
	destType := destVal.Elem().Type()

	// Prepare the slice of arguments for the Scan method
	var args []interface{}

	// Create a list of pointers to the struct fields
	for i := 0; i < destType.NumField(); i++ {
		field := destVal.Elem().Field(i)
		if field.CanSet() {
			args = append(args, field.Addr().Interface())
		}
	}

	// Perform the Scan operation
	err := row.Scan(args...)
	if err != nil {
		return errors.Wrap(err, "util.ScanRowIntoStruct")
	}

	return nil
}
