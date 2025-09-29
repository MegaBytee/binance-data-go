package types

var TimeFrames = []string{"1s", "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w"}

func flagFrom(value string) string {

	switch value {
	case "s":
		return Spot
	case "f":
		return Futures
	default:
		return Spot
	}

}

func flagDateTime(value string) string {
	switch value {
	case "d":
		return Daily
	case "m":
		return Monthly
	default:
		return Monthly
	}
}

func IsTimeFrameValidChoice(value string) bool {
	for _, c := range TimeFrames {
		if value == c {
			return true
		}
	}
	return false
}
