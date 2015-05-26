package main

import (
	"github.com/z-ray/alipaydemo/gateway"
	"github.com/z-ray/log"
	"net/http"
)

func main() {

	log.SetOutputLevel(log.Ldebug)
	addr := ":8080"
	log.Debugf("alipay service is running on %s", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	http.HandleFunc("/service/gateway.do", gateway.GatewayService)
	log.Fatal(http.ListenAndServe(addr, nil))

}
