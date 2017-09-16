package validator

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator func(name string, value string, inputs map[string]string, params []string) bool

var validators = map[string]Validator{
	"accepted":       validateAccepted,
	"active_url":     validateActiveUrl,
	"alpha":          validateAlpha,
	"alpha_dash":     validateAlphaDash,
	"alpha_num":      validateAlphaNumeric,
	"boolean":        validateBoolean,
	"chars":          validateChars,
	"chars_between":  validateCharsBetween,
	"confirmed":      validateConfirmed,
	"date":           validateDate,
	"different":      validateDifferent,
	"digits":         validateDigits,
	"digits_between": validateDigitsBetween,
	"email":          validateEmail,
	"in":             validateIn,
	"integer":        validateInteger,
	"ip":             validateIp,
	"max_chars":      validateMaxChars,
	"max_digits":     validateMaxDigits,
	"max_value":      validateMaxValue,
	"min_chars":      validateMinChars,
	"min_digits":     validateMinDigits,
	"min_value":      validateMinValue,
	"not_in":         validateNotIn,
	"numeric":        validateNumeric,
	"regex":          validateRegex,
	"required":       validateRequired,
	"same":           validateSame,
	"url":            validateUrl,
	"value":          validateValue,
	"value_between":  validateValueBetween,
}

var messages = map[string]string{
	"accepted":       "The %s must be accepted.",
	"active_url":     "The %s is not a valid URL.",
	"alpha":          "The %s may only contain letters.",
	"alpha_dash":     "The %s may only contain letters, numbers, and dashes.",
	"alpha_num":      "The %s may only contain letters and numbers.",
	"boolean":        "The %s field must be true or false.",
	"chars":          "The %s field must have %s characters.",
	"chars_between":  "The %s field must have between %s characters.",
	"confirmed":      "The %s confirmation does not match.",
	"date":           "The %s is not a valid date.",
	"different":      "The %s and %s must be different.",
	"digits":         "The %s must have %s digits.",
	"digits_between": "The %s must have between %s and %s digits.",
	"email":          "The %s must be a valid email address.",
	"in":             "The selected %s is invalid.",
	"max_chars":      "The %s must have fewer than %s characters.",
	"max_digits":     "The %s must have fewer than %s digits.",
	"max_value":      "The %s must be less than %s.",
	"min_chars":      "The %s must have more than %s characters.",
	"min_digits":     "The %s must have more than %s digits.",
	"min_value":      "The %s must be greater than %s.",
	"integer":        "The %s must be an integer.",
	"ip":             "The %s must be a valid IP address.",
	"not_in":         "The selected %s is invalid.",
	"numeric":        "The %s must be a number.",
	"regex":          "The %s format is invalid.",
	"required":       "The %s field is required.",
	"same":           "The %s and %s must match.",
	"url":            "The %s format is invalid.",
	"value":          "The %s must %s.",
	"value_between":  "The %s must be between %s and %s.",
}

func AddValidator(name string, fn Validator, message string) {
	validators[name] = fn
	messages[name] = message
}

func validateAccepted(name string, value string, inputs map[string]string, params []string) bool {
	valid := []string{"1", "true", "yes", "on"}
	return stringInSlice(value, valid)
}

func validateActiveUrl(name string, value string, inputs map[string]string, params []string) bool {
	lc := strings.ToLower(value)
	if validScheme := strings.HasPrefix(lc, "http://") || strings.HasPrefix(lc, "https://"); !validScheme {
		return false
	}
	// trim schema and then check dns
	lc = strings.TrimPrefix(lc, "http://")
	lc = strings.TrimPrefix(lc, "https://")
	_, err := net.LookupHost(lc)
	return err == nil
}

func validateAlpha(name string, value string, inputs map[string]string, params []string) bool {
	return regexp.MustCompile("^[a-zA-Z]+$").MatchString(value)
}

func validateAlphaDash(name string, value string, inputs map[string]string, params []string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9-_]+$").MatchString(value)
}

func validateAlphaNumeric(name string, value string, inputs map[string]string, params []string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(value)
}

func validateBoolean(name string, value string, inputs map[string]string, params []string) bool {
	_, err := strconv.ParseBool(value)
	return err == nil
}

func validateChars(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		charCount, err := strconv.ParseInt(params[0], 10, 16)
		return err == nil && int64(len(value)) == charCount
	}
	return false
}

func validateCharsBetween(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 2 {
		p1 := []string{params[0]}
		p2 := []string{params[1]}
		return validateMinChars(name, value, inputs, p1) && validateMaxChars(name, value, inputs, p2)
	}
	return false
}

func validateConfirmed(name string, value string, inputs map[string]string, params []string) bool {
	fieldValue, fieldExists := inputs[name+"_confirmation"]
	return fieldExists && fieldValue == value
}

func validateDate(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		_, err := time.Parse(params[0], value)
		return err == nil
	}
	return false
}

func validateDifferent(name string, value string, inputs map[string]string, params []string) bool {
	return !validateSame(name, value, inputs, params)
}

func validateDigits(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 && regexp.MustCompile("^[0-9]+$").MatchString(value) {
		digitCount, err := strconv.ParseInt(params[0], 10, 16)
		return err == nil && int64(len(value)) == digitCount
	}
	return false
}

func validateDigitsBetween(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 2 {
		p1 := []string{params[0]}
		p2 := []string{params[1]}
		return validateMinDigits(name, value, inputs, p1) && validateMaxDigits(name, value, inputs, p2)
	}
	return false
}

func validateEmail(name string, value string, inputs map[string]string, params []string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$").MatchString(value)
}

