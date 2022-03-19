package main

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/xorm"
)

type userID int64

var (
	DBdriver = "mysql"
	DBinfo   = "root:password@/chatServer?charset=utf8"
	DBengine *xorm.Engine
)

type User struct {
	ID       userID
	Name     string
	Salt     string
	Password string `xorm:"varchar(200)"`
}

type Message struct {
	fromUserID userID
	toUserID   userID
	Content    string
	Created    time.Time
	Modified   time.Time
}

type Group struct {
	groupID   int64
	leader    userID
	adminList []userID
	userList  []userID
	count     int
}

func InitDB() {
	DBengine, err := xorm.NewEngine(DBdriver, DBinfo)
	// xorm.NewEngine("postgres", "user=postgres password=123 dbname=test sslmode=disable")
	// xorm.NewEngine("sqlite3", "./test.db")
	// xorm.NewEngine("mssql", "sqlserver://sa:123@localhost:1433?database=test")
	if err != nil {
		panic(err)
	} else {
		// DBengine.ShowSQL(true)

		// check if table exists, will create it if not
		DBengine.Sync2(new(User), new(Message), new(Group))
		// err = DBengine.Sync2(new(User), new(Message), new(Group))
		// if err != nil {
		// 	panic(err)
		// }
	}
}

func (s *Group) SetAdmin(id userID) error {
	if has, _ := DBengine.Exist(&Group{
		groupID:  s.groupID,
		userList: []userID{id},
	}); has {
		s.adminList = append(s.adminList, id)
		RemoveSlice(s.userList, id)
		DBengine.Update(&Group{}, s)
		// e.Table(&Group{}).Update(s)
	}
	return nil
}

func (u *User) CheckExist() bool {
	has, _ := DBengine.Exist(u)
	return has
}

func (u *User) Add() error {
	u.Salt = GetRandomString(10)
	u.Password = GetMD5(u.Password + u.Salt)
	_, err := DBengine.Insert(u)
	return err
}

func (u *User) SetPassword(pwd string) error {
	u.Salt = GetRandomString(10)
	u.Password = GetMD5(pwd + u.Salt)
	_, err := DBengine.Update(u)
	return err
}

// 2 - false, err - user not found
// 1 - false, nil - password not match
// 0 - true, nil  - login success
func (e *User) Login() (bool, error) {
	var user User
	has, err := DBengine.Where("id = ?", e.ID).Get(&user)
	if err != nil {
		return false, err
	}
	if has {
		if user.Password == GetMD5(e.Password+user.Salt) {
			return true, nil
		}
	}
	return false, nil
}
