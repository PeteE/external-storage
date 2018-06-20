package provisioner

import (
	"fmt"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/kubernetes-incubator/external-storage/lib/util"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type freenasProvisioner struct {
    Config *FreeNasConfig
}

// NewFreenasProvisioner creates new iscsi provisioner
func NewFreenasProvisioner(config *FreeNasConfig) controller.Provisioner {
	initLog()

	return &freenasProvisioner{
		Config : config,
	}
}

// getAccessModes returns access modes iscsi volume supported.
func (p *freenasProvisioner) getAccessModes() []v1.PersistentVolumeAccessMode {
	return []v1.PersistentVolumeAccessMode{
		v1.ReadWriteOnce,
		v1.ReadOnlyMany,
	}
}

// Provision creates a storage asset and returns a PV object representing it.
func (p *freenasProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	if !util.AccessModesContainedInAll(p.getAccessModes(), options.PVC.Spec.AccessModes) {
		return nil, fmt.Errorf("invalid AccessModes %v: only AccessModes %v are supported", options.PVC.Spec.AccessModes, p.getAccessModes())
	}
    volName := p.getVolumeName(options)
	log.Debugln("new provision request received for pvc: ", volName)

    size := getSize(options)
    log.Debugf("got size: %d\n", size)
	vol, err := CreateVolume(p.Config, volName, size)
	if err != nil {
		log.Warnln(err)
		return nil, err
	}
	log.Debugln("volume created with vol: ", vol )

	annotations := make(map[string]string)
	annotations["volume_name"] = vol.Name
	annotations["pool"] = p.Config.Pool

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
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
			// set volumeMode from PVC Spec
			VolumeMode: options.PVC.Spec.VolumeMode,
			PersistentVolumeSource: v1.PersistentVolumeSource{
				ISCSI: &v1.ISCSIPersistentVolumeSource{
					TargetPortal:      p.Config.Portal,
                    IQN:               fmt.Sprintf("%s:%s", p.Config.IQN, vol.Name),
					ISCSIInterface:    "default",
					Lun:               0,
					ReadOnly:          false,
					FSType:            "ext4",
                    /*
					DiscoveryCHAPAuth: false,
					SessionCHAPAuth:   false
					SecretRef:         nil,
                    */
				},
			},
		},
	}
	return pv, nil
}

// Delete removes the storage asset that was created by Provision represented
// by the given PV.
func (p *freenasProvisioner) Delete(volume *v1.PersistentVolume) error {
	//vol from the annotation
    name := volume.GetName()
	log.Debugln("volume deletion request received: ", volume.GetName())
    err := DeleteVolume(p.Config, name)
    if err != nil {
        log.Fatal(err)
        return err
    }
	return nil
}

func getSize(options controller.VolumeOptions) int64 {
	q := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	return q.Value()
}

func (p *freenasProvisioner) getVolumeName(options controller.VolumeOptions) string {
	return options.PVName
}
