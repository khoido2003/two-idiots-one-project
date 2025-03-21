package routers

func (routerContext *RouterContext) CheckHealth() {
	routerContext.v1Router.Get("/checkhealth", routerContext.handlerContext.CheckHealth)

}
