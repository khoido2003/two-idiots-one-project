package routers

func(r *RouterContext) OrderRoute() {
    r.v1Router.Get("/orders", r.handlerContext.GetOrders)
}
