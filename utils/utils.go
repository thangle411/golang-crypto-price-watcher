package utils

func ShouldNotify(last float64, current float64, bottomRange float64, topRange float64, increment float64, percentThreshold float64) bool {
	if last == 0 {
		return false
	}
	thresholds := []float64{}
	for i := bottomRange; i < topRange+1; i = i + increment {
		thresholds = append(thresholds, i)
	}
	for _, threshold := range thresholds {
		if last < threshold && current > threshold {
			return true
		}
	}
	percent := (current - last) / last * 100
	if percent > percentThreshold {
		return true
	}
	return false
}
