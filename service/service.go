package service

import (
	"cache"
)

const overflow bool = false
const overflow_thresh int = 404

var overflow_exceed bool = false

func Count() int {
	if overflow {
		cnt := cache.Count()
		if cnt > overflow_thresh {
			overflow_exceed = true
		}
		return cnt
	} else {
		return 0
	}
}

func OverFlow() bool {
	if overflow {
		if overflow_exceed {
			return true
		} else {
			cnt := cache.GetCount()
			if cnt > overflow_thresh {
				overflow_exceed = true
			}
			return overflow_exceed
		}

	} else {
		return false
	}
}
