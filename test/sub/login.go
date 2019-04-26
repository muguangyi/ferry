package sub

type ILogin interface {
	Login(name string, pwd string) bool
	Logout()
}
