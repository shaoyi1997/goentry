package view

var Templates = &Template{
	Login:    "Login",
	Register: "Register",
	Profile:  "Profile",
}

type TemplateString string

type Template struct {
	Login, Register, Profile TemplateString
}
