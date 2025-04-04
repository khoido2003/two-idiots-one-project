package routers

func (r *RouterContext) CartsRoute() {
	r.v1Router.Post("/cart", r.handlerContext.AddToCart)
	r.v1Router.Get("/cart", r.handlerContext.GetCart)
    r.v1Router.Patch("/cart/{id}", r.handlerContext.UpdateCartItem)
    r.v1Router.Delete("/cart/{id}", r.handlerContext.RemoveFromCart)
}
