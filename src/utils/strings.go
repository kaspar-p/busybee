package utils

func StringInSlice(list []string, value string) (isIn bool, foundIndex int) {
	isIn = false

	for index, element := range list {
		if value == element {
			isIn = true
			foundIndex = index

			break
		}
	}

	return isIn, foundIndex
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
