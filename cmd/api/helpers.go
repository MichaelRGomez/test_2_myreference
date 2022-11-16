// Filename: test2/cmd/api/helpers.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"michaelgomez.net/internal/validator"
)

// defining envelope type
type envelope map[string]interface{}

// pulling any "id" from json input
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// write JSON function
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// reading the json
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//limiting the size of the request body to 1MB
	maxBytes := 1_048_576

	//decoding the request body into the target destination
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)

	//ensuring there's no bad request
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError): //checking syntax
			return fmt.Errorf("body contains badly-formated JSON(at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF): //end of file error
			return errors.New("body contains badly-formatted JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for fiend %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF): //empty body
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json:unkoiwn field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contrains unkown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	//decoding again
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single value")
	}
	return nil
}

// returns string value from the query parameters
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// splits the value into slice based on the comma separator
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	//the delimiter is ","
	return strings.Split(value, ",")
}

// reads int from the parameters
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return intValue
}
