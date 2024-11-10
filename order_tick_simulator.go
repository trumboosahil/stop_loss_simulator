package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"context"

	"github.com/go-redis/redis/v8"
)

type TickEvent struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type RedisOrder struct {
	OrderID       int     `json:"order_id"`
	StopLossPrice float64 `json:"stop_loss_price"`
	Expiry        int64   `json:"expiry"`
	Symbol        string  `json:"symbol"`
}

func main() {
	redisAddr := "localhost:6379"
	if len(os.Args) > 1 {
		redisAddr = os.Args[1]
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	fmt.Println("Starting Order and Tick Simulator...")

	// Step 1: Preload Millions of Orders
	preloadOrders(redisClient, 1000000)

	// Step 2: Start Simulating Tick Events
	simulateTickEvents(redisClient)
}

// Preload millions of orders with stop-loss conditions into Redis
func preloadOrders(redisClient *redis.Client, count int) {
	symbols := []string{"AAPL", "GOOG", "TSLA", "AMZN", "MSFT"}
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Preloading %d orders...\n", count)

	for i := 1; i <= count; i++ {
		symbol := symbols[rand.Intn(len(symbols))]
		stopLossPrice := rand.Float64()*100 + 50         // Random stop-loss price between 50 and 150
		expiry := time.Now().Add(5 * time.Minute).Unix() // Expiry as Unix timestamp

		order := RedisOrder{
			OrderID:       i,
			StopLossPrice: stopLossPrice,
			Expiry:        expiry,
			Symbol:        symbol,
		}

		// Serialize to JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			fmt.Printf("Failed to marshal order %d: %v\n", i, err)
			continue
		}

		// Use the stop-loss price as the score for sorting
		err = redisClient.ZAdd(context.Background(), "orderset", &redis.Z{
			Score:  stopLossPrice,
			Member: orderJSON,
		}).Err()

		if err != nil {
			fmt.Printf("Failed to add order %d to Redis: %v\n", i, err)
		}

		if i%10000 == 0 {
			fmt.Printf("Loaded %d orders...\n", i)
		}
	}

	fmt.Println("Order preloading complete.")
}

// Simulate tick events for predefined symbols
func simulateTickEvents(redisClient *redis.Client) {
	symbols := []string{"AAPL", "GOOG", "TSLA", "AMZN", "MSFT"}
	ticker := time.NewTicker(10 * time.Millisecond) // Adjust frequency as needed
	defer ticker.Stop()

	fmt.Println("Simulating Tick Events...")

	for {
		<-ticker.C
		for _, symbol := range symbols {
			latestPrice := generateRandomPrice()
			tickEvent := TickEvent{
				Symbol:    symbol,
				Price:     latestPrice,
				Timestamp: time.Now().Unix(),
			}

			tickPayload, err := json.Marshal(tickEvent)
			if err != nil {
				fmt.Printf("Failed to marshal tick event: %v\n", err)
				continue
			}

			// Publish directly to Redis Pub/Sub
			err = redisClient.Publish(context.Background(), "tick_events", string(tickPayload)).Err()
			if err != nil {
				fmt.Printf("Failed to publish tick event: %v\n", err)
			}
		}
	}
}

// generateRandomPrice generates a random price for simulation purposes
func generateRandomPrice() float64 {
	return rand.Float64()*100 + 50 // Random price between 50 and 150
}
