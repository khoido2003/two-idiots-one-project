package routers

func (routerContext *RouterContext) UsersRoute() {
	routerContext.v1Router.Post("/signin", routerContext.handlerContext.Signin)
	routerContext.v1Router.Post("/signup", routerContext.handlerContext.Signup)

}
