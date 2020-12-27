package routers

import (
	"phantom/controllers"
	"github.com/astaxie/beego"
	"phantom/models"
	"github.com/astaxie/beego/context"
)

var FilterSetup = func(ctx *context.Context) {
    game,_ := models.GetGameSetting()
    if !game.IfSetup && ctx.Request.RequestURI != "/setup"{
        ctx.Redirect(302, "/setup")
	}else if game.IfSetup && ctx.Request.RequestURI == "/setup"{
		ctx.Abort(404,"404")
	}
}

func init() {

	//404
	beego.ErrorController(&controllers.ErrorController{})
	//routers
	beego.InsertFilter("/*",beego.BeforeRouter,FilterSetup)
	//Installation page
	beego.Router("/setup",&controllers.SetupController{})

	//Home page
	beego.Router("/", &controllers.IndexController{})

	//Personal settings
	beego.Router("/usersetting", &controllers.UserSettingController{})
	beego.Router("/changepwd", &controllers.ChangePwdController{})

	//View user information (scores, rankings, answers)
	beego.Router("/user/:username", &controllers.UserInfoController{})

	//Login and registration related
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/signup", &controllers.RegisterController{})
	beego.Router("/active/:user/:activestring", &controllers.ActiveUserController{})

	//Management page
	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/gamesetting", &controllers.GameManageController{})
	beego.Router("/admin/subjects", &controllers.AdminSubjectsController{})
	beego.Router("/admin/subjects/add", &controllers.SubjectsAddController{})
	beego.Router("/admin/subjects/delete/:id", &controllers.SubjectsDeleteController{})
	beego.Router("/admin/subjects/edit/:id", &controllers.SubjectsEditController{})

	//Contest page
	beego.Router("/game", &controllers.GameController{})
	beego.Router("/rank", &controllers.RankController{})

	//Topic file upload, download and delete
	beego.Router("/admin/subjects/file/upload/:id", &controllers.SubjectsFileUploadController{})
	beego.Router("/game/file/download/:id", &controllers.SubjectFileDownloadController{})
	beego.Router("/admin/subjects/file/delete/:id", &controllers.SubjectsFileDeleteController{})	
}
