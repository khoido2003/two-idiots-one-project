package routers


func (routerContext *RouterContext) CategoriesRoute() {
	routerContext.v1Router.Get("/categories", routerContext.handlerContext.GetAllCategories)

}
