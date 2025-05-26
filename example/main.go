package main

import (
	"fmt"
	"log"
	"time"

	"github.com/javanhut/easyhttp"
	"github.com/javanhut/easyjson"
)

func main() {
	fmt.Println("EasyHTTP Library Examples")
	fmt.Println("========================")

	// Example 1: Simple GET request
	fmt.Println("\n1. Simple GET Request:")
	resp, err := easyhttp.GET("https://httpbin.org/get")
	if err != nil {
		log.Printf("GET request failed: %v", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		if resp.OK() {
			jsonValue, _ := resp.JSONValue()
			url := jsonValue.Get("url").AsString()
			fmt.Printf("URL: %s\n", url)
		}
	}

	// Example 2: GET with query parameters
	fmt.Println("\n2. GET with Query Parameters:")
	resp, err = easyhttp.GET("https://httpbin.org/get", &easyhttp.RequestOptions{
		Params: map[string]string{
			"name": "John Doe",
			"age":  "30",
		},
	})
	if err != nil {
		log.Printf("GET with params failed: %v", err)
	} else if resp.OK() {
		jsonValue, _ := resp.JSONValue()
		args := jsonValue.Get("args")
		name := args.Get("name").AsString()
		age := args.Get("age").AsString()
		fmt.Printf("Query params - Name: %s, Age: %s\n", name, age)
	}

	// Example 3: POST with JSON using easyjson
	fmt.Println("\n3. POST with JSON (easyjson):")
	jsonObj := easyjson.NewObject()
	jsonObj.Set("name", "Alice Johnson")
	jsonObj.Set("email", "alice@example.com")
	jsonObj.Set("age", 28)

	resp, err = easyhttp.POST("https://httpbin.org/post", &easyhttp.RequestOptions{
		JSON: jsonObj,
	})
	if err != nil {
		log.Printf("POST request failed: %v", err)
	} else if resp.OK() {
		jsonValue, _ := resp.JSONValue()
		sentName := jsonValue.Q("json", "name").AsString()
		sentEmail := jsonValue.Q("json", "email").AsString()
		fmt.Printf("Sent data - Name: %s, Email: %s\n", sentName, sentEmail)
	}

	// Example 4: POST with regular map
	fmt.Println("\n4. POST with Regular Map:")
	userData := map[string]interface{}{
		"username": "johndoe",
		"active":   true,
		"score":    95.5,
	}

	resp, err = easyhttp.POST("https://httpbin.org/post", &easyhttp.RequestOptions{
		JSON: userData,
	})
	if err != nil {
		log.Printf("POST with map failed: %v", err)
	} else if resp.OK() {
		jsonValue, _ := resp.JSONValue()
		username := jsonValue.Q("json", "username").AsString()
		active := jsonValue.Q("json", "active").AsBool()
		score := jsonValue.Q("json", "score").AsFloat()
		fmt.Printf("Map data - Username: %s, Active: %v, Score: %.1f\n", username, active, score)
	}

	// Example 5: Using a configured client
	fmt.Println("\n5. Configured Client with Base URL:")
	client := easyhttp.New().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetHeaders(map[string]string{
			"User-Agent": "EasyHTTP-Example/1.0",
		}).
		SetTimeout(10 * time.Second)

	resp, err = client.GET("/posts/1")
	if err != nil {
		log.Printf("Client GET failed: %v", err)
	} else if resp.OK() {
		jsonValue, _ := resp.JSONValue()
		title := jsonValue.Get("title").AsString()
		userId := jsonValue.Get("userId").AsInt()
		fmt.Printf("Post Title: %s (User ID: %d)\n", title, userId)
	}

	// Example 6: Authentication
	fmt.Println("\n6. Bearer Token Authentication:")
	resp, err = easyhttp.GET("https://httpbin.org/bearer", &easyhttp.RequestOptions{
		Auth: &easyhttp.Auth{
			Token: "example-bearer-token-123",
		},
	})
	if err != nil {
		log.Printf("Auth request failed: %v", err)
	} else {
		fmt.Printf("Auth request status: %d\n", resp.StatusCode)
		if resp.OK() {
			jsonValue, _ := resp.JSONValue()
			authenticated := jsonValue.Get("authenticated").AsBool()
			token := jsonValue.Get("token").AsString()
			fmt.Printf("Authenticated: %v, Token: %s\n", authenticated, token)
		}
	}

	// Example 7: Error handling
	fmt.Println("\n7. Error Handling:")
	resp, err = easyhttp.GET("https://httpbin.org/status/404")
	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		if resp.OK() {
			fmt.Println("Request was successful")
		} else {
			fmt.Printf("Request failed with status: %d\n", resp.StatusCode)
		}
	}

	// Example 8: Complex JSON manipulation
	fmt.Println("\n8. Complex JSON with easyjson:")
	complexData := easyjson.NewObject()

	// Add user info
	user := easyjson.NewObject()
	user.Set("id", 12345)
	user.Set("name", "Emma Wilson")
	user.Set("email", "emma@example.com")

	// Add nested address
	address := easyjson.NewObject()
	address.Set("street", "456 Oak Avenue")
	address.Set("city", "San Francisco")
	address.Set("state", "CA")
	address.Set("zip", "94102")
	user.Set("address", address.Raw())

	// Add array of interests
	interests := easyjson.NewArrayFrom([]interface{}{
		"technology", "photography", "travel", "music",
	})
	user.Set("interests", interests.Raw())

	complexData.Set("user", user.Raw())
	complexData.Set("timestamp", time.Now().Format(time.RFC3339))
	complexData.Set("version", "1.0")

	resp, err = easyhttp.POST("https://httpbin.org/post", &easyhttp.RequestOptions{
		JSON: complexData,
	})
	if err != nil {
		log.Printf("Complex JSON request failed: %v", err)
	} else if resp.OK() {
		jsonValue, _ := resp.JSONValue()

		// Use query chaining for easy nested access
		userName := jsonValue.Q("json", "user", "name").AsString()
		userCity := jsonValue.Q("json", "user", "address", "city").AsString()
		firstInterest := jsonValue.Q("json", "user", "interests", 0).AsString()
		timestamp := jsonValue.Q("json", "timestamp").AsString()

		fmt.Printf("Complex data sent successfully!\n")
		fmt.Printf("User: %s from %s\n", userName, userCity)
		fmt.Printf("First interest: %s\n", firstInterest)
		fmt.Printf("Timestamp: %s\n", timestamp)

		// Count interests
		interestsArray := jsonValue.Q("json", "user", "interests")
		fmt.Printf("Total interests: %d\n", interestsArray.Len())
	}

	fmt.Println("\n========================")
	fmt.Println("All examples completed!")
	fmt.Println("This demonstrates the ease of use of EasyHTTP library")
	fmt.Println("with the power of easyjson for flexible JSON handling.")
}
