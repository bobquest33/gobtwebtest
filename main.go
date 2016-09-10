package main

import (
        "log"
		"os"
        "github.com/tarm/serial"
		"net/http"
		
)

var (
	http_port string
	serial_port string
	serial_porto *serial.Port
)

func init() {
	http_port = os.Args[1]
	serial_port = os.Args[2]
	
}

func send_serial(cmd string){
 
	log.Println(cmd)
    _, err := serial_porto.Write([]byte(cmd))
	log.Println(err)
    if err != nil {
            log.Fatal(err)
    }

    log.Println(cmd)
	buf := make([]byte, 128)
    n, err := serial_porto.Read(buf)
	log.Println(err)

    if err != nil {
            log.Fatal(err)
    }
    log.Printf("%q", buf[:n])

}

func handleSwitch(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	data := r.PostFormValue("switch")
	log.Println(data)
	if data == "On" {
		send_serial("a")
		w.Write([]byte("On"))
	} else if data == "Off" {
		send_serial("b")
		w.Write([]byte("Off"))
		
	} else {
		w.Write([]byte("Error"))
	}	
}

func main() {
	var err error
    c := &serial.Config{Name: serial_port, Baud: 9600}
    serial_porto, err = serial.OpenPort(c)
    if err != nil {
            log.Fatal(err)
    }
	defer serial_porto.Close()
	
	http.Handle("/",http.FileServer(FS(false)))
	http.HandleFunc("/api",handleSwitch)
	
	log.Println("Listening...")
	http.ListenAndServe(http_port,nil)

}