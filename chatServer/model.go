package main

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"xorm.io/xorm"
)

var (
	DBdriver = "mysql"
	DBinfo   = "augists:A1b2c3d4,.@/chatServer?charset=utf8"
	DBengine *xorm.Engine
)

type User struct {
	ID       string
	Name     string
	Salt     string
	Password string `xorm:"varchar(200)"`
}

type Group struct {
	GroupID string
	Name    string
	Leader  string
	// store list in mysql
	AdminList []string `xorm:"varchar(200)"`
	UserList  []string `xorm:"varchar(200)"`
	Count     int
}

type Message struct {
	SenderID   string
	ReceiverID string
	Content    string
	Created    time.Time
	Modified   time.Time
}

type GroupMessage struct {
	GroupID  string
	UserID   string
	Content  string
	Created  time.Time
	Modified time.Time
}

func InitDB() {
	var err error
	DBengine, err = xorm.NewEngine(DBdriver, DBinfo)
	// xorm.NewEngine("postgres", "user=postgres password=123 dbname=test sslmode=disable")
	// xorm.NewEngine("sqlite3", "./test.db")
	// xorm.NewEngine("mssql", "sqlserver://sa:123@localhost:1433?database=test")
	if err != nil {
		panic(err)
	} else {
		DBengine.ShowSQL(true)

		// check if table exists, will create it if not
		DBengine.Sync2(new(User), new(Message), new(Group))
		log.Print("Database initialized")
		// err = DBengine.Sync2(new(User), new(Message), new(Group))
		// if err != nil {
		// 	panic(err)
		// }
	}
}

func (u *User) Login() (bool, error) {
	// log.Print(*u)
	// has, err := DBengine.Where("i_d = ?", e.ID).Get(&user)
	if u.CheckExist() {
		if u.Password == GetMD5(u.Password+u.Salt) {
			return true, nil
		} else {
			return false, errors.New("password is wrong")
		}
	} else {
		return false, errors.New("user does not exist")
	}
}

func (u *User) CheckExist() bool {
	// log.Print(*u)
	// log.Print(u)
	has, _ := DBengine.Exist(u)
	return has
}

func (u *User) Register() error {
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

func (e *User) PushMessage(msg string) error {
	_, err := DBengine.Insert(&Message{
		SenderID:   e.ID,
		ReceiverID: e.ID,
		Content:    msg,
		Created:    time.Now(),
		Modified:   time.Now(),
	})
	return err
}

func (s *Group) SetAdmin(id string) error {
	if has, _ := DBengine.Exist(&Group{
		GroupID:  s.GroupID,
		UserList: []string{id},
	}); has {
		s.AdminList = append(s.AdminList, id)
		RemoveSlice(s.UserList, id)
		DBengine.Update(&Group{}, s)
		// e.Table(&Group{}).Update(s)
	}
	return nil
}
