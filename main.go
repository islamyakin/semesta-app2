package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	uid := r.Header.Get("X-UID")
	if uid == "" {
		uid = r.URL.Query().Get("uid")
	}
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Gagal mendapatkan hostname", http.StatusInternalServerError)
		return
	}

	url := "https://api.ipify.org?format=text"

	respCh := make(chan *http.Response)
	errCh := make(chan error)

	go func() {
		resp, err := http.Get(url)
		if err != nil {
			errCh <- err
			return
		}

		respCh <- resp
	}()

	select {
	case resp := <-respCh:
		handleResponse(w, resp, hostname, uid)
	case err := <-errCh:
		if strings.Contains(err.Error(), "no such host") {
			fmt.Fprintf(w, `<title>Semesta App 2</title>`)
			fmt.Fprintf(w, "<center><h1>App ini Berjalan di : %s</h1></center>", hostname)
			http.Error(w, "<h1><center>Namun anda tidak terhubung ke internet</center><h1>", http.StatusInternalServerError)
		} else {
			http.Error(w, "Kesalahan saat melakukan GET ke URL", http.StatusInternalServerError)
		}
		fmt.Println("Gagal melakukan GET ke URL:", err)
		return
	case <-time.After(10 * time.Second):
		fmt.Fprintf(w, "<center><h1>App ini Berjalan di : %s</h1></center>", hostname)
		fmt.Fprintf(w, "<center><h1>Namun belum berhasil GET ke API, Silahkan refresh kembali</h1></center>")
		fmt.Println("Waktu tunggu habis")
		return
	}
}

func handleResponse(w http.ResponseWriter, resp *http.Response, hostname string, uid string) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Gagal menutup response body:", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Gagal membaca respons", http.StatusInternalServerError)
		fmt.Println("Gagal membaca respons:", err)
		return
	}

	fmt.Fprintf(w, `<title>Semesta App 2</title>`)
	fmt.Fprintf(w, "<center><h1>App ini Berjalan di : %s</h1></center>", hostname)
	fmt.Fprintf(w, "<center><h1>Dengan IP Public : %s</h1></center>", body)
	fmt.Fprintf(w, `<center><img src="https://www.unger.dev/assets/200ok_logo_big.png" alt="200OK"></center>`)
	if uid != "" {
		fmt.Fprintf(w, "<center><h1>UUID: %s</h1></center>", uid)
	}
}

func main() {
	server := &http.Server{
		Addr:              ":3001",
		Handler:           http.HandlerFunc(handlerFunc),
		ReadHeaderTimeout: 3 * time.Second,
	}
	fmt.Println("server running on port 3001")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	fmt.Println(run())
}

func run() string {
	return "Hackathon SEMESTA - System Administrator"
}
