// The emanage package provides access to the Emanage REST API.
//
// Object names and data types are closely related to the names exported by Emanage.
//
// NOTE: This package is far from final and needs to undergo a lot of changes,
// as outlined in comments throughout the package's code.
package emanage

import (
	"net"
	"net/url"
	"time"

	"github.com/elastifile/errors"
	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/elastifile/emanage-go/pkg/rest"
	"github.com/elastifile/emanage-go/pkg/retry"
)

var Log = log.New("package", "emanage")

func init() {
	Log.SetHandler(log.DiscardHandler())
}

type Client struct {
	ClusterReports    *clusterReports
	DataContainers    *dataContainers
	Enodes            *enodes
	ControlTasks      *controlTasks
	Exports           *exports
	Hosts             *hosts
	Events            *events
	NetworkInterfaces *netInterfaces
	Policies          *policies
	Sessions          *rest.Session
	Statistics        *statistics
	Systems           *systems
	Tenants           *tenants
	VMManagers        *vmManagers
	VMs               *vms

	log.Logger
}

func EmanageURL(host string) *url.URL {
	baseURL := &url.URL{
		Scheme: "http",
		Host:   host,
	}
	return baseURL
}

func NewClient(baseURL *url.URL) *Client {
	s := rest.NewSession(baseURL)
	return &Client{
		ClusterReports:    &clusterReports{s},
		DataContainers:    &dataContainers{s},
		Enodes:            &enodes{s},
		Events:            &events{s},
		ControlTasks:      &controlTasks{s},
		Exports:           &exports{s},
		Hosts:             &hosts{s},
		NetworkInterfaces: &netInterfaces{s},
		Policies:          &policies{s},
		Sessions:          s,
		Statistics:        &statistics{s},
		Systems:           &systems{s},
		Tenants:           &tenants{s},
		VMManagers:        &vmManagers{s},
		VMs:               &vms{s},

		Logger: Log.New("baseURL", baseURL.String()),
	}
}

func (client *Client) RetriedLogin(username string, password string, timeout time.Duration) error {
	interval := 1 * time.Second
	err := retry.Basic{
		Timeout: interval,
		Retries: int(timeout / interval),
	}.Do(func() error {
		Log.Info("Trying to login to emanage...")
		e := client.Sessions.Login(username, password)
		if ue, ok := errors.Inner(e).(*url.Error); ok {
			if _, ok := ue.Err.(*net.OpError); ok {
				Log.Info("Connection error, retrying", "err", e)
				return &retry.TemporaryError{Err: e}
			}
		}
		return e
	})
	return err
}
