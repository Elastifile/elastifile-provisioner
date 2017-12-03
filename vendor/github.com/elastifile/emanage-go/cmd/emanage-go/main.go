package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/elastifile/emanage-go/pkg/emanage"
)

type options struct {
	URL      string
	Username string
	Password string
}

func connect(opts options) (*emanage.Client, error) {
	baseURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	if baseURL.Scheme == "" {
		println("Didn't find Scheme in url, assuming http")
		baseURL.Scheme = "http"
	}

	mgmt := emanage.NewClient(baseURL)

	err = mgmt.Sessions.Login(opts.Username, opts.Password)
	return mgmt, err
}

func do() error {
	var (
		opts options
		ok   bool
	)

	if opts.URL, ok = os.LookupEnv("EMANAGE_URL"); !ok {
		return fmt.Errorf("EMANAGE_URL not set")
	}

	if opts.Username, ok = os.LookupEnv("EMANAGE_USERNAME"); !ok {
		return fmt.Errorf("EMANAGE_USERNAME not set")
	}

	if opts.Password, ok = os.LookupEnv("EMANAGE_PASSWORD"); !ok {
		return fmt.Errorf("EMANAGE_PASSWORD not set")
	}

	mgmt, err := connect(opts)
	if err != nil {
		return err
	}

	dcs, err := mgmt.DataContainers.GetAll(nil)
	if err != nil {
		return err
	}

	fmt.Printf("Data Containers:\n")
	for _, dc := range dcs {
		fmt.Printf("%v\n", dc.Name)
	}

	return nil
}

func main() {
	err := do()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
