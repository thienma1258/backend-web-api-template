package api

import "strings"

type LOG_MODE string

const LOG_DEBUG_MODE = "DEBUG"
const WARNING_DEBUG_MODE = "WARNING"
const INFO_DEBUG_MODE = "INFO"

func (mode *LOG_MODE) isDebugMode() bool {
	if mode == nil || len(*mode) == 0 {
		return false
	}
	if strings.ToUpper(string(*mode)) == LOG_DEBUG_MODE {
		return true
	}
	return false
}

func (mode *LOG_MODE) isWarningMode() bool {
	if mode == nil || len(*mode) == 0 {
		return false
	}
	if strings.ToUpper(string(*mode)) == WARNING_DEBUG_MODE {
		return true
	}
	return false
}

func (mode *LOG_MODE) isInfoMode() bool {
	if mode == nil || len(*mode) == 0 {
		return false
	}
	if strings.ToUpper(string(*mode)) == INFO_DEBUG_MODE {
		return true
	}
	return false
}
