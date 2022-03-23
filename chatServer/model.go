package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"xorm.io/xorm"
)

type userID int64
type groupID int64

var (
	DBdriver = "mysql"
	DBinfo   = "augists:password@/chatServer?charset=utf8"
	DBengine *xorm.Engine
)

type User struct {
	ID       userID
	Name     string
	Salt     string
	Password string `xorm:"varchar(200)"`
}

type Group struct {
	groupID groupID
	Name    string
	leader  userID
	// store list in mysql
	adminList []userID
	userList  []userID
	count     int
}

type Message struct {
	SenderID   userID
	ReceiverID userID
	Content    string
	Created    time.Time
	Modified   time.Time
}

type GroupMessage struct {
	GroupID  groupID
	UserID   userID
	Content  string
	Created  time.Time
	Modified time.Time
}

func InitDB() {
	DBengine, err := xorm.NewEngine(DBdriver, DBinfo)
	// xorm.NewEngine("postgres", "user=postgres password=123 dbname=test sslmode=disable")
	// xorm.NewEngine("sqlite3", "./test.db")
	// xorm.NewEngine("mssql", "sqlserver://sa:123@localhost:1433?database=test")
	if err != nil {
		panic(err)
	} else {
		DBengine.ShowSQL(true)

		// check if table exists, will create it if not
		// DBengine.Sync2(new(User), new(Message), new(Group))
		log.Fatal(DBengine.Sync2(new(User), new(Message), new(Group)))
		fmt.Println("Database initialized")
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

// false, err - user not found
// false, nil - password not match
// true, nil  - login success
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
