package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var v *validator.Validate

func init() {
	v = validator.New()
	v.RegisterValidation("ISO8601Date", IsISO8601Date)
	v.RegisterValidation("sanitize", Sanitize)
}

func Validate(c echo.Context, s interface{}) (err error) {
	ctx := c.Request().Context()

	if err = c.Bind(s); err != nil {
		errMsg := err.Error()
		slog.ErrorContext(ctx, "error bind", err.Error())
		if strings.Contains(errMsg, "invalid syntax") {
			re := regexp.MustCompile(`"([^"]+)"`)
			matches := re.FindAllStringSubmatch(errMsg, -1)
			if len(matches) > 1 {
				if len(matches[1]) > 1 {
					err = fmt.Errorf("data type for %s is invalid", matches[1][1])
				}
			}
			c.Set("invalid-format", true)
			return
		}
		err = fmt.Errorf("%s", "Something Went Wrong")
		return
	}

	errVal := c.Validate(s)
	if errVal != nil {
		err = castedValidate(errVal.(validator.ValidationErrors))
		slog.ErrorContext(ctx, "error validate [a]", err.Error())
		c.Set("invalid-format", true)
		return
	}

	return
}

func ValidateAll(c echo.Context, s interface{}) (err error) {
	ctx := c.Request().Context()

	if err = c.Bind(s); err != nil {
		slog.ErrorContext(ctx, "error bind", err.Error())
		err = fmt.Errorf("%s", "Something Went Wrong")
		return
	}

	errVal := c.Validate(s)
	structType := getStructType(s)
	if errVal != nil {
		errObj := castedValidateAll(errVal.(validator.ValidationErrors), structType)
		errStr, _ := json.Marshal(errObj)
		err = fmt.Errorf("%s", errStr)
		slog.ErrorContext(ctx, "error validate [b]", err)
		c.Set("invalid-format", true)
		return
	}

	return
}

func ValidateStruct(ctx context.Context, s interface{}) (err error) {
	return v.StructCtx(ctx, s)
}

func ValidateAllStruct(ctx context.Context, s interface{}) (err error) {
	errVal := v.StructCtx(ctx, s)
	structType := getStructType(s)

	if errVal != nil {
		errObj := castedValidateAll(errVal.(validator.ValidationErrors), structType)
		errStr, _ := json.Marshal(errObj)
		err = fmt.Errorf("%s", errStr)
		slog.ErrorContext(ctx, "error validate [c]", err)
		return
	}
	return
}

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^\\d{4}(-\\d\\d(-\\d\\d(T\\d\\d:\\d\\d(:\\d\\d)?(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?)?)?)?$"
	return regexp.MustCompile(ISO8601DateRegexString).MatchString(fl.Field().String())
}

// sanitizeString performs the actual sanitization logic
func sanitizeString(input string) string {
	input = strings.TrimSpace(input)
	input = template.HTMLEscapeString(input)
	return input
}

func Sanitize(fl validator.FieldLevel) bool {
	if str, ok := fl.Field().Interface().(string); ok {
		sanitizedStr := sanitizeString(str)
		fl.Field().SetString(sanitizedStr)
		return true
	}
	return false
}

func IsValidDate(ctx context.Context, s string) (err error) {
	_, err = time.Parse("02-01-2006", s)
	if err != nil {
		slog.ErrorContext(ctx, "error validate date format", err)
		err = fmt.Errorf("date is not valid format (dd-mm-yyyy)")
	}
	return
}

// func ValidateUUID(ctx context.Context, u string) (err error) {
// 	_, err = uuid.Parse(u)
// 	if err != nil {
// 		slog.ErrorContext(ctx, "failed to parse uuid", err)
// 		err = fmt.Errorf("invalid uuid")
// 		return
// 	}
// 	return
// }

