package models

import (
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
	_"fmt"
)

func init() {
	//mysql
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//databaseconfig := fmt.Sprintf("%s:%s@/%s?charset=UTF8MB4",beego.AppConfig.String("mysqluser"),beego.AppConfig.String("mysqlpass"),beego.AppConfig.String("databasename"))
	//orm.RegisterDataBase("default", "mysql", databaseconfig)
	//sqlite3
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./datas/phantom.db")

	orm.RegisterModel(new(User),new(Subject),new(WrongSubmitTable),new(RightSubmitTable),new(Game),new(SubjectFile))
	orm.RunSyncdb("default", false, true)
	orm.Debug = true
}

type State int

const (
	WellOp  State = iota	//Everything is ok
	DatabaseErr              // Internal database error
	NoSuchKey

	//user status
    PassWrong // wrong password
	UserRepeat            // User already exists (when registering)
	EmailRepeat            // Email already exists (when registering)
	NoExistUser   //User does not exist
	MarkEditWrong //Failed to modify score
	NoActive //inactivated
	FailActive //Activation fails
	ActiveRepeat//Repeat activation
	NewAndOldDiff//Inconsistent old and new passwords

	//Question status
	NoSuchSubject
	NoSuchId
	
	//Submit flag status
	FlagWrong //flag error
	NoRightSubmit //No record submitted successfully
	HasRightSubmit //Have a record of successful submission

	//Title file status
	FileDeleteError//Failed to delete title file
)