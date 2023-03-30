package backend

import (
	"github.com/WinnieVenice/packer/model"
)

func AddPlatformUser(user *model.PlatformUserInfo) error {
	return db.Table(model.TablePlatformUserInfo).Create(&user).Error
}

func UpdPlatformUser(newUser *model.PlatformUserInfo) error {
	var user model.PlatformUserInfo
	if err := db.Table(model.TablePlatformUserInfo).First(&user, "id=?", newUser.Id).Error; err != nil {
		return err
	}
	user.Convert(newUser)
	return db.Table(model.TablePlatformUserInfo).Save(&user).Error
}

func DelPlatformUser(user *model.PlatformUserInfo) error {
	return db.Table(model.TablePlatformUserInfo).Delete(&user).Error
}

func QuePlatformUserById(id string) (user *model.PlatformUserInfo, err error) {
	err = db.Table(model.TablePlatformUserInfo).First(&user, "id=?", id).Error
	return
}

func GetAllPlatformUserId() (ids []string, err error) {
	err = db.Table(model.TablePlatformUserInfo).Pluck("id", &ids).Error
	return
}

func GetAllPlatformUser() (users []model.PlatformUserInfo, err error) {
	err = db.Table(model.TablePlatformUserInfo).Find(&users).Error
	return
}
