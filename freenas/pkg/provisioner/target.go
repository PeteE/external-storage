package provisioner
import (
    "fmt"
	"net/http"
    "encoding/json"
    "bytes"
)

type TargetUtil struct {
    Config *FreeNasConfig
}
type Target struct {
    // common params
    Name string `json:"iscsi_target_name"`
    Id uint `json:"id,omitempty"`
    Alias string `json:"iscsi_target_alias,omitempty"`
    Mode string `json:"scsi_target_mode,omitempty"`
}

func (u *TargetUtil) List() ([]Target, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/target/?format=json", u.Config.Uri)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var targets []Target
    json.NewDecoder(resp.Body).Decode(&targets)
    return targets, nil
}
func (u *TargetUtil) Find(name string) (*Target, error) {
    targets, err := u.List()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    for _, target := range targets {
        if target.Name == name {
            return &target, nil
        }
    }
    return nil, fmt.Errorf("Unable to find target: Name = %s", name)
}

func (u *TargetUtil) Create(name string) (*Target, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/target/", u.Config.Uri)
    t := &Target{
        Name: name,
        Alias: name,
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(t)
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
        return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
    }
    json.NewDecoder(resp.Body).Decode(t)
    return t, nil
}
func (u *TargetUtil) Delete(targetId uint) (error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/target/%d/", u.Config.Uri, targetId)
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
        return fmt.Errorf("Request failed with status: %s", resp.Status)
    }
    return nil
}

type TargetGroupUtil struct {
    Config *FreeNasConfig
}
type TargetGroup struct {
    // common params
    TargetId uint `json:"iscsi_target,omitempty"`
    Id uint `json:"id,omitempty"`
    AuthGroup *string `json:"iscsi_target_authgroup,omitempty"`
    AuthType string `json:"iscsi_target_authtype,omitempty"`
    PortalGroup uint `json:"iscsi_target_portalgroup,omitempty"`
    InitiatorGroup uint `json:"iscsi_target_initiatorgroup,omitempty"`
    InitialDigest string `json:"iscsi_target_initialdigest,omitempty"`
}

func (u *TargetGroupUtil) List() ([]TargetGroup, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targetgroup/?format=json", u.Config.Uri)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var groups []TargetGroup
    json.NewDecoder(resp.Body).Decode(&groups)
    return groups, nil
}
func (u *TargetGroupUtil) Find(targetId uint) (*TargetGroup, error) {
    targets, err := u.List()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    for _, target := range targets {
        if target.TargetId == targetId {
            return &target, nil
        }
    }
    return nil, fmt.Errorf("Unable to find TargetGroup: Target = %d\n", targetId)
}
func (u *TargetGroupUtil) Create(targetId uint) (*TargetGroup, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targetgroup/", u.Config.Uri)
    g := &TargetGroup{
        TargetId: targetId,
        AuthGroup: nil,
        AuthType: "None",
        PortalGroup: 1,
        InitiatorGroup: 1,
        InitialDigest: "Auto",
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(g)
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
    json.NewDecoder(resp.Body).Decode(g)
    return g, nil
}
func (u *TargetGroupUtil) Delete(target uint) (error) {
    tgs, err := u.List()
    if err != nil {
        log.Fatal(err)
        return err
    }
    for _, tg := range tgs {
        if tg.TargetId == target {
            var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targetgroup/%d/", u.Config.Uri, tg.Id)
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
            if resp.StatusCode != 204 {
                return fmt.Errorf("Request failed with status: %s", resp.Status)
            }
            return nil
        }
    }
    return fmt.Errorf("Unable to find TargetGroup: Target = %d\n", target)
}

type TargetExtentUtil struct {
    Config *FreeNasConfig
}
type TargetExtentRequest struct {
    Target uint `json:"iscsi_target,omitempty"`
    Extent uint `json:"iscsi_extent,omitempty"`
    LunId uint `json:"iscsi_lunid"`
}
type TargetExtentResponse struct {
    Id uint `json:"id,omitempty"`
    Extent uint `json:"iscsi_extent,omitempty"`
    LunId uint `json:"scsi_lunid"`
    Target uint `json:"iscsi_target,omitempty"`
}

func (u *TargetExtentUtil) List() ([]TargetExtentResponse, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targettoextent/?format=json", u.Config.Uri)
    req, err := http.NewRequest(http.MethodGet, uri, nil)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    req.SetBasicAuth(u.Config.Username, u.Config.Password)
    resp, err := httpClient.Do(req)
    if err != nil{
        log.Fatal(err)
        return nil, err
    }
    var tes []TargetExtentResponse
    json.NewDecoder(resp.Body).Decode(&tes)
    return tes, nil
}

func (u *TargetExtentUtil) Create(targetId uint, extentId uint, lunId uint) (*TargetExtentResponse, error) {
    var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targettoextent/", u.Config.Uri)
    tereq := &TargetExtentRequest{
        Target: targetId,
        Extent: extentId,
        LunId: lunId,
    }
    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(tereq)
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
    teresp := &TargetExtentResponse{}
    json.NewDecoder(resp.Body).Decode(teresp)
    return teresp, nil
}
func (u *TargetExtentUtil) Delete(target uint, extent uint) (error) {
    // get list
    tes, err := u.List()
    if err != nil {
        log.Fatal(err)
        return err
    }
    for _, te := range tes {
        if te.Extent == extent && te.Target == target {
            // do delete
            var uri = fmt.Sprintf("%s/api/v1.0/services/iscsi/targettoextent/%d/", u.Config.Uri, te.Id)
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
            if resp.StatusCode != 204 {
                return fmt.Errorf("Request failed with status: %s", resp.Status)
            }
            return nil
        }
    }
    return fmt.Errorf("Unable to find TargetExtent: Target = %d, Extent = %d", target, extent)
}
