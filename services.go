package main

import (
	"crypto/md5"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

)



// Services structure is used to contain the available html handlers
type Services struct{}

// ServeHTTP serve method
func (s *Services) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.URL.Path == "/" {
		_, err = s.HandleUpload(w, r)
	} else {
		err = ErrPageNotFound
	}

	if err != nil && err == ErrPageNotFound {
		http.NotFoundHandler().ServeHTTP(w, r)
	} else if err != nil {
		log.Println(err)
		http.Error(w, "unable to complete request", http.StatusInternalServerError)
		return
	}
}

// HandleUpload handle file upload
func (s *Services) HandleUpload(w http.ResponseWriter, r *http.Request) (*Page, error) {
	title := r.URL.Path[len("/"):]
	p := &Page{Title: title}

	// Set the upload token
	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))

	p.Token = token

	if r.Method == http.MethodPost {
		message, err := s.HandlePost(r)
		if err != nil {
			return nil, err
		}
		p.Message = message
	} else if r.Method != http.MethodGet {
		return nil, ErrPageNotFound
	}

	// Run template
	t, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		return nil, err
	}
	t.Execute(w, p)

	return p, nil
}

// HandlePost handle posts
func (s *Services) HandlePost(r *http.Request) (string, error) {

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return "Error Retrieving the File", err
	}
	defer file.Close()
	ext := filepath.Ext(handler.Filename)
	if ext != ".pcap" {
		return "Invalid file extension. File extension should be - \".pcap\"", nil
	}
	err = s.ProcessPCAP(handler.Filename)
	if err != nil {
		return "Error while processing .pcap file", err
	}

	message := "done"

	return message, nil
}

func (s *Services) ProcessPCAP(filename string) error {
	log.Printf("Begin reading '%s'...", filename)

	if handle, err := pcap.OpenOffline(filename); err != nil {
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet)  // Do something with a packet here.
		}
	}

	return nil
}

//getReader get csv reader from file
func handlePacket(handlePacket gopacket.Packet) error {

	log.Println(handlePacket)
	return nil
}