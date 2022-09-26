package shared

// Example - GetService(func(s *some type*) {
//	  *outer scope variable* = s
// })
type IServiceProvider interface {
	GetService(function interface{}) error
	AddService(constructor *ServiceDescriptor) error
	RunAfterBuild(functionList []func()) error
}
