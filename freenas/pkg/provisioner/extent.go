package provisioner
import (
    "fmt"
	"net/http"
    "encoding/json"
    "bytes"
)

type ExtentUtil struct {
    Config *FreeNasConfig
}
type Extent struct {
    // common params
    EType string `json:"iscsi_target_extent_type"`
    Name  string `json:"iscsi_target_extent_name"`
    Disk  string `json:"iscsi_target_extent_disk"`

    Id uint `json:"id,omitempty"`
    BlockSize uint `json:"iscsi_target_extent_blocksize,omitempty"`
}

func (u *ExtentUtil) List() ([]Extent, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/?format=json", u.Config.Uri)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var extents[ ]Extent
    json.NewDecoder(resp.Body).Decode(&extents)
    return extents, nil
}
func (u *ExtentUtil) Find(name string) (*Extent, error) {
    extents, err := u.List()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    for _, extent := range extents {
        if extent.Name == name {
            return &extent, nil
        }
    }
    return nil, fmt.Errorf("Unable to find extent: Name = %s", name)
}

func (u *ExtentUtil) Create(name string) (*Extent, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/", u.Config.Uri)
    extent := &Extent{
        EType: "Disk",
        Name: name,
        Disk: fmt.Sprintf("zvol/%s/%s", u.Config.Pool, name),
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(extent)
    req, err := http.NewRequest(http.MethodPost, uri, b)
    if err != nil {
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    req.Header.Set("Content-Type", "application/json")
    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 201 {
        return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
    }
    json.NewDecoder(resp.Body).Decode(extent)
    return extent, nil
}

func (u *ExtentUtil) Delete(name string) (error) {
    if name == "" {
        return fmt.Errorf("name is empty")
    }
    extents, err := u.List()
    if err != nil {
        return err
    }
    for _, e := range extents {
        if e.Name == name {
            var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/%d/", u.Config.Uri, e.Id)
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
                return err
            }
            if resp.StatusCode != 200 && resp.StatusCode != 204 {
                return fmt.Errorf("Request failed with status: %s", resp.Status)
            }
            return nil
        }
    }
    return fmt.Errorf("No volume found matching: %s", name)
}
