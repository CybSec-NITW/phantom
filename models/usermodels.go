package models

import (
	//_ "github.com/go-sql-driver/mysql"
	_"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"phantom/tools"
	"time"
	_"fmt"
)

//The user's data model.
type User struct {
	Id	int
	Mark	int //User score
	Name	string `orm:"size(100)"` //username
	Email	string
	Stuid	string //student ID
	Username	string //username
	Hashpass	string //Hashed password
	Identity	int //Identifies whether it is an administrator (1) or a normal user (0)
	IfActive	bool//Activate now
	IfHidden	bool//Whether to hide
	ActiveString string//Activation link
	CreatedTime	time.Time	`orm:"auto_now_add;type(datetime)"`
	UpdatedTime	time.Time	`orm:"auto_now;type(datetime)"`
}

//User answer information, no databas
type UserSubmitInfo struct{
	SubjectName	string//Topic name
	SubjectPoint	int//Question score
	SubmitTime	time.Time//Submission time
}

//User information, no database
type UserInfo struct{
	Username	string//username
	Email	string//User mailbox
	Rank	int//Rank
	Mark	int//Score
	SubmitTable	[]UserSubmitInfo//Answer information
}

//Database operations for registered users
func RegisterUser(username, password, email ,activestring string,ifadmin int,ifactive bool)(state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	user.Email = email
	user.Hashpass = tools.Md5Encode(password)
	user.Identity = ifadmin
	if ifadmin == 1{
		user.IfActive = true
		user.IfHidden = true
	}else{
		user.IfActive = ifactive
		user.IfHidden = false
	}
	user.Mark = 0
	user.ActiveString = activestring
	var err error

	//First check if the username is duplicate
	err = o.Read(user,"Username")
	if err == orm.ErrNoRows {
	} else if err == orm.ErrMissPK {
		
	} else {
		state = UserRepeat
		return
	}

	//Check if the mailbox is duplicate
	err = o.Read(user,"Email")
	if err == orm.ErrNoRows {
		
	} else if err == orm.ErrMissPK {
		
	} else {
		state = EmailRepeat
		return
	}

	_, err = o.Insert(user)
	if err != nil {
		state = DatabaseErr
		return
	}else{
		state = WellOp
		return
	}
}

//Login database operation
func LoginUser(username,password string) (state State) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.Hashpass != tools.Md5Encode(password){
		state = PassWrong
		return
	}else if user.IfActive == false{
		state = NoActive
		return
	}else{
		state = WellOp
		return
	}
}

//Activate user's database operation
func ActiveUser(username,activestring string) (state State) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.IfActive{
		state = ActiveRepeat
		return
	}

	if user.ActiveString != activestring{
		state = FailActive
		return
	}else{
		user.IfActive = true
		_, err := o.Update(user,"IfActive")
		if err != nil {
			state = DatabaseErr
		}else{
			state = WellOp
		}
		return
	}
}


//Database operations to determine user roles

func IfAdmin(username string) (state State,isadmin bool) {
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	err := o.Read(user,"Username")
	if err != nil{
		state = NoExistUser
		return
	}
	if user.Identity == 1{
		state = WellOp
		isadmin = true
		return
	}else{
		state = WellOp
		isadmin = false
		return
	}

}

//Get the actions of all unhidden users (ranking list）
func GetUnhiddenUsers() (state State,users []User){
	o := orm.NewOrm()
	o.QueryTable("user").OrderBy("-Mark").Filter("IfHidden", false).All(&users, "Username","Mark")
	state = WellOp
	return
}

//Get all user operations (user management)
func GetUsers() (state State,users []User){
	o := orm.NewOrm()
	o.QueryTable("user").All(&users, "Username","Mark")
	state = WellOp
	return
}

//Obtain a user information (for further operations）
func GetUserInfo(username string)(state State,user User){
	o := orm.NewOrm()
	user = User{Username:username}
	err := o.Read(&user,"Username")

	if err == orm.ErrNoRows {
		state = NoExistUser
	} else if err == orm.ErrMissPK {
    	state = NoSuchKey
	} else {
    	state = WellOp
	}
	return
}

//Update user information
func UpdateUserInfo(user User)(state State){
	o := orm.NewOrm()
	var olduser User
	olduser.Username = user.Username
	if o.Read(&olduser,"Username") == nil {
		olduser.Name = user.Name
		olduser.Stuid = user.Stuid
		if _, err := o.Update(&olduser,"Name","Stuid"); err == nil {
			state = WellOp
		}
	}else{
		state = DatabaseErr
	}
	return
}

//change Password
func UpdatePassword(username,oldpassword,password string) (state State){
	o := orm.NewOrm()
	user := new(User)
	user.Username = username
	user.Hashpass = tools.Md5Encode(oldpassword)
	err := o.Read(user,"Username","Hashpass")

	if err == orm.ErrNoRows {
    	state = NewAndOldDiff
	} else if err == orm.ErrMissPK {
    	state = NoSuchKey
	} else {
			user.Hashpass = tools.Md5Encode(password)
			if _, err := o.Update(user,"Hashpass"); err == nil {
				state = WellOp
			}else{
				state = DatabaseErr
			}
	}
	return
}

//Find unhidden users by username
func FindUnHiddenUsersByUsername(username string)(userinfo UserInfo,state State){
	o := orm.NewOrm()
	user := User{Username: username,IfHidden:false}
	err := o.Read(&user,"Username","IfHidden")
	if err == orm.ErrNoRows {
		state = NoExistUser
		return
	} else if err == orm.ErrMissPK {
		state = NoSuchKey
		return
	} else {
		userinfo.Username = user.Username
		userinfo.Email = user.Email
		userinfo.Mark = user.Mark
		var users []User
		state,users = GetUnhiddenUsers()
		userinfo.Rank = 1
		for _,eachuser := range users{
			if eachuser.Username == user.Username{
				break
			}else{
				userinfo.Rank = userinfo.Rank + 1
			}
		}
		var rstable []RightSubmitTable
		_, err := o.QueryTable("right_submit_table").OrderBy("-CreatedTime").Filter("user_name",userinfo.Username).All(&rstable)
		if err != nil{
			state = DatabaseErr
		}else{
			userinfo.SubmitTable = make([]UserSubmitInfo,len(rstable))
			state = WellOp
			for key,rs :=range rstable{
				subject := Subject{Id:rs.SubjectId}
				err := o.Read(&subject)
				if err == orm.ErrNoRows{
					state = DatabaseErr
				} else if err == orm.ErrMissPK{
					state = NoSuchKey
				}else{
					userinfo.SubmitTable[key].SubjectName = subject.SubName
					userinfo.SubmitTable[key].SubjectPoint = subject.SubMark
					userinfo.SubmitTable[key].SubmitTime = rs.CreatedTime							
					state = WellOp
				}
			}
		}
		return
	}
}