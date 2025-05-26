package easyhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/javanhut/easyjson"
)

// Test server handlers
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.String(),
		"headers": map[string]string{
			"user-agent":    r.Header.Get("User-Agent"),
			"authorization": r.Header.Get("Authorization"),
			"content-type":  r.Header.Get("Content-Type"),
		},
	}

	if r.Body != nil {
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		if len(body) > 0 {
			var jsonBody interface{}
			if json.Unmarshal(body, &jsonBody) == nil {
				response["json"] = jsonBody
			} else {
				response["data"] = string(body)
			}
		}
	}

	json.NewEncoder(w).Encode(response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		code = "200"
	}

	switch code {
	case "200":
		w.WriteHeader(200)
	case "404":
		w.WriteHeader(404)
	case "500":
		w.WriteHeader(500)
	default:
		w.WriteHeader(200)
	}

	fmt.Fprintf(w, `{"status": %s}`, code)
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	delay := r.URL.Query().Get("seconds")
	if delay == "" {
		delay = "1"
	}

	if delay == "2" {
		time.Sleep(2 * time.Second)
	} else {
		time.Sleep(1 * time.Second)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"delayed": "%s seconds"}`, delay)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")
	if count == "" || count == "0" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"redirected": true}`)
		return
	}

	http.Redirect(w, r, "/redirect?count=0", http.StatusFound)
}

func createTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/status", statusHandler)
	mux.HandleFunc("/delay", delayHandler)
	mux.HandleFunc("/redirect", redirectHandler)

	return httptest.NewServer(mux)
}

// Test Response methods
func TestResponse_Text(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	text := resp.Text()
	if text == "" {
		t.Error("Expected non-empty text response")
	}

	if !strings.Contains(text, "method") {
		t.Error("Expected response to contain 'method'")
	}
}

func TestResponse_JSON(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	var result map[string]interface{}
	if err := resp.JSON(&result); err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	if result["method"] != "GET" {
		t.Errorf("Expected method=GET, got %v", result["method"])
	}
}

func TestResponse_JSONValue(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSONValue parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "GET" {
		t.Errorf("Expected method=GET, got %s", method)
	}

	// Test fluent query chaining
	userAgent := jsonValue.Q("headers", "user-agent").AsString()
	if userAgent == "" {
		t.Error("Expected non-empty user-agent")
	}
}

func TestResponse_OK(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Test successful response
	resp, err := GET(server.URL + "/status?code=200")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	if !resp.OK() {
		t.Error("Expected OK() to return true for 200 status")
	}

	// Test error response
	resp, err = GET(server.URL + "/status?code=404")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	if resp.OK() {
		t.Error("Expected OK() to return false for 404 status")
	}
}

func TestResponse_Bytes(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	bytes := resp.Bytes()
	if len(bytes) == 0 {
		t.Error("Expected non-empty byte response")
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Errorf("Response bytes are not valid JSON: %v", err)
	}
}

// Test HTTP methods
func TestGET(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "GET" {
		t.Errorf("Expected method=GET, got %s", method)
	}
}

func TestPOST(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	testData := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}

	resp, err := POST(server.URL+"/json", &RequestOptions{
		JSON: testData,
	})
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "POST" {
		t.Errorf("Expected method=POST, got %s", method)
	}

	// Check if JSON data was sent correctly
	sentName := jsonValue.Q("json", "name").AsString()
	if sentName != "John Doe" {
		t.Errorf("Expected name=John Doe, got %s", sentName)
	}

	sentAge := jsonValue.Q("json", "age").AsInt()
	if sentAge != 30 {
		t.Errorf("Expected age=30, got %d", sentAge)
	}
}

