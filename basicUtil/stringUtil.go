package basicUtil

func HandleString(callStationId string) string {
	size := len(callStationId)
	if size%2 == 1 {
		return ""
	}
	result := ""
	for i := 0; i < size/2; i++ {
		result += callStationId[2*i : 2*i+2]
		if i != size/2-1 {
			result += "-"
		}
	}
	return result
}
