package main

// This is to do Complex Type Assertions like
// []interface{} to []string
// interface{} to []string

func convertInterfaceToStringArray(rawValue interface{}) ([]string, bool) {
	rawArrayOfValues, arrayOk := rawValue.([]interface{})
	if !arrayOk {
		return nil, false
	}

	return convertInterfaceArrayToStringArray(rawArrayOfValues)
}

func convertInterfaceArrayToStringArray(rawArrayOfValues []interface{}) ([]string, bool) {
	var stringArray []string
	for _, rawStringValue := range rawArrayOfValues {
		stringValue, stringOk := rawStringValue.(string)
		if !stringOk {
			return nil, false
		}

		stringArray = append(stringArray, stringValue)
	}

	return stringArray, true
}
