package elastifile

import (
	"testing"
	"fmt"
	"github.com/elastifile/emanage-go/pkg/emanage"
	"github.com/elastifile/emanage-go/pkg/optional"
)

func TestConfigFromParamsBad(t *testing.T) {
	params := map[string]string{
		"foo":  "bar",
		"quux": "xyzzy",
	}

	_, err := newConfigFromParams(params)
	if err == nil {
		t.Fatal("expected error due to missing params")
	}
}

func TestConfigFromParamsGood(t *testing.T) {
	params := map[string]string{
		"nfsServer": "myserver",
		"restURL":   "myrest",
		"username":  "myuser",
		"password":  "mypass",
	}

	conf, err := newConfigFromParams(params)
	if err != nil {
		t.Fatal(err)
	}

	if conf.NFSServer != "myserver" {
		t.Fatal("wrong nfsServer")
	}
	if conf.EmanageURL != "myrest" {
		t.Fatal("wrong restURL")
	}
	if conf.Username != "myuser" {
		t.Fatal()
	}
	if conf.Password != "mypass" {
		t.Fatal()
	}
}

func TestConfigSetAnnotations(t *testing.T) {
	ann := map[string]string{
		"foo": "bar",
	}
	conf := config{
		NFSServer:  "myserver",
		EmanageURL: "myemanage",
	}
	conf.setAnnotations(ann)

	if ann["elastifile.com/restURL"] != "myemanage" {
		t.Fatal()
	}
}

func TestConfigFromAnnotations(t *testing.T) {
	ann := map[string]string{
		"elastifile.com/restURL":   "myemanage",
		"elastifile.com/nfsServer": "myserver",
		"elastifile.com/username":  "myuser",
		"elastifile.com/password":  "mypass",
	}

	conf, err := newConfigFromAnnotations(ann)
	if err != nil {
		t.Fatal(err)
	}

	expected := config{
		NFSServer:  "myserver",
		EmanageURL: "myemanage",
		Username:   "myuser",
		Password:   "mypass",
	}
	if *conf != expected {
		t.Fatalf("got:\n%+v\nexpected:\n%+v\n", *conf, expected)
	}
}
func TestEmanageClient(t *testing.T) {
	conf := config{
		EmanageURL: "https://10.11.209.216",
		Username:   "admin",
		Password:   "changeme",
	}
	_, err := getEmanage(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProvision(t *testing.T) {

	t.SkipNow()

	conf := config{
		EmanageURL: "https://10.11.209.211",
		Username:   "admin",
		Password:   "changeme",
	}
	mgmt, err := getEmanage(conf)
	if err != nil {
		t.Fatal(err)
	}


	// get and print  dcs
	dcs,err := mgmt.DataContainers.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _,dc := range dcs {
		fmt.Println(dc.Id)
	}

	// get and print exports
	exports, err := mgmt.Exports.GetAll(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _,export := range exports {
		fmt.Println(export.Id)
	}

	// create new export
	exportOpts := &emanage.ExportCreateOpts{
		DcId: dcs[0].Id,
		Path: "/dsdsd",
		Access: emanage.ExportAccessRW,
		UserMapping: emanage.UserMappingAll,
		Uid: optional.NewInt(0),
		}

	mgmt.Exports.Create("testprov",exportOpts)

}

