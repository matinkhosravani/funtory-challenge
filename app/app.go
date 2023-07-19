package app

type Application struct {
	Env *Env
}

var app Application

func Boot() {
	app.Env = NewEnv(false)
}

func GetEnv() *Env {
	return app.Env
}
