package common

func idFltr(ids []Id, idV Id) bool {
	if len(ids) == 0 {
		return true
	} else {
		for _, v := range ids {
			if v == idV {
				return true
			}
		}
	}
	return false
}

func flagFltr(fltrFlag, measFlag Flag) bool {
	return fltrFlag == 0 || fltrFlag == measFlag
}

func inTimeInterval(from, to, tstamp Time) bool {
	// [from tstamp to]
	if from <= tstamp && to >= tstamp {
		return true
	} else {
		return false
	}
}
