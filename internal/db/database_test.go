package db

import (
	"github.com/Dendyator/AntiBF/internal/logger"
	"testing"
)

func TestRedisOperations(t *testing.T) {
	log := logger.New("debug")
	InitRedis("localhost:6379", log)
	defer CloseRedis()

	if !UpdateRedis("whitelist", "192.168.1.1/25", true) {
		t.Errorf("Failed to add to whitelist")
	}

	inList, err := CheckInRedis("whitelist", "192.168.1.1/25")
	if err != nil {
		t.Errorf("Error checking whitelist: %v", err)
	} else if !inList {
		t.Errorf("IP not found in whitelist")
	}

	if !UpdateRedis("whitelist", "192.168.1.1/25", false) {
		t.Errorf("Failed to remove from whitelist")
	}

	inList, err = CheckInRedis("whitelist", "192.168.1.1/25")
	if err != nil {
		t.Errorf("Error checking whitelist: %v", err)
	} else if inList {
		t.Errorf("IP should not be in whitelist")
	}
}
