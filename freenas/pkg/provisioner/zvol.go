package provisioner
import (
    "fmt"
	"net/http"
    "encoding/json"
    "bytes"
)

// http://api.freenas.org/resources/storage.html#volume
type ZPool struct {
	ID       int64 `json:"id"`
	Avail    int64 `json:"avail"`
	Children []struct {
		ID       int64 `json:"id"`
		Avail    int64 `json:"avail"`
		Children []struct {
			ID         int64  `json:"id"`
			Avail      int64  `json:"avail"`
			Mountpoint string `json:"mountpoint"`
			Name       string `json:"name"`
			Path       string `json:"path"`
			Status     string `json:"status"`
			Type       string `json:"type"`
			Used       int64  `json:"used"`
			UsedPct    int64  `json:"used_pct"`
		} `json:"children"`
		Mountpoint string `json:"mountpoint"`
		Name       string `json:"name"`
		Path       string `json:"path"`
		Status     string `json:"status"`
		Type       string `json:"type"`
		Used       int64  `json:"used"`
		UsedPct    int64  `json:"used_pct"`
	} `json:"children"`
	IsDecrypted   bool    `json:"is_decrypted"`
	IsUpgraded    bool    `json:"is_upgraded"`
	Mountpoint    string  `json:"mountpoint"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	Used          int64   `json:"used"`
	UsedPct       string  `json:"used_pct"`
	VolEncrypt    int64   `json:"vol_encrypt"`
	VolEncryptkey string  `json:"vol_encryptkey"`
	VolFstype     string  `json:"vol_fstype"`
	VolGuid       float64 `json:"vol_guid,string"`
	VolName       string  `json:"vol_name"`
}

func CreateVolume(config *FreeNasConfig, vol string, size int64) (*ZVol, error){
    zvolUtil := &ZVolUtil{
        Config: config,
    }
    var err error
    v, err := zvolUtil.Create(vol, size, true, "4k")
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    log.Printf("Created vol: %s\n", v.Name)

    extentUtil := &ExtentUtil{
        Config: config,
    }
    e, err := extentUtil.Create(v.Name)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    targetUtil := &TargetUtil{
        Config: config,
    }
    t, err := targetUtil.Create(v.Name)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    log.Printf("Created target: %s\n", t.Name)

    targetGroupUtil := &TargetGroupUtil{
        Config: config,
    }
    tg, err := targetGroupUtil.Create(t.Id)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    fmt.Printf("Created target group: %d\n", tg.Id)

    teUtil := &TargetExtentUtil{
        Config: config,
    }
    te, err := teUtil.Create(t.Id, e.Id, 0)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    fmt.Printf("Created targetextent: %d\n", te.Id)
    return v, nil
}
func DeleteVolume(config *FreeNasConfig, name string) (error) {
    var err error
    // find the extent
    extentUtil := &ExtentUtil{
        Config: config,
    }
    extent, err := extentUtil.Find(name)
    if err != nil {
        log.Fatal(err)
        return err
    }
    // find the target
    targetUtil := &TargetUtil{
        Config: config,
    }
    target, err := targetUtil.Find(name)
    if err != nil {
        log.Fatal(err)
        return err
    }

    // delete the targettoextent
    teUtil := &TargetExtentUtil{
        Config: config,
    }
    err = teUtil.Delete(target.Id, extent.Id)
    if err != nil {
        log.Fatal(err)
        return err
    }
    fmt.Printf("Deleted targetextent: Target = %d, Extenet = %d\n", target.Id, extent.Id)

    // delete the targetgroup
    targetGroupUtil := &TargetGroupUtil{
        Config: config,
    }
    err = targetGroupUtil.Delete(target.Id)
    if err != nil {
        log.Fatal(err)
        return err
    }
    fmt.Printf("Deleted target group: Target = %d\n", target.Id)

    // delete the target
    targetUtil.Delete(target.Id)

    // delete the extent
    err = extentUtil.Delete(extent.Name)
    if err != nil {
        log.Fatal(err)
        return err
    }
    fmt.Printf("Deleted extent: Name = %d\n", extent.Name)

    // delete the zvol
    zvolUtil := &ZVolUtil{
        Config: config,
    }
    err = zvolUtil.Delete(name)
    if err != nil {
        log.Fatal(err)
        return err
    }
    log.Debugf("Deleted zvol: Name = %s\n", name)
    return nil
}
type ZPoolUtil struct {
    Config *FreeNasConfig
}

func (u *ZPoolUtil) Find(name string) (*ZPool, error) {
    zpools, err  := u.List()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    for _, zpool := range zpools {
        if zpool.Name == u.Config.Pool {
            for _, ds := range zpool.Children {
                if ds.Type == "dataset" && ds.Name == u.Config.Pool {
                    for _, zvol := range ds.Children {
                        if zvol.Type == "zvol" && zvol.Name == name {
                            z := &ZPool {
                                ID: zvol.ID,
                                Name: zvol.Name,
                                Avail: zvol.Avail,
                            }

                            return z, nil
                        }
                    }
                }

            }
        }
    }
    return nil, fmt.Errorf("Unable to find zvol with Name = %s", name)
}

func (u *ZPoolUtil) List() ([]ZPool, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/storage/volume/?format=json", u.Config.Uri)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var pools []ZPool
    json.NewDecoder(resp.Body).Decode(&pools)
    return pools, nil
}

type ZVolUtil struct {
    Config *FreeNasConfig
}

type ZVol struct {
    // common params
    Comments string `json:"comments,omitempty"`
    Name string `json:"name"`
    Volsize int64 `json:"volsize"`
    Compression string `json:"compression,omitempty"`

    // request only
    Sparse bool `json:"sparse,omitempty"`
    Force bool `json:"force,omitempty"`
    BlockSize string `json:"blocksize,omitempty"`

    // response only
    Avail int64 `json:"avail,omitempty"`
    Dedup string `json:"dedupe,omitempty"`
    Refer int64 `json:"refer,omitempty"`
    Used int64 `jons:"used,omitempty"`
}

func (u *ZVolUtil) List() ([]ZVol, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/storage/volume/%s/zvols/?format=json", u.Config.Uri, u.Config.Pool)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var vols []ZVol
    json.NewDecoder(resp.Body).Decode(&vols)
    return vols, nil
}
func (u *ZVolUtil) Create(name string, size int64, sparse bool, blocksize string) (*ZVol, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/storage/volume/%s/zvols/", u.Config.Uri, u.Config.Pool)
    vol := &ZVol{
        Name: name,
        Volsize: size,
        Sparse: sparse,
        BlockSize: blocksize,
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(vol)
    log.Printf("Posting: %s", b.String())
    req, err := http.NewRequest(http.MethodPost, uri, b)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    req.Header.Set("Content-Type", "application/json")
    resp, err := httpClient.Do(req)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    if resp.StatusCode != 201 {
        log.Fatal(resp.Status)
        return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
    }
    json.NewDecoder(resp.Body).Decode(vol)
    return vol, nil
}
func (u *ZVolUtil) Delete(name string) (error) {
    var uri = fmt.Sprintf("%s/api/v1.0/storage/volume/%s/zvols/%s/", u.Config.Uri, u.Config.Pool, name)
    req, err := http.NewRequest(http.MethodDelete, uri, nil)
    if err != nil {
        log.Fatal(err)
        return err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    req.Header.Set("Content-Type", "application/json")
    resp, err := httpClient.Do(req)
    if err != nil {
        log.Fatal(err)
        return  err
    }
    if resp.StatusCode != 200 && resp.StatusCode != 204 {
        log.Fatal(resp.Status)
        return fmt.Errorf("Request failed with status: %s", resp.Status)
    }
    return nil
}
