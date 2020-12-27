package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"fmt"
	"strconv"
	"time"
	"os"
)


//Topic database model
type Subject struct {
	Id       int
	IfHidden bool    //Whether the title is hidden
	SubName     string `orm:"size(100)"` //Title name
	SubMark		int //Question score
	SubFlag		string `orm:"size(200)"` //flag
	SubDescribe string `orm:"size(1000)"` //Title description
	SubType		string `orm:"size(50)"` //Question type
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`
	IfDone	bool	//Temporarily determine whether the current user answered the question correctly
}

//The database model of the game
type Game struct {
	Id	int
	IfSetup	bool//Whether the configuration is complete
	GameName	string//Game name
	GameUrl	string//Competition domai
	IfUseEmail	bool
	EmailHost	string//Mail Serve
	EmailPort	int//Server port
	EmailAcount	string//Mail account
	EmailPass	string//Mail password
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`			
}

//Database model of the title attachment
type SubjectFile struct{
	Id	int
	SubId	int//The ID of the corresponding question
	FileName	string//File name downloaded
	Md5FileName	string//File name stored after Md5
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`		
}


//Get an array of all questions for further manipulation or rendering of the page
func GetSubjects() (state State,subjects []Subject){
	o := orm.NewOrm()
	//o.QueryTable("subject").Filter("Status", 1).All(&subjects, "Id", "Title")
	o.QueryTable("subject").All(&subjects, "Id", "SubName","SubType")
	state = WellOp
	return
}

//Get information on a certain topic for editing
func GetSubject(id int) (state State,subject Subject){
	o := orm.NewOrm()
	subject.Id = id
	err := o.Read(&subject)
	if err == orm.ErrNoRows {
		state = NoSuchSubject
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
	} else {
		state = WellOp
	}
	return
}

//Add title
func AddSubject(subject Subject) (state State){
	o := orm.NewOrm()
	_, err := o.Insert(&subject)
	if err != nil{
		state = DatabaseErr
		return
	}else{
		state = WellOp
		return
	}
}

//Modify question
func EditSubject(subject Subject) (state State){

	o := orm.NewOrm()
	oldsubject := Subject{Id:subject.Id}
	if o.Read(&oldsubject) == nil {
		oldsubject.IfHidden = subject.IfHidden
		oldsubject.SubName = subject.SubName
		oldsubject.SubMark = subject.SubMark
		oldsubject.SubFlag	 = subject.SubFlag
		oldsubject.SubDescribe = subject.SubDescribe
		oldsubject.SubType = subject.SubType
    if _, err := o.Update(&oldsubject); err == nil {
		state = WellOp
		return
	}
	state = NoSuchSubject
	return
}else{
	state = DatabaseErr
	return
}
}

//Delete the title and delete the corresponding file
func DeleteSubject(id int) (state State){
	o := orm.NewOrm()
	var subfiles []SubjectFile
	subfiles,state = GetSubjectFile(id)
	if state != WellOp{
		return
	}
	_, err := o.QueryTable("subject_file").Filter("sub_id", id).Delete()
	if err !=nil{
		state = DatabaseErr
		return
	}
	for i,_ := range subfiles{
		err := os.Remove("upload/" + subfiles[i].Md5FileName)
		if err != nil{
			state = FileDeleteError
			return
		}
	}
	if id, err := o.Delete(&Subject{Id: id}); err == nil {
	state = WellOp
	if id == 0{
		state = NoSuchSubject
	}
	}else{
		state = DatabaseErr
	}
	return
}

//Get all unhidden questions
func GetUnhiddenSubjects() (state State,subjects []Subject){
	o := orm.NewOrm()
	//o.QueryTable("subject").Filter("Status", 1).All(&subjects, "Id", "Title")
	o.QueryTable("subject").Filter("IfHidden", false).All(&subjects, "Id", "SubName","SubType","SubMark","SubDescribe")
	state = WellOp
	return
}

//Submit flag and record submission history
func UserCommitFlag(subjectId,userFlag,userName string) (state State){
	o := orm.NewOrm()
	subject := new(Subject)
	var errors error
	subject.Id,errors = strconv.Atoi(subjectId)
	if errors != nil{
		state = NoSuchId
		return
	}
	errors = o.Read(subject,"Id")
	if errors != nil{
		state = NoSuchSubject
		return
	}
	if subject.SubFlag != userFlag{
		state = FlagWrong
		WrongSubmit(subject.Id,userName,userFlag)
		return
	}else{
		if IfSolved(subject.Id,userName) == NoRightSubmit{
			state = EditUserMark(userName,subject.SubMark)
			RightSubmit(subject.Id,userName)
			return
		}else{
			return
		}
	}

}

//Add and subtract points operation
func EditUserMark(userName string,userMark int) (state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = userName
	if o.Read(user,"Username") == nil {
		user.Mark += userMark 
		if _, err := o.Update(user,"Mark"); err != nil {
			state = MarkEditWrong
			return
		}
		state = WellOp
		return
	}else{
		state = MarkEditWrong
		return
	}
}

//Game global settings
func GameSetting(game Game)(state State){
	o := orm.NewOrm()
	oldgame := game
	if created, _, err := o.ReadOrCreate(&oldgame, "Id"); err == nil {
		if created {
			state = WellOp
		} else {
			oldgame = game
			if _, err := o.Update(&oldgame); err == nil {
				state = WellOp
			}else{
				state = DatabaseErr
			}
		}
	}else{
		state = DatabaseErr
	}
	return
}

//Get global game settings
func GetGameSetting()(game Game,state State){
	o := orm.NewOrm()
	game.Id = 1
	err := o.Read(&game)
	if err == orm.ErrNoRows {
		state = NoSuchId
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
	} else {
		state = WellOp
	}
	return
}

//Record the file corresponding to the title
func UploadSubjectFile(filename,md5filename string,subjectid int)(state State){
	o := orm.NewOrm()
	var subfile SubjectFile
	subfile.Md5FileName = md5filename
	subfile.FileName = filename
	subfile.SubId = subjectid
	_, err := o.Insert(&subfile)
	if err == nil {
		state = WellOp
	}else{
		state = DatabaseErr
	}
	return
}

//Delete the file corresponding to the question
func DeleteSubjectFile(fileid int)(state State){
	o := orm.NewOrm()
	if _, err := o.Delete(&SubjectFile{Id: fileid}); err == nil {
		state = WellOp
	}else{
		state = DatabaseErr
	}
	return
}

//Get the file corresponding to the question
func GetSubjectFile(subjectid int)(subfile []SubjectFile,state State){
	o := orm.NewOrm()
	_, err := o.QueryTable("subject_file").Filter("sub_id", subjectid).All(&subfile)
	if err != nil{
		state = DatabaseErr
	}else{
		state = WellOp
	}
	return
}

//Get all files
func GetAllFiles()(subfile []SubjectFile,state State){
	o := orm.NewOrm()
	_, err := o.QueryTable("subject_file").All(&subfile)
	if err != nil{
		state = DatabaseErr
	}else{
		state = WellOp
	}
	return
}

//Get the md5 file name corresponding to the file id
func GetFileById(fileid int)(subfile SubjectFile,state State){
	o := orm.NewOrm()
	subfile.Id = fileid
	err := o.Read(&subfile,"Id")
	if err == orm.ErrNoRows {
		state = NoSuchId
	} else if err == orm.ErrMissPK{
		state = NoSuchKey
	}else{
		state = WellOp
		}
	return
}