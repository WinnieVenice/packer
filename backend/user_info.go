package backend

import (
	"fmt"

	"github.com/WinnieVenice/packer/model"
)

func AddUser(user *model.UserInfo) error {
	return db.Table(model.TableUserInfo).Create(&user).Error
}

func UpdUser(newUser *model.UserInfo) error {
	var user model.UserInfo
	if err := db.Table(model.TableUserInfo).First(&user, "id=?", newUser.Id).Error; err != nil {
		return err
	}
	user.Convert(newUser)
	return db.Table(model.TableUserInfo).Save(&user).Error
}

func DelUser(user *model.UserInfo) error {
	return db.Table(model.TableUserInfo).Delete(&user).Error
}

func QueUserById(id string) (user *model.UserInfo, err error) {
	err = db.Table(model.TableUserInfo).First(&user, "id=?", id).Error
	return
}

func QueUserByPlatformId(platform, id string) (user *model.UserInfo, err error) {
	sql := fmt.Sprintf("%s_id = %s", platform, id)
	err = db.Table(model.TableUserInfo).Where(sql).Find(&user).Error
	return
}

func GetAllUser() (users []model.UserInfo, err error) {
	err = db.Table(model.TableUserInfo).Find(&users).Error
	return
}
