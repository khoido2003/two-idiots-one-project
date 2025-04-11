package handlers

import (
    "encoding/json"
    "net/http"
    "os"
    "server/models"
    "strconv"

    "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/client"
)

type CheckoutRequest struct {
    CartItems []struct {
        ProductID uint `json:"productId"`
        Quantity  int  `json:"quantity"`
    } `json:"cartItems"`
    Address struct {
        Line1      string `json:"line1"`
        Line2      string `json:"line2,omitempty"`
        City       string `json:"city"`
        State      string `json:"state"`
        PostalCode string `json:"postalCode"`
        Country    string `json:"country"`
    } `json:"address"`
}

func (h *HandlerContext) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
    userIdQuery := r.URL.Query().Get("userId")
    userId, err := strconv.ParseUint(userIdQuery, 10, 32)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid userId")
        return
    }

    var req CheckoutRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }

    // Fetch cart items
    var cartItems []models.Cart
    if err := h.db.Preload("Product").Where("user_id = ?", userId).Find(&cartItems).Error; err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to fetch cart items: "+err.Error())
        return
    }

    // Validate cart items match request
    if len(cartItems) != len(req.CartItems) {
        respondWithError(w, http.StatusBadRequest, "Cart items mismatch")
        return
    }
    total := int64(0)
    for i, item := range cartItems {
        if item.ProductID != req.CartItems[i].ProductID || item.Quantity != req.CartItems[i].Quantity {
            respondWithError(w, http.StatusBadRequest, "Cart items mismatch")
            return
        }
        total += int64(item.Product.Price * float64(item.Quantity) * 100)
    }
    total += 1000 // $10 shipping in cents
    total += int64(float64(total) * 0.08) // 8% tax

    // Initialize Stripe client
    stripeKey := os.Getenv("STRIPE_SECRET")
    if stripeKey == "" {
        respondWithError(w, http.StatusInternalServerError, "Stripe secret key not configured")
        return
    }
    stripe.Key = stripeKey
    sc := &client.API{}
    sc.Init(stripe.Key, nil)

    // Create payment intent
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(total),
        Currency: stripe.String(string(stripe.CurrencyUSD)),
        Metadata: map[string]string{
            "userId": strconv.FormatUint(userId, 10),
        },
    }
    paymentIntent, err := sc.PaymentIntents.New(params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create payment intent: "+err.Error())
        return
    }

    // Create order (temporary placement; ideally use webhook)
    order := models.Order{
        UserID:    uint(userId),
        Total:     total,
        AddressLine1: req.Address.Line1,
        AddressLine2: req.Address.Line2,
        City:      req.Address.City,
        State:     req.Address.State,
        PostalCode: req.Address.PostalCode,
        Country:   req.Address.Country,
    }
    for _, item := range cartItems {
        order.Items = append(order.Items, models.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     int64(item.Product.Price * 100), // Store in cents
            Name:      item.Product.Name,
        })
    }
    if err := h.db.Create(&order).Error; err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create order: "+err.Error())
        return
    }

    // Clear cart after order creation
    if err := h.db.Where("user_id = ?", userId).Delete(&models.Cart{}).Error; err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to clear cart: "+err.Error())
        return
    }

    respondWithJson(w, http.StatusOK, map[string]interface{}{
        "message":      "Payment intent created",
        "clientSecret": paymentIntent.ClientSecret,
    })
}
