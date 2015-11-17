package service

import (
	"cache"
)

var overflow bool = true
var overflow_exceed bool = false
var overflow_thresh int = 404

func Count() int {
	cnt := cache.Count()
	if cnt > overflow_thresh {
		overflow_exceed = true
	}
	return cnt
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
