package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		defer resp.Body.Close()
		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Gagal membaca respons", http.StatusInternalServerError)
			fmt.Println("Gagal membaca respons:", err)
			return
		}

		fmt.Fprintf(w, "<center><h1>App ini Berjalan di : %s</h1></center>", hostname)
		fmt.Fprintf(w, "<center><h1>Dengan IP Public : %s</h1></center>", ip)
	case err := <-errCh:
		if strings.Contains(err.Error(), "no such host") {
			http.Error(w, "Tidak ada koneksi internet", http.StatusInternalServerError)
		} else {
			http.Error(w, "Kesalahan saat melakukan GET ke URL", http.StatusInternalServerError)
		}
		fmt.Println("Gagal melakukan GET ke URL:", err)
		return
	case <-time.After(3 * time.Second):
		http.Error(w, "Waktu tunggu habis", http.StatusInternalServerError)
		fmt.Println("Waktu tunggu habis")
		return
	}

}

func main() {
	http.HandleFunc("/aboutus", handlerFunc)
	http.ListenAndServe(":3001", nil)
}