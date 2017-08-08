package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/percona/pmm-manage/configurator/config"
	"github.com/percona/pmm-manage/configurator/sshkey"
)

// SSHKey global variable with KeyPath and KeyOwner
var SSHKey sshkey.Handler
var c config.PMMConfig

func main() {
	c = config.ParseConfig()
	SSHKey = sshkey.Init(c)
	SSHKey.RunSSHKeyChecks()

	router := mux.NewRouter().PathPrefix(c.PathPrefix).Subrouter()

	router.HandleFunc("/v1/sshkey", getSSHKeyHandler).Methods("GET")
	router.HandleFunc("/v1/sshkey", setSSHKeyHandler).Methods("POST")

	log.WithFields(log.Fields{
		"address": c.ListenAddress,
	}).Info("DB Configurator is started")
	log.Fatal(http.ListenAndServe(c.ListenAddress, router))
}

func returnError(w http.ResponseWriter, req *http.Request, httpStatus int, title string, err error) {
	response := jsonResponce{
		Code:   httpStatus,
		Status: http.StatusText(httpStatus),
		Title:  title,
	}
	if err != nil {
		response.Detail = err.Error()
	}

	responseJSON, _ := json.Marshal(response)
	responseJSONQuoted := strings.Trim(strconv.Quote(string(responseJSON)), "\"")
	log.Errorf("%s %s: %s", req.Method, req.URL.String(), responseJSONQuoted)

	http.Error(w, string(responseJSON)+"\n", httpStatus)
}
