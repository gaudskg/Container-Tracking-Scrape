package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"io"
	"strings"
	"os"
)

// PageData holds the data to be displayed on the webpage
type PageData struct {
	Input  string
	Output string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		// Default to port 8080 if PORT environment variable is not set
		port = "8080"
	}
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/process", handleProcess)
	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//fmt.Println("Server is running on port 8080...")
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleIndex displays the index page with a form for input
func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

// handleProcess processes the form input and displays the result
func handleProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}


	input := r.FormValue("input")
	
	client := &http.Client{}
	json_string := fmt.Sprintf(`{"variables":{"{\"__wwtype\":\"f\",\"code\":\"variables['465276e0-a224-4a28-bf2c-098ddca8ca4b']\"}":"62bd38e5-6c23-4f5b-839e-57086cc492f1","{\"__wwtype\":\"f\",\"code\":\"variables['148f038e-244a-44c4-bcb9-750a64e74167-value']||variables['654fb0ae-1573-42f0-abce-648701780a13-value']||variables['3300cf30-27b8-4c61-9de6-b7d1bd59af48']\"}":"%s","{\"__wwtype\":\"f\",\"code\":\"variables['1bb4ebcb-2021-4770-99f3-c7a00de242ca-value']||\\\"\\\"\"}":"","{\"__wwtype\":\"f\",\"code\":\"variables['a03e8a5d-fd4b-43b8-b206-1f95f441753c-value']||\\\"\\\"\"}":""}}`,input)
	var data = strings.NewReader(json_string)
	req, err := http.NewRequest("POST", "https://tools.sinay.ai/ww/cms_data_sets/97105c3b-3a4b-41e4-a2a9-f102ad05c44f/fetch?limit=100&offset=0", data)
	if err != nil {
	    log.Fatal(err)
	}
	req.Header.Set("authority", "tools.sinay.ai")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en;q=0.7")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://tools.sinay.ai")
	req.Header.Set("referer", "https://tools.sinay.ai/container-tracking/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Brave";v="121", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
	    log.Fatal(err)
	}
	defer resp.Body.Close()
	output, err := io.ReadAll(resp.Body)
	if err != nil {
	    log.Fatal(err)
	}

	// Printing the Output
	//fmt.Printf("%s\n", output)


	output_data := PageData{
		Input:  input,
		Output: string(output),
	}

	tmpl := template.Must(template.ParseFiles("result.html"))
	tmpl.Execute(w, output_data)
}

