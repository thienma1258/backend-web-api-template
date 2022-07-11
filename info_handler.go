package api

import "fmt"

func (app *App) InfoHandler(ctx *RequestContext) {
	_ = ctx.SendSuccess(mapObject{
		"service":     fmt.Sprintf("SERVICE : %v", app.Config.ServiceName),
		"commitHash":  app.serviceInfo.CommitHash,
		"containerID": app.serviceInfo.NodeID,
		"ipAddress":   ctx.RemoteIP,
	})
}

func (app *App) Ping(ctx *RequestContext) {
	_ = ctx.SendSuccess(mapObject{
		"service":     fmt.Sprintf("SERVICE : %v", app.Config.ServiceName),
		"commitHash":  app.serviceInfo.CommitHash,
		"containerID": app.serviceInfo.NodeID,
		"message":     "pong",
	})
}
