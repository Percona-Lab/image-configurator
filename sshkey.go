package main

import (
	"encoding/json"
	"net/http"

	selinux "github.com/opencontainers/selinux/go-selinux"
)

func getSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	parsedSSHKey, result, err := SSHKey.Read()
	if result == "success" {
		json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
	} else {
		returnError(w, req, http.StatusInternalServerError, result, err)
	}
}

func setSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	parsedSSHKey, result, err := SSHKey.Write(req.Body)
	if result != "success" {
		returnError(w, req, http.StatusInternalServerError, result, err)
	}

	if err = selinux.SetFileLabel(SSHKey.KeyPath, "system_u:object_r:httpd_sys_content_t:s0"); err != nil { // nolint
		returnError(w, req, http.StatusInternalServerError, "Cannot set SELinux context", err)
	}

	w.Header().Set("Location", req.URL.String())
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
}
