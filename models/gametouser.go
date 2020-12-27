package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"time"
)

//Error submission record form
type WrongSubmitTable struct{
	Id	int
	SubmitFlag	string	//Submitted content
	UserName	string	//username
	SubjectId	int	//Topic id
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
}

//Submit the record form correctly
type RightSubmitTable struct{
	Id	int
	UserName	string
	SubjectId	int
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
}

//Determine whether the user has successfully submitted the correct flag
func IfSolved(subjectId int,userName string)(state State){
	o := orm.NewOrm()
	rstable := RightSubmitTable{UserName: userName,SubjectId:subjectId}
	err := o.Read(&rstable,"UserName","SubjectId")
	if err == orm.ErrNoRows {
		state = NoRightSubmit
		return
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
		return
	} else {
		state = HasRightSubmit
		return
	}
}

//Record correctly submitted records
func RightSubmit(subjectId int,userName string)(state State){
	o := orm.NewOrm()
	var rstable RightSubmitTable
	rstable.SubjectId = subjectId
	rstable.UserName = userName
	_, err := o.Insert(&rstable)
	if err == nil {
		state = WellOp
		return
	}else{
		state = DatabaseErr
		return
	}
}

//Record the wrong submission record
func WrongSubmit(subjectId int,userName string,submitFlag string)(state State){
	o := orm.NewOrm()
	var wstable WrongSubmitTable
	wstable.SubjectId = subjectId
	wstable.UserName = userName
	wstable.SubmitFlag = submitFlag
	_, err := o.Insert(&wstable)
	if err == nil {
		state = WellOp
		return
	}else{
		state = DatabaseErr
		return
	}
}