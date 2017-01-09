package elastifile

import (
	"fmt"
	"net/url"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/nfs-provisioner/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/elastifile/emanage-go/pkg/emanage"
)

const (
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "nfs-dynamic-provisioner"
	annPrefix    = "elastifile.com/" // Prefix for Elastifile volume annotations
)

func NewProvisioner(client kubernetes.Interface) controller.Provisioner {
	return &provisioner{
		clientSet: client,
	}
}

type provisioner struct {
	clientSet kubernetes.Interface
}

// Fill in the password using the matching secret.
func (p *provisioner) fillSecret(conf *config) error {
	secrets := p.clientSet.CoreV1().Secrets(conf.SecretNamespace)
	if secrets == nil {
		return fmt.Errorf("no secrets found in namespace %q", conf.SecretNamespace)
	}

	secret, err := secrets.Get(conf.SecretName)
	if err != nil {
		return fmt.Errorf("secret %q: %v", conf.SecretName, err)
	}

	const passwordKey = "password.txt"
	password, ok := secret.Data[passwordKey]
	if !ok {
		return fmt.Errorf("secret %q: no value found for key %q", conf.SecretName, passwordKey)
	}

	conf.Password = string(password)
	return nil
}

// Provision creates a volume i.e. the storage asset and returns a PV object
// for the volume
func (p *provisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	conf, err := newConfigFromParams(options.Parameters)
	if err != nil {
		return nil, err
	}

	if err := p.fillSecret(conf); err != nil {
		return nil, err
	}

	volSource, err := p.createVolume(options, *conf)
	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy
	conf.setAnnotations(annotations)

	pv := &v1.PersistentVolume{
		ObjectMeta: v1.ObjectMeta{
			Name:        options.PVName,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: volSource,
			},
		},
	}

	return pv, nil
}

// Connect to Elastifile system
func getEmanage(conf config) (*emanage.Client, error) {
	baseURL, err := url.Parse(conf.EmanageURL)
	if err != nil {
		return nil, err
	}

	mgmt := emanage.NewClient(baseURL)
	err = mgmt.Sessions.Login(conf.Username, conf.Password)
	return mgmt, err
}

func (p *provisioner) createVolume(options controller.VolumeOptions, conf config) (*v1.NFSVolumeSource, error) {
	name := options.PVName

	glog.Infof("Creating Elastifile Data Container '%v'", name)

	mgmt, err := getEmanage(conf)
	if err != nil {
		return nil, err
	}

	capacity := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]

	// Create DataContainer for volume
	dc, err := mgmt.DataContainers.Create(name, 1, &emanage.DcCreateOpts{
		SoftQuota: int(capacity.Value()),
		HardQuota: int(capacity.Value()),
	})
	if err != nil {
		return nil, err
	}

	export := "root" // TODO: Make this a parameter?
	_, err = mgmt.Exports.Create(export, &emanage.ExportCreateOpts{
		DcId:        dc.Id,
		Path:        "/",
		UserMapping: emanage.UserMappingNone,
	})
	if err != nil {
		return nil, err
	}

	// Return values
	// path.Join()

	return &v1.NFSVolumeSource{
		Server:   conf.NFSServer,
		Path:     fmt.Sprintf("/%v/%v", name, export), // TODO: Use path.Join()
		ReadOnly: false,
	}, nil
}

// Delete removes the storage asset that was created by Provision backing the
// given PV. Does not delete the PV object itself.
//
// May return IgnoredError to indicate that the call has been ignored and no
// action taken. In case multiple provisioners are serving the same storage
// class, provisioners may ignore PVs they are not responsible for (e.g. ones
// they didn't create). The controller will act accordingly, i.e. it won't
// emit a misleading VolumeFailedDelete event.
func (p *provisioner) Delete(pv *v1.PersistentVolume) error {
	name := pv.Name

	conf, err := newConfigFromAnnotations(pv.Annotations)
	if err != nil {
		return err
	}

	if err := p.fillSecret(conf); err != nil {
		return err
	}

	glog.Infof("Going to delete Elastifile Data Container '%v'", name)
	mgmt, err := getEmanage(*conf)
	if err != nil {
		return err
	}

	dcs, err := mgmt.DataContainers.GetAll(nil)
	if err != nil {
		return err
	}

	found := false
	for _, dc := range dcs {
		if dc.Name == name {
			if err := deleteAllExports(mgmt, dc); err != nil {
				return err
			}

			if _, err := mgmt.DataContainers.Delete(&dc); err != nil {
				return err
			}

			glog.Infof("Deleted Elastifile Data Container '%v'", name)

			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no such Data Container: '%v'", name)
	}

	return nil
}

func deleteAllExports(mgmt *emanage.Client, dc emanage.DataContainer) error {
	exps, err := mgmt.Exports.GetAll(nil)
	if err != nil {
		return err
	}

	for _, e := range exps {
		if e.DataContainerId == dc.Id {
			if _, err := mgmt.Exports.Delete(&e); err != nil {
				return err
			}
			glog.Infof("Deleted Elastifile Export '%v'", e.Name)
		}
	}

	return nil
}
