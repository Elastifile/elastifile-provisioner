package emanage

import (
	"fmt"
	"path/filepath"

	"github.com/elastifile/emanage-go/pkg/rest"
)

const (
	managersUri = "/api/vm_managers"
)

type vmManagers struct {
	conn *rest.Session
}

type VMManager struct {
	Id          int    `json:"id"`
	Server      string `json:"server"`
	Login       string `json:"login"`
	Fingerprint string `json:"fingerprint"`
	Secure      bool   `json:"secure"`
}

type VMManagerLoginOpts struct {
	Server   string `json:"server"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Secure   bool   `json:"secure"`
}

func (vm *vmManagers) GetAll(opts *GetAllOpts) ([]VMManager, error) {
	var result []VMManager
	err := vm.conn.Request(rest.MethodGet, managersUri, &opts, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (vm *vmManagers) TestConnection(vmManagerID int) error {
	var result VMManager

	managersTestUri := filepath.Join( // /api/vm_managers/1/test_connection
		managersUri,
		fmt.Sprintf("%v", vmManagerID),
		"test_connection")
	err := vm.conn.Request(rest.MethodPost, managersTestUri, nil, &result)
	if err != nil {
		return err
	}
	return nil
}
