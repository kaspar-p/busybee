package utils

func StringInSlice(array []string, value string) (bool, int) {
	for index, element := range array {
		if value == element {
			return true, index
		}
	}

	return false, 0
}

func RemoveStringFromSlice(array []string, value string) []string {
	found, index := StringInSlice(array, value)
	if !found {
		return array
	} else {
		newArray := make([]string, len(array)-1)
		newArray = append(newArray, array[:index]...)
		newArray = append(newArray, array[index+1:]...)
		return newArray;
	}
}