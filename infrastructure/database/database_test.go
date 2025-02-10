package database

import (
	"testing"

	"github.com/Dendyator/AntiBF/pkg/logger"
)

func TestRedisOperations(t *testing.T) {
	log := logger.New("debug")

	// Создаем экземпляр DB
	db := NewDB("postgres://user:password@localhost:5432/antibruteforce", "localhost:6379", log)
	defer db.Close()

	// Тест добавления в whitelist
	if !db.UpdateRedis("whitelist", "192.168.1.1/25", true) {
		t.Errorf("Failed to add to whitelist")
	}

	// Проверка наличия в whitelist
	inList, err := db.CheckInRedis("whitelist", "192.168.1.1/25")
	if err != nil {
		t.Errorf("Error checking whitelist: %v", err)
	} else if !inList {
		t.Errorf("IP not found in whitelist")
	}

	// Тест удаления из whitelist
	if !db.UpdateRedis("whitelist", "192.168.1.1/25", false) {
		t.Errorf("Failed to remove from whitelist")
	}

	// Проверка отсутствия в whitelist
	inList, err = db.CheckInRedis("whitelist", "192.168.1.1/25")
	if err != nil {
		t.Errorf("Error checking whitelist: %v", err)
	} else if inList {
		t.Errorf("IP should not be in whitelist")
	}
}
