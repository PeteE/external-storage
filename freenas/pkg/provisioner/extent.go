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
    var url = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/?format=json", u.Config.Url)
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
        //log.Fatal(err)
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
    var url = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/", u.Config.Url)
    extent := &Extent{
        EType: "Disk",
        Name: name,
        Disk: fmt.Sprintf("zvol/%s/%s", u.Config.Pool, name),
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(extent)
    //log.Printf("Posting: %s", b.String())
    req, err := http.NewRequest(http.MethodPost, url, b)
    if err != nil {
        //log.Fatal(err)
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    req.Header.Set("Content-Type", "application/json")
    resp, err := httpClient.Do(req)
    if err != nil {
        //log.Fatal(err)
        return nil, err
    }
    if resp.StatusCode != 201 {
        //log.Fatal(resp.Status)
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
        //log.Fatal(err)
        return err
    }
    for _, e := range extents {
        if e.Name == name {
            //log.Printf("Found match: %s, %d", e.Name, e.Id)
            var url = fmt.Sprintf("%s/api/v1.0/services/iscsi/extent/%d/", u.Config.Url, e.Id)
            req, err := http.NewRequest(http.MethodDelete, url, nil)
            if err != nil {
                log.Fatal(err)
                return err
            }
            req.SetBasicAuth(u.Config.Username, u.Config.Password)
            req.Header.Set("Content-Type", "application/json")
            resp, err := httpClient.Do(req)
            if err != nil {
                //log.Fatal(err)
                return err
            }
            if resp.StatusCode != 200 && resp.StatusCode != 204 {
                //log.Fatal(resp.Status)
                return fmt.Errorf("Request failed with status: %s", resp.Status)
            }
            return nil
        }
    }
    return fmt.Errorf("No volume found matching: %s", name)
}
