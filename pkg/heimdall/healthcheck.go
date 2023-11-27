package heimdall

import "net/http"

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if IsHealthy {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			Fatal(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("NOT OK"))
		if err != nil {
			Fatal(err.Error())
		}
	}
}

func StartHealthCheckServer() {
	http.HandleFunc("/health", healthCheckHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		Fatal(err.Error())
	}
}
