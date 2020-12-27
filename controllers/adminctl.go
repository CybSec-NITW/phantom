package controllers

import (
	"phantom/models"
	_"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"strconv"
	"strings"
	"io"
	"os"
	"time"
	"phantom/tools"
)

//Management page homepage
type AdminController struct {
	beego.Controller
}

func (c *AdminController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("user")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *AdminController) Get() {
	c.TplName = "admin/index.html"
}

//Title management page
type AdminSubjectsController struct {
	beego.Controller
}

func (c *AdminSubjectsController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *AdminSubjectsController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["welldone"] = true //Flash the success message on the topic page, modify, add, delete the success message
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["baddone"] = true //Failure message
	}
	var subjects []models.Subject
	_, subjects = models.GetSubjects()
	c.Data["Subjects"] = subjects
	c.TplName = "admin/subjects.html"
	return
}

//Title editing page, dynamic routing
type SubjectsEditController struct {
	beego.Controller
}

func (c *SubjectsEditController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Get method to the controller.
func (c *SubjectsEditController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["UploadOk"] = true
	} else if _, ok = flash.Data["error"];ok  {
		c.Data["EditError"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	EditId, error := strconv.Atoi(c.Ctx.Input.Param(":id"))
	c.Data["EditId"] = EditId
	if error != nil {
		c.Ctx.Redirect(302, "/admin/subjects")
		return
	}
	state, subject := models.GetSubject(EditId)
	if state != models.WellOp {
		c.Ctx.Redirect(302, "/admin/subjects")
		return
	}
	subjectid, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	subfiles,_ := models.GetSubjectFile(subjectid)
	c.Data["SubFiles"] = subfiles
	c.Data["Subject"] = subject
	c.TplName = "admin/subjectedit.html"
	return
}

//Post method to the controller.
func (c *SubjectsEditController) Post() {
	var subject models.Subject
	var errors error
	flash := beego.ReadFromRequest(&c.Controller)
	subject.Id, errors = strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	subject.SubName = strings.TrimSpace(c.GetString("subname"))
	subject.SubType = strings.TrimSpace(c.GetString("subtype"))
	subject.SubFlag = strings.TrimSpace(c.GetString("subflag"))
	subject.SubDescribe = strings.TrimSpace(c.GetString("subdescribe"))
	subject.SubMark, errors = c.GetInt("submark")
	if errors != nil {
		flash.Error("The score must be an integer!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/edit/"+c.Ctx.Input.Param(":id"))
		return
	}
	if strings.TrimSpace(c.GetString("ifhidden")) == "on" {
		subject.IfHidden = true
	} else {
		subject.IfHidden = false
	}
	models.EditSubject(subject)
	flash.Notice("Successfully modified!")
	flash.Store(&c.Controller)
	c.Ctx.Redirect(302, "/admin/subjects/")
}

//Add title
type SubjectsAddController struct {
	beego.Controller
}

func (c *SubjectsAddController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

//Get method to the controller.
func (c *SubjectsAddController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["adderror"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "admin/subjectadd.html"
	return
}

//Post method to the controller.
func (c *SubjectsAddController) Post() {
	var subject models.Subject
	var errors error
	flash := beego.NewFlash()
	subject.SubName = strings.TrimSpace(c.GetString("subname"))
	subject.SubType = strings.TrimSpace(c.GetString("subtype"))
	subject.SubFlag = strings.TrimSpace(c.GetString("subflag"))
	subject.SubDescribe = strings.TrimSpace(c.GetString("subdescribe"))
	subject.SubMark, errors = c.GetInt("submark")
	if errors != nil {
		flash.Error("The score must be an integer!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/add")
		return
	}
	if strings.TrimSpace(c.GetString("ifhidden")) == "on" {
		subject.IfHidden = true
	} else {
		subject.IfHidden = false
	}
	if models.AddSubject(subject) != models.WellOp {
		flash.Error("Failed to add title!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/add")
		return
	} else {
		flash.Notice("Succeeded in adding question!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}

}

//Delete question
type SubjectsDeleteController struct {
	beego.Controller
}

func (c *SubjectsDeleteController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdmin"] = true
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

//Get method to the controller.
func (c *SubjectsDeleteController) Get() {
	flash := beego.NewFlash()
	subjectId, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("failed to delete!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	if models.DeleteSubject(subjectId) == models.WellOp {
		flash.Notice("successfully deleted!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
	} else {
		flash.Error("failed to delete!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	return
}

//Management competition
type GameManageController struct {
	beego.Controller
}

func (c *GameManageController) Prepare() {
	c.EnableXSRF = true
	adminSess := c.GetSession("admin")
	if adminSess != nil {
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	} else {
		c.Data["IsLogin"] = false
		c.Ctx.Abort(404,"404")
	}
}

func (c *GameManageController) Get() {
	game, _ := models.GetGameSetting()
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["game"] = game
	c.TplName = "admin/gamesettings.html"
	return
}

func (c *GameManageController) Post() {
	gamename := strings.TrimSpace(c.GetString("gamename"))
	domainname := strings.TrimSpace(c.GetString("domainname"))
	emailserver := strings.TrimSpace(c.GetString("emailserver"))
	emailport, _ := c.GetInt("emailport")
	emailaccount := strings.TrimSpace(c.GetString("emailaccount"))
	emailpass := strings.TrimSpace(c.GetString("emailpass"))
	var game models.Game
	game.GameName = gamename
	game.GameUrl = domainname
	game.EmailHost = emailserver
	game.EmailPort = int(emailport)
	game.EmailAcount = emailaccount
	game.EmailPass = emailpass
	game.Id = 1
	game.IfSetup = true
	if strings.TrimSpace(c.GetString("ifuseemail")) == "on" {
		game.IfUseEmail = true
	} else {
		game.IfUseEmail = false
	}
	flash := beego.NewFlash()
	if models.GameSetting(game) == models.WellOp {
		flash.Notice("Successfully modified！")
		flash.Store(&c.Controller)
		gamecommon, _ = models.GetGameSetting()
	} else {
		flash.Error("fail to edit！")
		flash.Store(&c.Controller)
	}
	c.Ctx.Redirect(302, "/admin/gamesetting/")
	return
}


//Topic file upload
type SubjectsFileUploadController struct {
	beego.Controller
}

func (c *SubjectsFileUploadController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

//Post method to the controller.
func (c *SubjectsFileUploadController) Post() {
	flash := beego.ReadFromRequest(&c.Controller)
	subjectid, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("Path error!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	files, err := c.GetFiles("files")
	if err != nil {
		c.Ctx.WriteString("Invalid file")
		return
	}
	for i,_ := range files{
		file,err := files[i].Open()
		defer file.Close()
		if err != nil {
			flash.Error("upload failed!")
			flash.Store(&c.Controller)
			break
		}
		//Determine whether the folder exists, create the folder if it does not exist
		_, err = os.Stat("upload")
		if os.IsNotExist(err) {
			err := os.Mkdir("upload", os.ModePerm)
			if err != nil {
				flash.Error("upload failed!")
				flash.Store(&c.Controller)
				break
				}
			}
		md5filename := tools.Md5Encode(files[i].Filename+time.Now().String())
		dst, err := os.Create("upload/" + md5filename)
		defer dst.Close()
		if err != nil {
			flash.Error("upload failed!")
			flash.Store(&c.Controller)
			break
		}
		if _, err := io.Copy(dst, file); err != nil {
			flash.Error("upload failed!")
			flash.Store(&c.Controller)
			break
		}
		if models.UploadSubjectFile(files[i].Filename,md5filename,subjectid) != models.WellOp{
			flash.Error("Database error!")
			flash.Store(&c.Controller)
			break			
		}
	}
	if err == nil{
		flash.Notice("Upload successfully!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/edit/"+strconv.Itoa(subjectid)+"#theupload")
	}else{
		c.Ctx.Redirect(302, "/admin/subjects")
	}

}

//Title file deletio
type SubjectsFileDeleteController struct {
	beego.Controller
}

func (c *SubjectsFileDeleteController) Prepare() {
	adminSess := c.GetSession("admin")
	if adminSess == nil {
		c.Ctx.Abort(404,"404")
	} else {
		adminSess := c.GetSession("admin")
		c.Data["IsLogin"] = true
		c.Data["IsAdmin"] = true
		c.Data["UserName"] = adminSess
		c.Data["IsAdminPage"] = true
	}
}

func (c *SubjectsFileDeleteController) Get(){
	flash := beego.ReadFromRequest(&c.Controller)
	fileid, errors := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errors != nil {
		flash.Error("Path error!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/admin/subjects/")
		return
	}
	subfile,state := models.GetFileById(fileid)
	md5filename := subfile.Md5FileName
	if state != models.WellOp{
		flash.Error("Database error!")
		flash.Store(&c.Controller)			
	}else{
		if models.DeleteSubjectFile(fileid) != models.WellOp{
			flash.Error("Attachment deletion failed!")
			flash.Store(&c.Controller)	
		}else{
			err := os.Remove("upload/" + md5filename)
			if err != nil {
				flash.Error("Attachment deletion failed!")
				flash.Store(&c.Controller)
			}else {
				flash.Notice("Attachment deleted successfully！")
				flash.Store(&c.Controller)
			}
		}
	}
	c.Ctx.Redirect(302, "/admin/subjects/")
	return
}