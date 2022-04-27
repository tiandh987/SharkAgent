package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	DNS1123SubdomainMaxLength int = 253
	dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	dns1123SubdomainFmt string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"
	dns1123SubdomainErrorMsg string = "a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and   must start and end with an alphanumeric character"
)

const (
	qualifiedNameErrMsg    string = "must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
	qualifiedNameMaxLength int    = 63
)

const (
	qnameCharFmt     string = "[A-Za-z0-9]"
	qnameExtCharFmt  string = "[-A-Za-z0-9_.]"
	qualifiedNameFmt string = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
)

var (
	dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")
)

var (
	qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")
)

const (
	minPassLength = 8
	maxPassLength = 16
)

// IsQualifiedName 测试传递的值是否是 IAM 所说的“qualified name”。
// 这是在整个系统的各个地方使用的格式。 如果该值无效，则返回错误字符串列表。
// 否则返回一个空列表（或 nil）。
func IsQualifiedName(value string) []string {
	var errs []string

	// 以 / 切分value
	parts := strings.Split(value, "/")

	var name string

	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			errs = append(errs, "prefix part" + EmptyError())
		} else if msgs := IsDNS1123Subdomain(prefix); len(msgs) != 0 {
			errs = append(errs, prefixEach(msgs, "prefix part ")...)
		}
	default:
		return append(
			errs,
			"a qualified name "+ RegexError(
				qualifiedNameErrMsg,
				qualifiedNameFmt,
				"MyName",
				"my.name",
				"123-abc",
			) + " with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')",
		)
	}

	if len(name) == 0 {
		errs = append(errs, "name part "+EmptyError())
	} else if len(name) > qualifiedNameMaxLength {
		errs = append(errs, "name part "+MaxLenError(qualifiedNameMaxLength))
	}

	if !qualifiedNameRegexp.MatchString(name) {
		errs = append(
			errs,
			"name part "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"),
		)
	}
	return errs
}

// EmptyError returns a string explanation of a "must not be empty" validation
// failure.
func EmptyError() string {
	return "must be non-empty"
}

// MaxLenError returns a string explanation of a "string too long" validation
// failure.
func MaxLenError(length int) string {
	return fmt.Sprintf("must be no more than %d characters", length)
}

// RegexError returns a string explanation of a regex validation failure.
func RegexError(msg string, fmt string, examples ...string) string {
	if len(examples) == 0 {
		return msg + " (regex used for validation is '" + fmt + "')"
	}
	msg += " (e.g. "
	for i := range examples {
		if i > 0 {
			msg += " or "
		}
		msg += "'" + examples[i] + "', "
	}
	msg += "regex used for validation is '" + fmt + "')"
	return msg
}

// IsDNS1123Subdomain 测试符合 DNS (RFC 1123) 中子域定义的字符串。
func IsDNS1123Subdomain(value string) []string {
	var errs []string
	if len(value) > DNS1123SubdomainMaxLength {
		errs = append(errs, MaxLenError(DNS1123SubdomainMaxLength))
	}
	if !dns1123SubdomainRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123SubdomainErrorMsg, dns1123SubdomainFmt, "example.com"))
	}
	return errs
}

func prefixEach(msgs []string, prefix string) []string {
	for i := range msgs {
		msgs[i] = prefix + msgs[i]
	}
	return msgs
}

// IsValidPassword validate password.
func IsValidPassword(password string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
			passLen++
		case unicode.IsUpper(ch):
			hasUpper = true
			passLen++
		case unicode.IsLower(ch):
			hasLower = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !hasLower {
		appendError("lowercase letter missing")
	}
	if !hasUpper {
		appendError("uppercase letter missing")
	}
	if !hasNumber {
		appendError("at least one numeric character required")
	}
	if !hasSpecial {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(
			fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength),
		)
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}

	return nil
}