func ValidateMultipartFormValue(ctx context.Context, values map[string][]string, target interface{}) (err error) {
	rawRequest := make(map[string]string, 0)
	for key, vals := range values {
		if len(vals) > 0 {
			rawRequest[key] = vals[0]
		}
	}
	rawByte, err := json.Marshal(rawRequest)
	if err != nil {
		slog.ErrorContext(ctx, "failed to validate multipart form value", err)
		err = fmt.Errorf("something went wrong")
		return
	}
	err = json.Unmarshal(rawByte, target)
	if err != nil {
		slog.ErrorContext(ctx, "failed to validate multipart form value", err)
		err = fmt.Errorf("something went wrong")
		return
	}
	err = ValidateStruct(ctx, target)
	if err != nil {
		slog.ErrorContext(ctx, "failed to validate multipart form value", err)
	}
	return
}

func castedValidate(valErr validator.ValidationErrors) (err error) {
	for _, v := range valErr {
		switch v.Tag() {
		case "required":
			err = fmt.Errorf("%s is required", v.Field())
		case "email":
			err = fmt.Errorf("%s is not a valid email", v.Field())
		case "min":
			err = fmt.Errorf("%s is too short, minimum %s digit", v.Field(), v.Param())
		case "max":
			err = fmt.Errorf("%s is too long maximum %s digit", v.Field(), v.Param())
		case "len":
			err = fmt.Errorf("%s length is not valid", v.Field())
		case "eqfield":
			err = fmt.Errorf("%s is not equal to %s", v.Field(), v.Param())
		case "eq":
			err = fmt.Errorf("%s is not equal to %s", v.Field(), v.Param())
		case "gt":
			err = fmt.Errorf("%s is not greater than %s", v.Field(), v.Param())
		case "gte":
			err = fmt.Errorf("%s is not greater than or equal to %s", v.Field(), v.Param())
		case "lt":
			err = fmt.Errorf("%s is not less than %s", v.Field(), v.Param())
		case "lte":
			err = fmt.Errorf("%s is not less than or equal to %s", v.Field(), v.Param())
		case "ne":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "nfeq":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "oneof":
			err = fmt.Errorf("%s is not one of %s", v.Field(), v.Param())
		case "uuid":
			err = fmt.Errorf("%s is not a valid uuid", v.Field())
		case "ISO8601Date":
			err = fmt.Errorf("%s is not a valid ISO8601Date", v.Field())
		case "nefield":
			err = fmt.Errorf("%s is equal to %s", v.Field(), v.Param())
		case "validInprogressStatus":
			err = fmt.Errorf("field is not valid status")
		default:
			err = fmt.Errorf("%s is not valid", v.Field())
		}
	}
	return

}

type ErrValidation struct {
	Key   string `json:"key"`
	Type  string `json:"types"`
	Error string `json:"error"`
}

func getStructType(s interface{}) (typ reflect.Type) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ = val.Type()
	return typ
}

func unmarshalErrValidationStr(input string) {

}

// getFieldPath constructs the full JSON path for a validation error
func getFieldPathAndType(namespace string, typ reflect.Type) (string, string) {
	parts := strings.Split(namespace, ".")
	if len(parts) < 2 {
		return "", ""
	}

	var jsonTags []string
	var fieldType string
	for _, part := range parts[1:] {
		field, found := typ.FieldByName(part)
		if !found {
			return "", ""
		}

		jsonTag := field.Tag.Get("json")
		// if jsonTag == "" skip, we assume it as extended struct
		if jsonTag != "" {
			jsonTag = strings.Split(jsonTag, ",")[0]
			jsonTags = append(jsonTags, jsonTag)
		}

		typ = field.Type
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Struct {
			// Continue to the next part
		} else {
			fieldType = typ.String()
			break
		}
	}

	return strings.Join(jsonTags, "."), fieldType
}

type SetErrorStruct struct {
	Namespace string
	Data      interface{}
	Message   string
	DataType  string //reflect not working for file type
}

