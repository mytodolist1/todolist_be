package component

import "regexp"

func ValidatePhoneNumber(phoneNumber string) (bool, error) {
	numericPattern := `^[0-9]+$`
	numericRegexp, err := regexp.Compile(numericPattern)
	if err != nil {
		return false, err
	}
	if !numericRegexp.MatchString(phoneNumber) {
		return false, nil
	}

	pattern := `^62\d{6,13}$`
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	isValid := regexp.MatchString(phoneNumber)

	return isValid, nil
}