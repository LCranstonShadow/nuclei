package expressions

import (
	"regexp"

	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/generators"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/replacer"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/starlight"
)

var templateExpressionRegex = regexp.MustCompile(`(?m)\{\{[^}]+\}\}["'\)\}]*`)

// Evaluate checks if the match contains a dynamic variable, for each
// found one we will check if it's an expression and can
// be compiled, it will be evaluated and the results will be returned.
//
// The provided keys from finalValues will be used as variable names
// for substitution inside the expression.
func Evaluate(data string, base map[string]interface{}) (string, error) {
	data = replacer.Replace(data, base)

	dynamicValues := make(map[string]interface{})
	for _, match := range templateExpressionRegex.FindAllString(data, -1) {
		expr := generators.TrimDelimiters(match)
		result, err := starlight.Eval(expr, base)
		if err != nil {
			continue
		}
		dynamicValues[expr] = result
	}
	// Replacer dynamic values if any in raw request and parse  it
	return replacer.Replace(data, dynamicValues), nil
}
