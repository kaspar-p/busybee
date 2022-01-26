package utils

func StringInSlice(list []string, value string) (bool, int) {
	result := false
	finalIndex := 0

	for index, element := range list {
		if value == element {
			result = true
			finalIndex = index

			break
		}
	}

	return result, finalIndex
}

func RemoveStringFromSlice(array []string, value string) []string {
	found, index := StringInSlice(array, value)
	if !found {
		return array
	}

	newArray := make([]string, 0)
	newArray = append(newArray, array[:index]...)
	newArray = append(newArray, array[index+1:]...)

	return newArray
}

func WrapStringInCodeBlock(s string) string {
	return "```" + s + "```"
}