func SetErrorStructValidation(data SetErrorStruct) (err ErrValidation) {
	requestStructType := getStructType(data.Data)
	fieldPath, fieldType := getFieldPathAndType(data.Namespace, requestStructType)
	err = ErrValidation{
		Key:   fieldPath,
		Type:  fieldType,
		Error: "Invalid",
	}
	if data.DataType == "File" {
		err.Type = "File"
	}
	if data.DataType == "integer" || data.DataType == "uint" {
		err.Type = "int"
	}
	if data.DataType == "string" {
		err.Type = "string"
	}
	if data.Message != "" {
		err.Error = data.Message
	}
	return
}

func castedValidateAll(valErr validator.ValidationErrors, structType reflect.Type) (err []ErrValidation) {
	for _, v := range valErr {
		fieldPath, _ := getFieldPathAndType(v.Namespace(), structType)
		typeField := v.Type().String()
		if typeField == "uint" {
			typeField = "int"
		}
		errObj := ErrValidation{
			Key:  fieldPath,
			Type: typeField,
		}
		switch v.Tag() {
		case "required":
			errObj.Error = "Required"
			err = append(err, errObj)
		case "email":
			errObj.Error = "Invalid email"
			err = append(err, errObj)
		case "min":
			errObj.Error = fmt.Sprintf("Minimum %s digit/char", v.Param())
			err = append(err, errObj)
		case "max":
			errObj.Error = fmt.Sprintf("Maximum %s digit/char", v.Param())
			err = append(err, errObj)
		case "len":
			errObj.Error = fmt.Sprintf("Length must be %s digit/char", v.Param())
			err = append(err, errObj)
		case "eqfield":
			errObj.Error = fmt.Sprintf("%s is not equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "eq":
			errObj.Error = fmt.Sprintf("%s is not equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "gt":
			errObj.Error = fmt.Sprintf("%s is not greater than %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "gte":
			errObj.Error = fmt.Sprintf("%s is not greater than or equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "lt":
			errObj.Error = fmt.Sprintf("%s is not less than %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "lte":
			errObj.Error = fmt.Sprintf("%s is not less than or equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "ne":
			errObj.Error = fmt.Sprintf("%s is equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "nfeq":
			errObj.Error = fmt.Sprintf("%s is equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "oneof":
			errObj.Error = fmt.Sprintf("%s is not one of %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "uuid":
			errObj.Error = "Invalid UUID"
			err = append(err, errObj)
		case "ISO8601Date":
			errObj.Error = "Invalid ISO8601Date"
			err = append(err, errObj)
		case "nefield":
			errObj.Error = fmt.Sprintf("%s is equal to %s", v.Field(), v.Param())
			err = append(err, errObj)
		case "validInprogressStatus":
			errObj.Error = "Invalid inprogress status"
			err = append(err, errObj)
		default:
			errObj.Error = "Invalid"
			err = append(err, errObj)
		}
	}
	return

}

func IsNumericOnly(str string) bool {
	regex, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		fmt.Println("error compiling the regular expression:", err)
		return false
	}

	// use the MatchString method to check if the string matches the pattern
	return regex.MatchString(str)
}

type CustomValidation struct {
	ValidationErrs map[string]interface{}
	Fields         map[string]interface{}
	ErrorCodes     map[string]interface{}
}

func NewCustomValidation(fileName string) (*CustomValidation, error) {
	validationErr, err := readCustomValidationFiles(fileName)
	if err != nil {
		return nil, err
	}

	return &CustomValidation{
		ValidationErrs: validationErr["validations"],
		Fields:         validationErr["fields"],
		ErrorCodes:     validationErr["errorCodes"],
	}, nil
}

func readCustomValidationFiles(fileName string) (map[string]map[string]interface{}, error) {
	var validationErr map[string]map[string]interface{}
	file, err := os.ReadFile(fmt.Sprintf("lang/en/%s-validation.json", fileName))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error occurred reading validation file: %s", err))
	}
	err = json.Unmarshal(file, &validationErr)
	if err != nil {
		return nil, err
	}

	return validationErr, nil
}

func CamelCaseToWords(input string) string {
	var words []string
	var currentWord string

	for i, char := range input {
		if i == 0 {
			currentWord += string(unicode.ToUpper(char))
		} else if unicode.IsUpper(char) {
			words = append(words, currentWord)
			currentWord = string(char)
		} else {
			currentWord += string(char)
		}
	}
	words = append(words, currentWord)
	return strings.Join(words, " ")
}

func PascalCase(input string) string {
	return strings.ToUpper(string(input[0])) + input[1:]
}

func ConvertToPascalCase(input string) string {
	var delimiters []rune
	words := strings.FieldsFunc(input, func(r rune) bool {
		if unicode.IsSpace(r) || r == '_' || r == '-' || r == '.' {
			delimiters = append(delimiters, r)
			return true
		}
		return false
	})

	var result strings.Builder
	for i, word := range words {
		if len(word) > 0 {
			word = strings.ToUpper(string(word[0])) + word[1:]
			// Handle special case where "Id" should be "ID"
			word = strings.Replace(word, "Id", "ID", 1)
			result.WriteString(word)
		}
		// Append the corresponding delimiter if it exists
		if i < len(delimiters) {
			result.WriteRune(delimiters[i])
		}
	}

	return result.String()
}

func (t *CustomValidation) TranslateValidationError(err validator.FieldError, errorCode ...string) string {
	namespace := err.Namespace()
	field, has := t.Fields[err.Field()]
	if !has {
		if field, has = t.Fields[namespace].(string); !has {
			field = CamelCaseToWords(err.Field())
		}
	}
	message, has := t.ValidationErrs[err.Tag()]
	if !has {
		message = err.Tag()
	}
	messageStr := message.(string)

	if len(errorCode) > 0 {
		// for temporary only handle 1 args
		key := errorCode[0]
		namespace = strings.Replace(namespace, key+".", "", 1)
		if typeCode, ok := t.ErrorCodes[key].(map[string]interface{}); ok {
			if errCodeStr, ok := typeCode[namespace].(string); ok {
				messageStr += " [" + errCodeStr + "]"
			}

		}
	}

	if err.Param() == "" {
		return fmt.Sprintf(
			messageStr,
			field,
		)

	}

	return fmt.Sprintf(
		messageStr,
		field,
		err.Param(),
	)

}

func (t *CustomValidation) CustomValidate(c echo.Context, s interface{}, errorCode ...string) []string {
	c.Bind(s)
	if err := c.Validate(s); err != nil {
		var validationErrs []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			message := t.TranslateValidationError(fieldError, errorCode...)
			validationErrs = append(validationErrs, message)
		}
		return validationErrs
	}
	return nil
}

func (t *CustomValidation) CustomValidateStruct(ctx context.Context, s interface{}, errorCode ...string) []string {
	var trimSpaces func(reflect.Value)
	trimSpaces = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.String:
			v.SetString(strings.TrimSpace(v.String()))
		case reflect.Ptr:
			if !v.IsNil() {
				trimSpaces(v.Elem())
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				trimSpaces(v.Field(i))
			}
		}
	}
	val := reflect.ValueOf(s).Elem()
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			trimSpaces(val.Field(i))
		}
	}
	if err := v.StructCtx(ctx, s); err != nil {
		var validationErrs []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			message := t.TranslateValidationError(fieldError, errorCode...)
			validationErrs = append(validationErrs, message)
		}
		return validationErrs
	}
	return nil
}

func (t *CustomValidation) CustomErrorMessage(errorKey string, parent string, customVld ...string) string {
	var messageStr string
	if msg, ok := t.Fields[errorKey].(string); ok {
		messageStr = msg
	} else {
		key := parent + "." + errorKey
		if msg, ok := t.Fields[key].(string); ok {
			messageStr = msg
		}
	}

	finalMsg := messageStr
	if len(customVld) > 0 {
		finalMsg = fmt.Sprintf(t.ValidationErrs[customVld[0]].(string), messageStr)
	}
	if typeCode, ok := t.ErrorCodes[parent].(map[string]interface{}); ok {
		if errCodeStr, ok := typeCode[errorKey].(string); ok {
			finalMsg += " [" + errCodeStr + "]"
		}
	}
	return finalMsg
}
