package routers

func (routerContext *RouterContext) CreateProduct() {
	routerContext.v1Router.Post("/products", routerContext.handlerContext.CreateProduct)
	routerContext.v1Router.Get("/products", routerContext.handlerContext.GetProducts)
	routerContext.v1Router.Get("/products/{id}", routerContext.handlerContext.GetProductByID)

}
