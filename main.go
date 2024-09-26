package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"stripe-test/config"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
)

func init() {
	config.LoadEnv("resources/service.properties")
}

func main() {
	stripeKey := os.Getenv("stripe.Key")
	fmt.Println(stripeKey)
	r := gin.Default()
	stripe.Key = stripeKey
	// Webhook route for handling Stripe events
	r.POST("/webhook", handleWebhook)
	r.POST("/create-payment-intent", createPaymentIntent)
	r.POST("/create-customer", createCustomer)
	r.Run(":4242")

}

func createCustomer(c *gin.Context) {
	params := &stripe.CustomerParams{
		Email: stripe.String("navyatest@gmail.com"),
		Name:  stripe.String("Navya"),
	}
	nc, err := customer.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"customer_id": nc.ID})
}

func createPaymentIntent(c *gin.Context) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(1099), // Amount in cents ($10.99)
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"client_secret": pi.ClientSecret})
}

func handleWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), "")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Webhook signature verification failed: %v", err))
		return
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		// Handle successful payment intent
		fmt.Println("PaymentIntent was successful!")
	case "payment_intent.payment_failed":
		// Handle failed payment intent
		fmt.Println("PaymentIntent failed.")
	case "payment_intent.created":
		fmt.Println("PaymentIntent created.")
	// Add more cases for other event types as needed
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	// Respond with a 200 status to acknowledge receipt of the event
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
