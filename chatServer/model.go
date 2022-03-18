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

func initDB() *xorm.Engine {
	engine, err := xorm.NewEngine(DBdriver, DBinfo)
	// xorm.NewEngine("postgres", "user=postgres password=123 dbname=test sslmode=disable")
	// xorm.NewEngine("sqlite3", "./test.db")
	// xorm.NewEngine("mssql", "sqlserver://sa:123@localhost:1433?database=test")
	if err != nil {
		panic(err)
	} else {
		err = engine.Sync2(new(User), new(Message), new(Group))
		if err != nil {
			panic(err)
		} else {
			return engine
		}
	}
}

func (s *Group) setAdmin(e *xorm.Engine, id userID) error {
	if has, _ := e.Exist(&Group{
		groupID:  s.groupID,
		userList: []userID{id},
	}); has {
		s.adminList = append(s.adminList, id)
		RemoveSlice(s.userList, id)
		e.Update(&Group{}, s)
		// e.Table(&Group{}).Update(s)
	}
	return nil
}

func RemoveSlice(slice []userID, i userID) []userID {
	for j, v := range slice {
		if v == i {
			slice = append(slice[:j], slice[j+1:]...)
			break
		}
	}
	return slice
}

func (u *User) Add(e *xorm.Engine) error {
	u.Salt = GetRandomString(10)
	u.Password = GetMD5(u.Password + u.Salt)
	_, err := e.Insert(u)
	return err
}

func (u *User) setPassword(e *xorm.Engine, pwd string) error {
	u.Salt = GetRandomString(10)
	u.Password = GetMD5(pwd + u.Salt)
	_, err := e.Update(u)
	return err
}
