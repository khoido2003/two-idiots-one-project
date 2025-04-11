package routers

func (r *RouterContext) CheckoutRoute() {
	r.v1Router.Post("/checkout", r.handlerContext.CreatePaymentIntent)

}
