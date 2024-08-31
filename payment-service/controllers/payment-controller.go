package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go"

	"github.com/stripe/stripe-go/charge"
)

func CreatePayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		stripe.Key = os.Getenv("STRIPE_KEY")
		ReceiptEmail := "abidaziz9876@gmail.com"
		charge, err := charge.New(&stripe.ChargeParams{
			Amount:       stripe.Int64(100),
			Currency:     stripe.String(string(stripe.CurrencyUSD)),
			Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
			ReceiptEmail: stripe.String(ReceiptEmail),
		})

		if err != nil {
			ctx.String(http.StatusBadRequest, "Request failed")
			return
		}
		fmt.Println(charge.BillingDetails)
		ctx.String(http.StatusCreated, "Successfully charged")
	}
}
