package model

import "time"

// Employee 员共模型

type Employee struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;comment:主键" json:"id"`
	Name       string    `gorm:"type:varchar(32);not null;comment:姓名" json:"name"`
	Username   string    `gorm:"type:varchar(32);not null;uniqueIndex:idx_username;comment:用户名" json:"username"`
	Password   string    `gorm:"type:varchar(64);not null;comment:密码" json:"password"`
	Phone      string    `gorm:"type:varchar(11);not null;comment:手机号" json:"phone"`
	Sex        string    `gorm:"type:varchar(2);not null;comment:性别" json:"sex"`
	IDNumber   string    `gorm:"type:varchar(18);not null;comment:身份证号" json:"id_number"`
	Status     int       `gorm:"not null;default:1;comment:状态 0:禁用，1:启用" json:"status"`
	CreateTime time.Time `gorm:"type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"type:datetime;comment:更新时间" json:"update_time"`
	CreateUser int64     `gorm:"type:bigint;default:null;comment:创建人" json:"create_user"`
	UpdateUser int64     `gorm:"type:bigint;default:null;comment:修改人" json:"update_user"`
}

func (m *Employee) TableName() string {
	return "employee"
}
