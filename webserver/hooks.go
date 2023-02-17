package webserver

type (
	HookFunction func(r Resource) error
	Hook         interface {
		BeforeRun(hf HookFunction) *WebServer
		AfterRun(hf HookFunction) *WebServer
		BeforeExit(hf HookFunction) *WebServer
		AfterExit(hf HookFunction) *WebServer
	}
)

func (w *WebServer) BeforeRun(hf HookFunction) *WebServer {
	w.beforeRun = append(w.beforeRun, hf)
	return w
}

func (w *WebServer) AfterRun(hf HookFunction) *WebServer {
	w.afterRun = append(w.afterRun, hf)
	return w
}

func (w *WebServer) BeforeExit(hf HookFunction) *WebServer {
	w.beforeExit = append(w.beforeExit, hf)
	return w
}

func (w *WebServer) AfterExit(hf HookFunction) *WebServer {
	w.afterExit = append(w.afterExit, hf)
	return w
}
