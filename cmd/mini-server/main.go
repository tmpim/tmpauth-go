package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/tmpim/tmpauth-go"
)

func main() {
	tmpauthInstances := make(map[string]*tmpauth.Tmpauth)

	http.HandleFunc("/parse-wrapped-auth-jwt", func(w http.ResponseWriter, r *http.Request) {
		configID := r.Header.Get(tmpauth.ConfigIDHeader)
		if configID == "" {
			log.Println("missing config ID")
			http.Error(w, "missing config ID", http.StatusBadRequest)
			return
		}

		ta, ok := tmpauthInstances[configID]
		if !ok {
			log.Println("invalid config ID:", configID)
			http.Error(w, "invalid config ID", http.StatusPreconditionFailed)
			return
		}

		token, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cachedToken, err := ta.ParseWrappedAuthJWT(string(token))
		if err != nil {
			log.Println("error parsing token:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(cachedToken)
		if err != nil {
			log.Println("error encoding response:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/start-auth", func(w http.ResponseWriter, r *http.Request) {
		configID := r.Header.Get(tmpauth.ConfigIDHeader)
		if configID == "" {
			log.Println("missing config ID")
			http.Error(w, "missing config ID", http.StatusBadRequest)
			return
		}

		requestURI := r.Header.Get(tmpauth.RequestURIHeader)
		if requestURI == "" {
			log.Println("missing request URI")
			http.Error(w, "missing request URI", http.StatusBadRequest)
			return
		}

		host := r.Header.Get(tmpauth.HostHeader)
		if host == "" {
			log.Println("missing host")
			http.Error(w, "missing host", http.StatusBadRequest)
			return
		}

		var err error
		r.URL, err = url.ParseRequestURI(requestURI)
		if err != nil {
			log.Println("error parsing request URI:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Header.Set("Host", host)

		ta, ok := tmpauthInstances[configID]
		if !ok {
			log.Println("invalid config ID:", configID)
			http.Error(w, "invalid config ID", http.StatusPreconditionFailed)
			return
		}

		code, err := ta.StartAuth(w, r)
		if err != nil {
			if code == 0 {
				code = http.StatusInternalServerError
			}
			w.WriteHeader(code)
			log.Println("error starting auth:", err)
			http.Error(w, err.Error(), code)
			return
		}
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		var config tmpauth.UnserializableConfig
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			log.Println("config error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		configData, err := json.Marshal(config)
		if err != nil {
			log.Println("config re-marshal error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// SHA1 the configData
		sum := sha1.Sum(configData)
		configID := hex.EncodeToString(sum[:])

		if ta, ok := tmpauthInstances[configID]; ok {
			log.Println("config already registered:", configID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(tmpauth.RemoteConfig{
				ConfigID: configID,
				Secret:   ta.Config.Secret,
				ClientID: ta.Config.ClientID,
			})

			return
		}

		parsedConfig, err := config.Parse()
		if err != nil {
			log.Println("config error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ta := tmpauth.NewTmpauth(parsedConfig, nil)

		tmpauthInstances[configID] = ta

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(tmpauth.RemoteConfig{
			ConfigID: configID,
			Secret:   ta.Config.Secret,
			ClientID: ta.Config.ClientID,
		})

		log.Println("successfully registered config:", configID)
	})

	http.HandleFunc("/.well-known/tmpauth/", func(w http.ResponseWriter, r *http.Request) {
		configID := r.Header.Get(tmpauth.ConfigIDHeader)
		if configID == "" {
			log.Println("missing config ID")
			http.Error(w, "missing config ID", http.StatusBadRequest)
			return
		}

		ta, ok := tmpauthInstances[configID]
		if !ok {
			log.Println("invalid config ID:", configID)
			http.Error(w, "invalid config ID", http.StatusPreconditionFailed)
			return
		}

		ta.Stdlib().ServeHTTP(w, r)
	})

	http.HandleFunc("/tmpauth/cache", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			CacheMinIat int64 `json:"cacheMinIat"`
		}{tmpauth.MinValidationTime().UnixMilli()})
	})

	http.HandleFunc("/tmpauth/whomst", func(w http.ResponseWriter, r *http.Request) {
		configID := r.Header.Get(tmpauth.ConfigIDHeader)
		if configID == "" {
			log.Println("missing config ID")
			http.Error(w, "missing config ID", http.StatusBadRequest)
			return
		}

		ta, ok := tmpauthInstances[configID]
		if !ok {
			log.Println("invalid config ID:", configID)
			http.Error(w, "invalid config ID", http.StatusPreconditionFailed)
			return
		}

		whomstData, err := ta.Whomst()
		if err != nil {
			log.Println("error getting whomst:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(whomstData)
	})

	log.Println("starting server on :4600")
	http.ListenAndServe(":4600", nil)
}
