package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var port = flag.Int("port", 8080, "Port number to use")
var cert = flag.String("cert", "", "Certificate file")
var key = flag.String("key", "", "Private key file")
var fs = flag.String("fs", "", "Location of file system")
var timeout = flag.Int("timeout", 60, "Timeout")
var test = flag.Bool("test", false, "Use internal template")

func home(w http.ResponseWriter, r *http.Request) {
	scheme := "https://"
	if r.TLS == nil {
		scheme = "http://"
	}
	log.Println(fmt.Sprintf("%s%s", scheme, r.Host))
	data := struct {
		BaseURL string
		Title   string
	}{
		BaseURL: fmt.Sprintf("%s%s", scheme, r.Host),
		Title:   "Webserver",
	}
	homeTemplate.Execute(w, data)
}

func main() {
	flag.Parse()

	if *fs == "" {
		fmt.Println("Usage SimpleHTTPServer -port <port number> -cert <cert file> -key <privkey file> -fs <home folder>")
		fmt.Println("No file system specified, will serve internal index.html file")
		*fs = "./"
		*test = true
	}

	fsys := http.FileServer(http.Dir(*fs))
	if *test {
		http.HandleFunc("/", home)
		log.Println("Test on")
		fsys = nil
	}

	if *cert != "" && *key != "" {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		srv := &http.Server{
			Addr:           fmt.Sprintf(":%d", *port),
			Handler:        fsys,
			ReadTimeout:    time.Duration(*timeout) * time.Second,
			WriteTimeout:   time.Duration(*timeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      cfg,
			TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}

		log.Printf("Webserver listen port :%d, with TLS", *port)
		err := srv.ListenAndServeTLS(*cert, *key)
		if err != nil {
			log.Fatal(fmt.Sprintf("ListenAndServeTLS: %s", err.Error()))
		}
	} else {
		log.Printf("Webserver listen port :%d, with no TLS", *port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), fsys)
		if err != nil {
			log.Fatal(fmt.Sprintf("ListenAndServe: %s", err.Error()))
		}
	}
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Webserver</title>
    <!-- CSS only -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">
    <!-- JavaScript Bundle with Popper -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js" integrity="sha384-ygbV9kiqUc6oa4msXn9868pTtWMgiQaeYH7/t7LECLbyPA2x65Kgf80OJFdroafW" crossorigin="anonymous"></script>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
body {
  background: linear-gradient(-45deg, #ee7752, #e73c7e, #23a6d5, #23d5ab, #1d251f);
  background-size: 400 400%;
  animation: gradient 15s ease infinite;
  height: 100vh;
}

@keyframes gradient {
	0% {
	  background-position: 0% 50%;
	}
	50% {
	  background-position: 100% 50%;
	}
	100% {
	  background-position: 0% 50%;
	}
  }
  </style>
  </head>
<body>
    <div class="d-flex flex-column justify-content-center w-100 h-100">
    
      <div class="d-flex flex-column justify-content-center align-items-center">
        <h1 class="fw-light text-white m-0">Simple HTTP/HTTPS Webserver</h1>
        <div class="btn-group my-5">
          <a href="https://albievan.com/SimpleHTTPServer.tar.gz" class="btn btn-outline-light" aria-current="page"><i class="fas fa-file-download me-2"></i> SOURCE CODE</a>
          <a href="https://albievan.com/SimpleHTTPServer.MacOS.tar.gz" class="btn btn-outline-light" aria-current="page"><i class="fas fa-file-download me-2"></i> MacOS Binary</a>
        </div>
        <a href="https://albievan.com/" class="text-decoration-none">
          <h5 class="fw-light text-white m-0"> Albie van der Merwe </h5>
        </a>
      </div>
    </div>
    </div>
  </body>
</html>
`))
