package routers

import "github.com/go-chi/chi/v5"

func (r *RouterContext) AdminRoute() {
	r.v1Router.Group(func(adminRouter chi.Router) {
		adminRouter.Use(r.handlerContext.AdminMiddleware)
		adminRouter.Get("/admin/dashboard", r.handlerContext.GetDashboardData)
		adminRouter.Delete("/admin/products/{id}", r.handlerContext.DeleteProduct)
		adminRouter.Patch("/admin/products/{id}", r.handlerContext.UpdateProduct)
		adminRouter.Patch("/admin/order/{id}", r.handlerContext.UpdateOrderStatus)

		adminRouter.Get("/admin/orders/{id}", r.handlerContext.GetOrder)
		adminRouter.Get("/admin/orders", r.handlerContext.GetAllOrders)

		adminRouter.Patch("/admin/users/{id}", r.handlerContext.UpdateUserRole)
	})
}