func validateIn(name string, value string, inputs map[string]string, params []string) bool {
	return stringInSlice(value, params)
}

func validateInteger(name string, value string, inputs map[string]string, params []string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

func validateIp(name string, value string, inputs map[string]string, params []string) bool {
	return net.ParseIP(value) != nil
}

func validateMinChars(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		minChars, err := strconv.ParseInt(params[1], 10, 16)
		return err == nil && int64(len(value)) >= minChars
	}
	return false
}

func validateMinDigits(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 && regexp.MustCompile("^[0-9]+$").MatchString(value) {
		minDigits, err := strconv.ParseInt(params[0], 10, 16)
		return err == nil && int64(len(value)) >= minDigits
	}
	return false
}

func validateMinValue(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		minValue, mvErr := strconv.ParseFloat(params[0], 64)
		floatValue, fvErr := strconv.ParseFloat(value, 64)
		return mvErr == nil && fvErr == nil && floatValue >= minValue
	}
	return false
}

func validateMaxChars(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		maxChars, err := strconv.ParseInt(params[1], 10, 16)
		return err == nil && int64(len(value)) <= maxChars
	}
	return false
}

func validateMaxDigits(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 && regexp.MustCompile("^[0-9]+$").MatchString(value) {
		maxDigits, err := strconv.ParseInt(params[0], 10, 16)
		return err == nil && int64(len(value)) <= maxDigits
	}
	return false
}

func validateMaxValue(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		maxValue, mvErr := strconv.ParseFloat(params[0], 64)
		floatValue, fvErr := strconv.ParseFloat(value, 64)
		return mvErr == nil && fvErr == nil && floatValue <= maxValue
	}
	return false
}

func validateNotIn(name string, value string, inputs map[string]string, params []string) bool {
	return !stringInSlice(value, params)
}

func validateNumeric(name string, value string, inputs map[string]string, params []string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func validateRegex(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		rx, err := regexp.Compile(params[0])
		return err == nil && rx.MatchString(value)
	}
	return false
}

func validateRequired(name string, value string, inputs map[string]string, params []string) bool {
	return value != ""
}

func validateSame(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		fieldValue, fieldExists := inputs[params[0]]
		return fieldExists && fieldValue == value
	}
	return false
}

func validateUrl(name string, value string, inputs map[string]string, params []string) bool {
	lc := strings.ToLower(value)
	if validScheme := strings.HasPrefix(lc, "http://") || strings.HasPrefix(lc, "https://"); validScheme {
		return true // todo
	}
	return false
}

func validateValue(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 1 {
		expectedValue, evErr := strconv.ParseFloat(params[0], 16)
		actualValue, avErr := strconv.ParseFloat(value, 16)
		return evErr == nil && avErr == nil && expectedValue == actualValue
	}
	return false
}

func validateValueBetween(name string, value string, inputs map[string]string, params []string) bool {
	if len(params) == 2 {
		p1 := []string{params[0]}
		p2 := []string{params[1]}
		return validateMinValue(name, value, inputs, p1) && validateMaxValue(name, value, inputs, p2)
	}
	return false
}

func Validate(inputs map[string]string, rules map[string]string) (bool, map[string]string) {
	// initialize an error messages map
	messages := make(map[string]string)
	for fieldName, fieldRulesRaw := range rules {
		// process each rule
		// start by extracting relevant field info
		fieldValue, fieldExists := inputs[fieldName]
		fieldRules := strings.Split(fieldRulesRaw, "|")
		fieldIsRequired := stringInSlice("required", fieldRules)
		alwaysValidate := stringInSlice("always", fieldRules)
		if fieldIsRequired && !fieldExists {
			// add message saying field is required
			// don't worry about the value, that will be handled below
			messages[fieldName] = buildErrorMessage(fieldName, "required", []string{})
		} else if fieldIsRequired || alwaysValidate || (!fieldIsRequired && fieldValue != "") {
			// process the rules
			for i := 0; i < len(fieldRules); i++ {
				rule, params := splitRuleParams(fieldRules[i])
				if vFn, vExists := validators[rule]; vExists {
					// specified validator exists, call it
					if !vFn(fieldName, fieldValue, inputs, params) {
						// validation failed, add a message
						messages[fieldName] = buildErrorMessage(fieldName, rule, params)
					}
				}
			}
		}
	}
	// validation succeeded if there are no error messages
	return len(messages) == 0, messages
}

func RulesFromStruct(s interface{}) map[string]string {
	var rules = make(map[string]string)
	rv := reflect.ValueOf(s)
	for i := 0; i < rv.NumField(); i++ {
		rti := rv.Type().Field(i)
		jsonName := rti.Tag.Get("json")
		validate := rti.Tag.Get("validate")
		if jsonName != "" && validate != "" {
			rules[jsonName] = validate
		}
	}
	return rules
}

func buildErrorMessage(fieldName string, rule string, params []string) string {
	message, exists := messages[rule]
	if !exists || strings.Count(message, "%s") != 1+len(params) {
		return fmt.Sprintf("The %s is invalid.", fieldName)
	}
	args := make([]interface{}, 1+len(params))
	args[0] = fieldName
	for i := 0; i < len(params); i++ {
		args[i+1] = params[i]
	}
	return fmt.Sprintf(message, args...)
}

func splitRuleParams(ruleWithParams string) (string, []string) {
	var params []string
	if split := strings.Split(ruleWithParams, ":"); len(split) == 2 {
		// params were passed
		rule := split[0]
		params := strings.Split(split[1], ",")
		return rule, params
	}
	return ruleWithParams, params
}

func stringInSlice(needle string, haystack []string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}
