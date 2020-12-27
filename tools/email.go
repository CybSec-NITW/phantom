package tools

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

//Send activation email
func SendEmailActive(reciever, username, activestring, emailhost, gameurl, gameport, emailacount, emailpass string, emailport int) {
	config := fmt.Sprintf(`{"username":"%s","password":"%s","host":"%s","port":%d}`, emailacount, emailpass, emailhost, emailport)
	// Create an Email object by storing a string of configuration information
	temail := utils.NewEMail(config)
	// Specify the basic information of the mail
	temail.To = []string{reciever} //Specify the recipient's email address
	temail.From = "ctf@cybsec.in"   //Specify the email address of the sender
	temail.Subject = "Activate your phantom account" //Specify the subject of the message
	temail.HTML = fmt.Sprintf(`<html>
		<head>
		</head>
			 <body>
			   <h1>Hello，%s，Welcome to participate in this CTF competition!</h1>
			   <br>
			   <h2>Click the hyperlink to complete the activation <a href="http://%s:%s/active/%s/%s" target="_brank">CLICK</a></h2>
			   <br>
			   <h2>Hack and have fun!</h2>
	     	</body>
	 	</html>`, username, gameurl, gameport, username, activestring) //Specify the message content
	// send email
	err := temail.Send()
	if err != nil {
		beego.Error("Failed to send mail:", err)
		return
	}
}
