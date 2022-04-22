package git_release

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s", value, k))
		}
		return
	}
}

// return the pieces of id `left:right` as left, right
func parseTwoPartID(id, left, right string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Unexpected ID format (%q). Expected %s:%s", id, left, right)
	}

	return parts[0], parts[1], nil
}

// format the strings into an id `a:b`
func buildTwoPartID(a, b string) string {
	return fmt.Sprintf("%s:%s", a, b)
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func flattenStringList(v []string) []interface{} {
	c := make([]interface{}, 0, len(v))
	for _, s := range v {
		c = append(c, s)
	}
	return c
}

func unconvertibleIdErr(id string, err error) *unconvertibleIdError {
	return &unconvertibleIdError{OriginalId: id, OriginalError: err}
}

type unconvertibleIdError struct {
	OriginalId    string
	OriginalError error
}

func (e *unconvertibleIdError) Error() string {
	return fmt.Sprintf("Unexpected ID format (%q), expected numerical ID. %s",
		e.OriginalId, e.OriginalError.Error())
}

func validateTeamIDFunc(v interface{}, keyName string) (we []string, errors []error) {
	teamIDString, ok := v.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", keyName)}
	}
	// Check that the team ID can be converted to an int
	if _, err := strconv.ParseInt(teamIDString, 10, 64); err != nil {
		return nil, []error{unconvertibleIdErr(teamIDString, err)}
	}

	return
}

// https://docs.github.com/en/actions/reference/encrypted-secrets#naming-your-secrets
var secretNameRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

func validateSecretNameFunc(v interface{}, keyName string) (we []string, errs []error) {
	name, ok := v.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", keyName)}
	}

	if !secretNameRegexp.MatchString(name) {
		errs = append(errs, errors.New("Secret names can only contain alphanumeric characters or underscores and must not start with a number"))
	}

	if strings.HasPrefix(strings.ToUpper(name), "GITHUB_") {
		errs = append(errs, errors.New("Secret names must not start with the GITHUB_ prefix"))
	}

	return we, errs
}
