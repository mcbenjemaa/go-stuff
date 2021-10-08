package rest

import (
	"fmt"

	resty "github.com/go-resty/resty/v2"
	"github.com/mcbenjemaa/go-stuff/internal/album"
)

const baseUrl = "http://localhost:8080/albums"

// HTTP GET client example
func HttpGet() {
	// Create a Resty Client
	client := resty.New()

	var albums []album.Album

	resp, err := client.R().
		EnableTrace().
		SetResult(&albums). // Json Marshal
		Get(baseUrl)

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Printf("  albums      %#v:", albums)
	fmt.Println()

	// -- Explore trace info
	// fmt.Println("Request Trace Info:")
	// ti := resp.Request.TraceInfo()
	// fmt.Println("  DNSLookup     :", ti.DNSLookup)
	// fmt.Println("  ConnTime      :", ti.ConnTime)
	// fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	// fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	// fmt.Println("  ServerTime    :", ti.ServerTime)
	// fmt.Println("  ResponseTime  :", ti.ResponseTime)
	// fmt.Println("  TotalTime     :", ti.TotalTime)
	// fmt.Println("  IsConnReused  :", ti.IsConnReused)
	// fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	// fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	// fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	// fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

// HTTP GET client example
func HttpError() {

	// Create a Resty Client
	client := resty.New()

	var album album.Album
	var authError AuthError

	resp, err := client.R().
		SetResult(&album). // Json Marshal
		SetError(&authError).
		Get(fmt.Sprintf("%v/%v", baseUrl, "8")) // pass parameter

	if err != nil {
		fmt.Println(fmt.Errorf("Error %w:", err))
	} else if resp.StatusCode() != 200 {
		fmt.Println(fmt.Errorf("Error: status %v,  %v", resp.StatusCode(), authError))
	} else {
		fmt.Println(fmt.Sprintf("Album: %#v", album))
	}

}

type AuthSuccess struct {
}

type AuthError struct {
	message string
}

// HTTP GET client example
func HttpPost() {
	// Create a Resty Client
	client := resty.New()

	// POST JSON string
	// No need to set content type, if you have client level setting
	// resp, err := client.R().
	// SetHeader("Content-Type", "application/json").
	// SetBody(`{"username":"testuser", "password":"testpass"}`).
	// SetResult(&AuthSuccess{}).    // or SetResult(AuthSuccess{}).
	// Post("http://lo.com/login")
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Error: POST request not working %w", err))
	// }

	// POST Struct, default is JSON content type. No need to set one
	_, err := client.R().
		SetBody(album.Album{ID: "4", Title: "K8s", Artist: "Med Bj", Price: 100.4}).
		SetResult(&AuthSuccess{}).
		SetError(&AuthError{}).
		Post(baseUrl)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: POST request not working %w", err))
	}

}

func HttpPut() {
	// Note: This is one sample of PUT method usage, refer POST for more combination

	// Create a Resty Client
	//client := resty.New()

	// Request goes as JSON content type
	// No need to set auth token, error, if you have client level settings
	// resp, err := client.R().
	// 	SetBody(Article{
	// 		Title: "go-resty",
	// 		Content: "This is my article content, oh ya!",
	// 		Author: "Jeevanandam M",
	// 		Tags: []string{"article", "sample", "resty"},
	// 	}).
	// 	SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
	// 	//SetError(&fmt.Errorf("Is not working %w")).       // or SetError(Error{}).
	// 	Put("https://myapp.com/article/1234")
}
