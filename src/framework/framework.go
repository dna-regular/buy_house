package framework

type Module interface {
	Init()
	NextUrl() (url string, done bool, err error)
	ParsePage(html string) (interface{}, error)
}

func ModuleRegister(modules []interface{}, module interface{}) {
	modules = append(modules, module)
}
