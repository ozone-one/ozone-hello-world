package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.elastic.co/apm/module/apmgorilla"
)

func hello(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/"))
	defer span.Finish()
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() > 0.80 {
		fmt.Println("500 - error Something bad happened!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - error Something bad happened!"))
		span.Finish(tracer.WithError(fmt.Errorf("500 - Something bad happened!")))
	} else {
		fmt.Println("200 - Success")
		fmt.Fprintf(w, `{"route":"/hello","response":"hello world main branch"}`)
	}
}

func empty(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpan("web.request", tracer.ResourceName("/hello"))
	defer span.Finish()
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() > 0.98 {
		fmt.Println("500 - error Something bad happened!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - error Something bad happened!"))
		span.Finish(tracer.WithError(fmt.Errorf("500 - Something bad happened!")))
	} else {
		fmt.Println("200 - Success")
		fmt.Fprintf(w, `{"route":"/","branch":"main"}`)
	}
}

func main() {

	apmProvider := os.Getenv("APM_TYPE")
	r := mux.NewRouter()
	switch apmProvider {
	case "newrelic":
		fmt.Println("Initializing new relic apm agent")
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("in2tive-go-hello-world"),
			newrelic.ConfigLicense("edc3adcbb4eec86331b9f44890108b628554NRAL"),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			panic(err)
		}
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/", empty))
		r.HandleFunc(newrelic.WrapHandleFunc(app, "/hello", hello))

	case "elastic":
		fmt.Println("Initializing elastic apm agent")
		r.Use(apmgorilla.Middleware())
		r.HandleFunc("/", empty)
		r.HandleFunc("/hello", hello)

	case "datadog":
		fmt.Println("Initializing datadog apm agent")
		tracer.Start()
		defer tracer.Stop()
		r.HandleFunc("/", empty)
		r.HandleFunc("/hello", hello)
	default:
		fmt.Println("No apm agent initialized")
		r.HandleFunc("/", empty)
		r.HandleFunc("/hello", hello)
	}
	fmt.Print(http.ListenAndServe(":3000", r))
}
