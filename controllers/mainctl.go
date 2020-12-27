package controllers

import (
	"phantom/models"
	"phantom/tools"
	_ "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"html/template"
	"net/url"
	"strings"
	"time"
)

//Verification code generator
var cpt *captcha.Captcha
var gamecommon models.Game

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
	gamecommon, _ = models.GetGameSetting()
}

//Error page controller
type ErrorController struct {
    beego.Controller
}

func (c *ErrorController) Error404() {
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
    c.Data["content"] = "The page is lost, I can't find it QAQ"
    c.TplName = "404.html"
}

//Install the controller of the page
type SetupController struct {
	beego.Controller
}

func (c *SetupController) Prepare() {
	c.EnableXSRF = true
}

//Get method to the controller.
func (c *SetupController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "setup.html"
}

//Post method to the controller.
func (c *SetupController) Post() {
	flash := beego.NewFlash()
	gamename := strings.TrimSpace(c.GetString("gamename"))
	username := strings.TrimSpace(c.GetString("adminname"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("The two input passwords are inconsistent!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
	email := strings.TrimSpace(c.GetString("email"))
	activestring := tools.Md5Encode(time.Now().String())
	if status := models.RegisterUser(username, password, email, activestring,1,true); status == models.WellOp {
		var game models.Game
		game.IfSetup = true
		game.GameName = gamename
		game.Id = 1
		game.IfUseEmail = false
		if status := models.GameSetting(game); status == models.WellOp{
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/")
			gamecommon, _ = models.GetGameSetting()
			return
		}else {
			flash.Error("Database error")
			//TODO:Clear user table after database error
			flash.Store(&c.Controller)
			c.Ctx.Redirect(302, "/setup")
			return
		}
	} else {
		flash.Error("Database error")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/setup")
		return
	}
}

//Homepage controller
type IndexController struct {
	beego.Controller
}

func (c *IndexController) Prepare() {
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *IndexController) Get() {
	c.TplName = "index.html"
}

//Controller for login page
type LoginController struct {
	beego.Controller
}

func (c *LoginController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
		c.Ctx.Redirect(302, "/")
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *LoginController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "login.html"
}

//Post method to the controller
func (c *LoginController) Post() {
	flash := beego.NewFlash()
	if cpt.VerifyReq(c.Ctx.Request) {

	} else {
		flash.Error("Verification code error!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	}
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	var state models.State
	state = models.LoginUser(username, password)
	if state == models.PassWrong {
		flash.Error("wrong password!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else if state == models.NoExistUser {
		flash.Error("User does not exist!!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else if state == models.WellOp {
		c.SetSession("user", username)

		state, isadmin := models.IfAdmin(username)
		if state != models.WellOp {
			c.Ctx.Redirect(302, "/")
			return
		}
		if isadmin {
			c.SetSession("admin", username)
		}
		c.Ctx.Redirect(302, "/")
	} else if state == models.NoActive {
		flash.Error("The user is not activated, please go to the mailbox to activate first!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	} else {
		flash.Error("Database error!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/login")
		return
	}
}

//Controller for registration page
type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

//Get method to the controller.
func (c *RegisterController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "signup.html"
}

//Post method to the controller.
func (c *RegisterController) Post() {
	flash := beego.NewFlash()
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("The two input passwords are inconsistent!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
	email := strings.TrimSpace(c.GetString("email"))
	activestring := tools.Md5Encode(time.Now().String())
	gameconfig, _ := models.GetGameSetting()
	if status := models.RegisterUser(username, password, email, activestring,0,!gameconfig.IfUseEmail); status == models.WellOp {
		if gameconfig.IfUseEmail{
			emailhost, _ := url.Parse(gameconfig.EmailHost)
			gameurl, _ := url.Parse(gameconfig.GameUrl)
			tools.SendEmailActive(email, username, activestring, emailhost.Hostname(), gameurl.Hostname(), gameurl.Port(), gameconfig.EmailAcount, gameconfig.EmailPass, gameconfig.EmailPort)
			flash.Notice("registration success! Please go to your email to activate.")
		}else{
			flash.Notice("registration success!")
		}

		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	} else {
		if status == models.EmailRepeat {
			flash.Error("The mailbox has been registered!")
		} else if status == models.UserRepeat {
			flash.Error("Username already exists!")
		} else {
			flash.Error("Database error")
		}
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/signup")
		return
	}
}

//Logout controller
type LogoutController struct {
	beego.Controller
}

func (c *LogoutController) Prepare() {
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/")
	}
}

func (c *LogoutController) Get() {
	c.DestroySession()
	c.Ctx.Redirect(302, "/")
}

func (c *LogoutController) Post() {
	c.DestroySession()
	c.Ctx.Redirect(302, "/")
}

//Active page controller
type ActiveUserController struct {
	beego.Controller
}

func (c *ActiveUserController) Prepare() {
	userSess := c.GetSession("user")
	if userSess != nil {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *ActiveUserController) Get() {
	username := c.Ctx.Input.Param(":user")
	activestring := c.Ctx.Input.Param(":activestring")
	if status := models.ActiveUser(username, activestring); status == models.WellOp {
		c.Data["Info"] = "Your account has been activated successfully!"
	} else if status == models.NoExistUser {
		c.Ctx.Redirect(302, "/")
		return
	} else if status == models.FailActive {
		c.Data["Info"] = "Activation of your account failed!"
	} else if status == models.DatabaseErr {
		c.Data["Info"] = "Database error, please contact the administrator!"
	} else if status == models.ActiveRepeat {
		c.Data["Info"] = "You have already activated, please do not activate again!"
	}
	c.TplName = "active.html"
}

//Controller for personal settings page
type UserSettingController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *UserSettingController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *UserSettingController) Get() {
	username := c.GetSession("user").(string)
	_, user := models.GetUserInfo(username)
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["username"] = user.Username
	c.Data["name"] = user.Name
	c.Data["userid"] = user.Stuid
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.TplName = "usersetting.html"
}

func (c *UserSettingController) Post() {
	flash := beego.NewFlash()
	name := strings.TrimSpace(c.GetString("name"))
	stuid := strings.TrimSpace(c.GetString("stuid"))
	username := c.GetSession("user").(string)
	var user models.User
	user.Name = name
	user.Stuid = stuid
	user.Username = username
	if models.UpdateUserInfo(user) == models.WellOp {
		flash.Notice("Successfully modified!")
		flash.Store(&c.Controller)
	} else {
		flash.Error("fail to edit!")
		flash.Store(&c.Controller)
	}
	c.Ctx.Redirect(302, "/usersetting")
}

//Modify the controller of the password page
type ChangePwdController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *ChangePwdController) Prepare() {
	c.EnableXSRF = true
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *ChangePwdController) Get() {
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	c.TplName = "changepass.html"
}

//Post method to the controller.
func (c *ChangePwdController) Post() {
	username := c.GetSession("user").(string)
	flash := beego.NewFlash()
	oldpass := strings.TrimSpace(c.GetString("oldpass"))
	password := strings.TrimSpace(c.GetString("password"))
	vrpassword := strings.TrimSpace(c.GetString("veripassword"))
	if password != vrpassword {
		flash.Error("The two input passwords are inconsistent!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}
	if password == oldpass {
		flash.Error("The old and new passwords cannot be the same!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}
	if status := models.UpdatePassword(username, oldpass, password); status == models.WellOp {
		flash.Notice("Successfully modified!")
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	} else {
		if status == models.NewAndOldDiff {
			flash.Error("The old and new passwords are inconsistent!")
		} else {
			flash.Error("Database error")
		}
		flash.Store(&c.Controller)
		c.Ctx.Redirect(302, "/changepwd")
		return
	}
}


//Controller for user information page
type UserInfoController struct {
	beego.Controller
}

//Prepare method to the controller.
func (c *UserInfoController) Prepare() {
	userSess := c.GetSession("user")
	if userSess == nil {
		c.Ctx.Redirect(302, "/login")
	} else {
		c.Data["IsLogin"] = true
		c.Data["UserName"] = userSess
		adminSess := c.GetSession("admin")
		if adminSess != nil {
			c.Data["IsAdmin"] = true
		}
	}
	c.Data["GameName"] = gamecommon.GameName
}

func (c *UserInfoController) Get() {
	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["Error"] = true
	}
	username := c.Ctx.Input.Param(":username")
	userinfo,state := models.FindUnHiddenUsersByUsername(username)
	if state == models.WellOp{
		c.Data["UserInfo"] = userinfo
	}else{
		c.Ctx.Abort(404,"404")
	}
	c.TplName = "userinfo.html"	
}