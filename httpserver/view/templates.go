package view

var Templates = &Template{
	Login:    "Login",
	Register: "Register",
	Profile:  "Profile",
	Edit:     "Edit",
}

type TemplateString string

type Template struct {
	Login, Register, Profile, Edit TemplateString
}
