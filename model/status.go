package model

import "math/rand"

var statusList = []string{"PENDING", "ACTIVE", "VALIDATION_REQUIRED", "DELETED", "UNKNOWN"}

func randStatus() string {
	return statusList[rand.Intn(len(statusList)-1)]
}