func TestPOST_WithEasyJSONValue(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Create JSON using easyjson
	jsonObj := easyjson.NewObject()
	jsonObj.Set("name", "Alice Johnson")
	jsonObj.Set("email", "alice@example.com")
	jsonObj.Set("active", true)

	resp, err := POST(server.URL+"/json", &RequestOptions{
		JSON: jsonObj,
	})
	if err != nil {
		t.Fatalf("POST with easyjson failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	// Verify the data was sent correctly
	sentName := jsonValue.Q("json", "name").AsString()
	if sentName != "Alice Johnson" {
		t.Errorf("Expected name=Alice Johnson, got %s", sentName)
	}

	sentEmail := jsonValue.Q("json", "email").AsString()
	if sentEmail != "alice@example.com" {
		t.Errorf("Expected email=alice@example.com, got %s", sentEmail)
	}

	sentActive := jsonValue.Q("json", "active").AsBool()
	if !sentActive {
		t.Error("Expected active=true")
	}
}

func TestPUT(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := PUT(server.URL + "/json")
	if err != nil {
		t.Fatalf("PUT failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "PUT" {
		t.Errorf("Expected method=PUT, got %s", method)
	}
}

func TestDELETE(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := DELETE(server.URL + "/json")
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "DELETE" {
		t.Errorf("Expected method=DELETE, got %s", method)
	}
}

func TestPATCH(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := PATCH(server.URL + "/json")
	if err != nil {
		t.Fatalf("PATCH failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "PATCH" {
		t.Errorf("Expected method=PATCH, got %s", method)
	}
}

func TestHEAD(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := HEAD(server.URL + "/json")
	if err != nil {
		t.Fatalf("HEAD failed: %v", err)
	}

	// HEAD responses typically have empty bodies
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// Test Client functionality
func TestClient_SetBaseURL(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	client := New().SetBaseURL(server.URL)

	resp, err := client.GET("/json")
	if err != nil {
		t.Fatalf("GET with base URL failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	url := jsonValue.Get("url").AsString()
	if !strings.Contains(url, "/json") {
		t.Errorf("Expected URL to contain /json, got %s", url)
	}
}

func TestClient_SetHeaders(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	client := New().SetHeaders(map[string]string{
		"User-Agent":      "TestAgent/1.0",
		"X-Custom-Header": "test-value",
	})

	resp, err := client.GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET with headers failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	userAgent := jsonValue.Q("headers", "user-agent").AsString()
	if userAgent != "TestAgent/1.0" {
		t.Errorf("Expected User-Agent=TestAgent/1.0, got %s", userAgent)
	}
}

func TestClient_SetAuth(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Test Bearer token
	client := New().SetAuth(&Auth{
		Token: "test-token-123",
	})

	resp, err := client.GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET with auth failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	auth := jsonValue.Q("headers", "authorization").AsString()
	expected := "Bearer test-token-123"
	if auth != expected {
		t.Errorf("Expected Authorization=%s, got %s", expected, auth)
	}
}

func TestClient_SetTimeout(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	client := New().SetTimeout(500 * time.Millisecond)

	// This should timeout
	_, err := client.GET(server.URL + "/delay?seconds=2")
	if err == nil {
		t.Error("Expected timeout error, but request succeeded")
	}
}

// Test RequestOptions
func TestRequestOptions_Params(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL+"/json", &RequestOptions{
		Params: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})
	if err != nil {
		t.Fatalf("GET with params failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	url := jsonValue.Get("url").AsString()
	if !strings.Contains(url, "key1=value1") {
		t.Errorf("Expected URL to contain key1=value1, got %s", url)
	}
	if !strings.Contains(url, "key2=value2") {
		t.Errorf("Expected URL to contain key2=value2, got %s", url)
	}
}

func TestRequestOptions_Headers(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL+"/json", &RequestOptions{
		Headers: map[string]string{
			"X-Test-Header": "test-value",
		},
	})
	if err != nil {
		t.Fatalf("GET with headers failed: %v", err)
	}

	// The test server doesn't echo back all headers, but we can verify the request was made
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestRequestOptions_Data(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	testData := "key1=value1&key2=value2"

	resp, err := POST(server.URL+"/json", &RequestOptions{
		Data: testData,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})
	if err != nil {
		t.Fatalf("POST with data failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	data := jsonValue.Get("data").AsString()
	if data != testData {
		t.Errorf("Expected data=%s, got %s", testData, data)
	}
}

func TestRequestOptions_Auth(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Test request-specific auth overriding client auth
	client := New().SetAuth(&Auth{Token: "default-token"})

	resp, err := client.GET(server.URL+"/json", &RequestOptions{
		Auth: &Auth{Token: "override-token"},
	})
	if err != nil {
		t.Fatalf("GET with override auth failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	auth := jsonValue.Q("headers", "authorization").AsString()
	expected := "Bearer override-token"
	if auth != expected {
		t.Errorf("Expected Authorization=%s, got %s", expected, auth)
	}
}

func TestRequestOptions_Timeout(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// This should timeout
	_, err := GET(server.URL+"/delay?seconds=2", &RequestOptions{
		Timeout: 500 * time.Millisecond,
	})
	if err == nil {
		t.Error("Expected timeout error, but request succeeded")
	}
}

func TestRequestOptions_AllowRedirect(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Test with redirects disabled - should get redirect status instead of following
	resp, err := GET(server.URL+"/redirect?count=1", &RequestOptions{
		AllowRedirect: false,
	})
	if err != nil {
		t.Fatalf("GET with no redirect failed: %v", err)
	}

	// Should get redirect status code instead of final destination
	if resp.StatusCode != 302 {
		t.Errorf("Expected status 302 (redirect), got %d", resp.StatusCode)
	}
}

// Test Authentication
func TestAuth_BasicAuth(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	client := New().SetAuth(&Auth{
		Username: "testuser",
		Password: "testpass",
	})

	resp, err := client.GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET with basic auth failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	// Basic auth sets Authorization header, which our test server captures
	auth := jsonValue.Q("headers", "authorization").AsString()
	if !strings.HasPrefix(auth, "Basic ") {
		t.Errorf("Expected Basic auth header, got %s", auth)
	}
}

func TestAuth_BearerToken(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL+"/json", &RequestOptions{
		Auth: &Auth{
			Token: "test-bearer-token",
		},
	})
	if err != nil {
		t.Fatalf("GET with bearer token failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	auth := jsonValue.Q("headers", "authorization").AsString()
	expected := "Bearer test-bearer-token"
	if auth != expected {
		t.Errorf("Expected Authorization=%s, got %s", expected, auth)
	}
}

// Test JSON handling with different data types
func TestJSON_ComplexData(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Create complex JSON structure using easyjson
	complexData := easyjson.NewObject()

	// Add user info
	user := easyjson.NewObject()
	user.Set("id", 12345)
	user.Set("name", "Emma Watson")
	user.Set("active", true)
	user.Set("score", 95.5)

	// Add nested address
	address := easyjson.NewObject()
	address.Set("street", "123 Main St")
	address.Set("city", "New York")
	address.Set("zip", "10001")
	user.Set("address", address.Raw())

	// Add array of tags
	tags := easyjson.NewArrayFrom([]interface{}{"vip", "premium", "verified"})
	user.Set("tags", tags.Raw())

	complexData.Set("user", user.Raw())
	complexData.Set("timestamp", "2024-01-15T10:30:00Z")

	resp, err := POST(server.URL+"/json", &RequestOptions{
		JSON: complexData,
	})
	if err != nil {
		t.Fatalf("POST with complex JSON failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	// Test query chaining for nested access
	userName := jsonValue.Q("json", "user", "name").AsString()
	if userName != "Emma Watson" {
		t.Errorf("Expected name=Emma Watson, got %s", userName)
	}

	userCity := jsonValue.Q("json", "user", "address", "city").AsString()
	if userCity != "New York" {
		t.Errorf("Expected city=New York, got %s", userCity)
	}

	firstTag := jsonValue.Q("json", "user", "tags", 0).AsString()
	if firstTag != "vip" {
		t.Errorf("Expected first tag=vip, got %s", firstTag)
	}

	userScore := jsonValue.Q("json", "user", "score").AsFloat()
	if userScore != 95.5 {
		t.Errorf("Expected score=95.5, got %f", userScore)
	}

	timestamp := jsonValue.Q("json", "timestamp").AsString()
	if timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected timestamp=2024-01-15T10:30:00Z, got %s", timestamp)
	}
}

func TestJSON_StringInput(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	jsonString := `{"message": "Hello World", "count": 42}`

	resp, err := POST(server.URL+"/json", &RequestOptions{
		JSON: jsonString,
	})
	if err != nil {
		t.Fatalf("POST with JSON string failed: %v", err)
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	message := jsonValue.Q("json", "message").AsString()
	if message != "Hello World" {
		t.Errorf("Expected message=Hello World, got %s", message)
	}

	count := jsonValue.Q("json", "count").AsInt()
	if count != 42 {
		t.Errorf("Expected count=42, got %d", count)
	}
}

// Test error conditions
func TestErrorHandling_InvalidURL(t *testing.T) {
	_, err := GET("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, but got none")
	}
}

func TestErrorHandling_InvalidJSON(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	_, err := POST(server.URL+"/json", &RequestOptions{
		JSON: make(chan int), // channels can't be marshaled to JSON
	})
	if err == nil {
		t.Error("Expected error for invalid JSON, but got none")
	}
}

func TestErrorHandling_InvalidJSONString(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	_, err := POST(server.URL+"/json", &RequestOptions{
		JSON: `{"invalid": json}`, // malformed JSON string
	})
	if err == nil {
		t.Error("Expected error for invalid JSON string, but got none")
	}
}

func TestErrorHandling_NetworkError(t *testing.T) {
	// Try to connect to a non-existent server
	_, err := GET("http://localhost:99999/nonexistent")
	if err == nil {
		t.Error("Expected network error, but got none")
	}
}

// Test Response body reading multiple times
func TestResponse_MultipleReads(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	// Read as text first
	text1 := resp.Text()
	if text1 == "" {
		t.Error("Expected non-empty text response")
	}

	// Read as text again - should return same content
	text2 := resp.Text()
	if text1 != text2 {
		t.Error("Multiple Text() calls should return the same content")
	}

	// Read as bytes - should also work
	bytes := resp.Bytes()
	if len(bytes) == 0 {
		t.Error("Expected non-empty byte response")
	}

	// Read as JSON - should also work
	var result map[string]interface{}
	if err := resp.JSON(&result); err != nil {
		t.Errorf("JSON parsing failed after multiple reads: %v", err)
	}

	// Read as JSONValue - should also work
	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Errorf("JSONValue parsing failed after multiple reads: %v", err)
	}

	method := jsonValue.Get("method").AsString()
	if method != "GET" {
		t.Errorf("Expected method=GET, got %s", method)
	}
}

// Test method chaining
func TestClient_MethodChaining(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	// Test that all Set methods return the client for chaining
	client := New().
		SetBaseURL(server.URL).
		SetTimeout(5 * time.Second).
		SetHeaders(map[string]string{
			"User-Agent": "ChainTest/1.0",
		}).
		SetAuth(&Auth{Token: "chain-token"})

	resp, err := client.GET("/json")
	if err != nil {
		t.Fatalf("Chained client GET failed: %v", err)
	}

	if !resp.OK() {
		t.Error("Expected successful response from chained client")
	}

	jsonValue, err := resp.JSONValue()
	if err != nil {
		t.Fatalf("JSON parsing failed: %v", err)
	}

	// Verify all settings were applied
	userAgent := jsonValue.Q("headers", "user-agent").AsString()
	if userAgent != "ChainTest/1.0" {
		t.Errorf("Expected User-Agent=ChainTest/1.0, got %s", userAgent)
	}

	auth := jsonValue.Q("headers", "authorization").AsString()
	if auth != "Bearer chain-token" {
		t.Errorf("Expected Authorization=Bearer chain-token, got %s", auth)
	}
}

// Test URL building
func TestClient_URLBuilding(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	client := New().SetBaseURL(server.URL)

	testCases := []struct {
		path     string
		expected string
	}{
		{"/json", "/json"},
		{"json", "/json"},
	}

	for _, tc := range testCases {
		resp, err := client.GET(tc.path)
		if err != nil {
			t.Fatalf("GET %s failed: %v", tc.path, err)
		}

		if !resp.OK() {
			t.Fatalf("GET %s returned status %d", tc.path, resp.StatusCode)
		}

		jsonValue, err := resp.JSONValue()
		if err != nil {
			t.Fatalf("JSON parsing failed for %s: %v", tc.path, err)
		}

		url := jsonValue.Get("url").AsString()
		if !strings.Contains(url, tc.expected) {
			t.Errorf("Expected URL to contain %s, got %s", tc.expected, url)
		}
	}

	// Test URL building without making requests (to avoid 404s)
	testURLs := []struct {
		baseURL  string
		path     string
		expected string
	}{
		{"https://api.example.com", "/users", "https://api.example.com/users"},
		{"https://api.example.com/", "/users", "https://api.example.com/users"},
		{"https://api.example.com", "users", "https://api.example.com/users"},
		{"https://api.example.com/", "users", "https://api.example.com/users"},
	}

	for _, tc := range testURLs {
		client := New().SetBaseURL(tc.baseURL)
		// Test the internal URL building by checking what would be built
		fullURL := client.buildURL(tc.path, nil)
		if fullURL != tc.expected {
			t.Errorf("Expected URL %s, got %s", tc.expected, fullURL)
		}
	}
}

// Benchmark tests
func BenchmarkGET(b *testing.B) {
	server := createTestServer()
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := GET(server.URL + "/json")
		if err != nil {
			b.Fatalf("GET failed: %v", err)
		}
		_ = resp.Text()
	}
}

func BenchmarkPOST_JSON(b *testing.B) {
	server := createTestServer()
	defer server.Close()

	testData := map[string]interface{}{
		"name": "Benchmark Test",
		"id":   12345,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := POST(server.URL+"/json", &RequestOptions{
			JSON: testData,
		})
		if err != nil {
			b.Fatalf("POST failed: %v", err)
		}
		_ = resp.Text()
	}
}

func BenchmarkJSONParsing_Traditional(b *testing.B) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		b.Fatalf("GET failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := resp.JSON(&result); err != nil {
			b.Fatalf("JSON parsing failed: %v", err)
		}
	}
}

func BenchmarkJSONParsing_EasyJSON(b *testing.B) {
	server := createTestServer()
	defer server.Close()

	resp, err := GET(server.URL + "/json")
	if err != nil {
		b.Fatalf("GET failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonValue, err := resp.JSONValue()
		if err != nil {
			b.Fatalf("JSONValue parsing failed: %v", err)
		}
		_ = jsonValue.Get("method").AsString()
	}
}

func BenchmarkClient_Reuse(b *testing.B) {
	server := createTestServer()
	defer server.Close()

	client := New().SetBaseURL(server.URL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.GET("/json")
		if err != nil {
			b.Fatalf("GET failed: %v", err)
		}
		_ = resp.Text()
	}
}

// Test concurrent requests
func TestConcurrentRequests(t *testing.T) {
	server := createTestServer()
	defer server.Close()

	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			resp, err := GET(server.URL + "/json")
			if err != nil {
				results <- fmt.Errorf("request %d failed: %v", id, err)
				return
			}

			if !resp.OK() {
				results <- fmt.Errorf("request %d got status %d", id, resp.StatusCode)
				return
			}

			jsonValue, err := resp.JSONValue()
			if err != nil {
				results <- fmt.Errorf("request %d JSON parsing failed: %v", id, err)
				return
			}

			method := jsonValue.Get("method").AsString()
			if method != "GET" {
				results <- fmt.Errorf("request %d expected method=GET, got %s", id, method)
				return
			}

			results <- nil
		}(i)
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent request failed: %v", err)
		}
	}
}
