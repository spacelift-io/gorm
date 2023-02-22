package gorm_test

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestDelete(t *testing.T) {
	user1, user2 := User{Name: "delete1"}, User{Name: "delete2"}
	DB.Save(&user1)
	DB.Save(&user2)

	if err := DB.Delete(&user1).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	}

	if !errors.Is(DB.Where("name = ?", user1.Name).First(&User{}).Error, gorm.ErrRecordNotFound) {
		t.Errorf("User can't be found after delete")
	}

	if errors.Is(DB.Where("name = ?", user2.Name).First(&User{}).Error, gorm.ErrRecordNotFound) {
		t.Errorf("Other users that not deleted should be found-able")
	}
}

func TestInlineDelete(t *testing.T) {
	user1, user2 := User{Name: "inline_delete1"}, User{Name: "inline_delete2"}
	DB.Save(&user1)
	DB.Save(&user2)

	if DB.Delete(&User{}, user1.Id).Error != nil {
		t.Errorf("No error should happen when delete a record")
	} else if !errors.Is(DB.Where("name = ?", user1.Name).First(&User{}).Error, gorm.ErrRecordNotFound) {
		t.Errorf("User can't be found after delete")
	}

	if err := DB.Delete(&User{}, "name = ?", user2.Name).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	} else if !errors.Is(DB.Where("name = ?", user2.Name).First(&User{}).Error, gorm.ErrRecordNotFound) {
		t.Errorf("User can't be found after delete")
	}
}

func TestSoftDelete(t *testing.T) {
	type User struct {
		Id        int64
		Name      string
		DeletedAt gorm.DeletedAt
	}
	DB.AutoMigrate(&User{})

	user := User{Name: "soft_delete"}
	DB.Save(&user)
	DB.Delete(&user)

	if DB.First(&User{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	var retrievedUser User
	if err := DB.Unscoped().First(&retrievedUser, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but err=%s", err)
	}

	if !retrievedUser.DeletedAt.Valid || retrievedUser.DeletedAt.Time.IsZero() {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but DeletedAt is not set")
	}

	DB.Unscoped().Delete(&user)
	if !errors.Is(DB.Unscoped().First(&User{}, "name = ?", user.Name).Error, gorm.ErrRecordNotFound) {
		t.Errorf("Can't find permanently deleted record")
	}
}

func TestSoftDelete_UsingTime_InsteadOf_GormTime_Is_Ignored(t *testing.T) {
	creditCard := CreditCard{Number: "411111111234567"}
	DB.Save(&creditCard)
	DB.Delete(&creditCard)

	var retrieved CreditCard
	result := DB.Unscoped().First(&retrieved, "number = ?", creditCard.Number)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Errorf("Must be entirely wiped from the database. Expected ErrRecordNotFound, got %v", result.Error)
	}
}
