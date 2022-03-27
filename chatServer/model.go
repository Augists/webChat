package main

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"time"
	"xorm.io/xorm"
)

var (
	DBdriver = "mysql"
	DBinfo   = "augists:A1b2c3d4,.@/chatServer?charset=utf8"
	DBengine *xorm.Engine
)

type User struct {
	ID        string
	Name      string
	Salt      string
	Password  string `xorm:"varchar(200)"`
	GroupList string `xorm:"varchar(200)"`
}

type Group struct {
	GroupID string
	Name    string
	/*
	 * Leader and AdminList should also be in UserList
	 */
	Leader    string
	AdminList string `xorm:"varchar(200)"`
	UserList  string `xorm:"varchar(200)"`
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
		DBengine.Sync2(new(User), new(Group))
		DBengine.Sync2(new(Message), new(GroupMessage))
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

// func (u *User) PushMessage(msg string) error {
// 	_, err := DBengine.Insert(&Message{
// 		SenderID:   u.ID,
// 		ReceiverID: u.ID,
// 		Content:    msg,
// 		Created:    time.Now(),
// 		Modified:   time.Now(),
// 	})
// 	return err
// }

/*
 * only use in go routine
 * should be called after wg.Add(n)
 */
func (u *User) PullMessage(conn net.Conn) {
	/*
	 * get all message whose receiver is u.ID and send to conn
	 */
	msg := []Message{}
	DBengine.Where("receiver_id = ?", u.ID).Find(&msg)
	for _, m := range msg {
		conn.Write([]byte(m.Content))
	}
	/*
	 * delete all message whose receiver is u.ID
	 */
	DBengine.Where("receiver_id = ?", u.ID).Delete(&Message{})
	wg.Done()
}

/*
 * only use in go routine
 * should be called after wg.Add(n)
 */
func (u *User) PullGroupMessage(conn net.Conn) {
	/*
	 * complete u *User infomation especially GroupList
	 */
	DBengine.Get(u)
	groupIDList := StringToList(u.GroupList)
	for _, gid := range groupIDList {
		msg := []GroupMessage{}
		DBengine.Where("group_id = ?", gid).Find(&msg)
		for _, m := range msg {
			conn.Write([]byte(m.Content))
		}
		/*
		 * delete all message whose receiver is u.ID
		 */
		DBengine.Where("group_id = ?", gid).Delete(&GroupMessage{})
	}
	wg.Done()
}

func (u *User) JoinGroup(gid string) error {
	u.GroupList = u.GroupList + "," + gid
	_, err := DBengine.Update(u)
	return err
}

func (u *User) QuitGroup(gid string) error {
	groupList := StringToList(u.GroupList)
	for i, g := range groupList {
		if g == gid {
			groupList = append(groupList[:i], groupList[i+1:]...)
			break
		}
	}
	u.GroupList = ListToString(groupList)
	_, err := DBengine.Update(u)
	return err
}

func (s *Group) SetAdmin(id string) error {
	if has, _ := DBengine.Exist(&Group{
		GroupID: s.GroupID,
	}); has {
		adminList := StringToList(s.AdminList)
		adminList = append(adminList, id)
		s.AdminList = ListToString(adminList)
		/*
		 * Leader and AdminList should also be in UserList
		 */
		// RemoveSlice(s.UserList, id)
		DBengine.Update(&Group{}, s)
		// e.Table(&Group{}).Update(s)
	}
	return nil
}

func (m *Message) Store() {
	DBengine.Insert(m)
}

func (m *GroupMessage) Store() {
	DBengine.Insert(m)
}

func GetGroupUserList(gid string) []string {
	group := Group{}
	DBengine.Where("group_id = ?", gid).Get(&group)
	return StringToList(group.UserList)
}
