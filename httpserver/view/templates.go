package view

var Templates = &Template{
	LoginRegister: "LoginRegister",
	Profile:       "Profile",
	Edit:          "Edit",
}

type TemplateString string

type Template struct {
	LoginRegister, Profile, Edit TemplateString
}

type LoginRegisterPageData struct {
	IsLoginPage     bool
	UsernameTaken   bool
	InvalidUsername bool
	InvalidPassword bool
}
