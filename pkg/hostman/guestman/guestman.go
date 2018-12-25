package guestman

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/util/regutils"
	"yunion.io/x/pkg/util/seclib"

	"yunion.io/x/onecloud/pkg/cloudcommon/httpclients"
	"yunion.io/x/onecloud/pkg/cloudcommon/sshkeys"
	"yunion.io/x/onecloud/pkg/cloudcommon/workmanager"
	"yunion.io/x/onecloud/pkg/hostman/guestfs"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/util/timeutils2"
)

const VNC_PORT_BASE = 5900

type SGuestManager struct {
	ServersPath      string
	Servers          map[string]*SKVMGuestInstance
	CandidateServers map[string]*SKVMGuestInstance
	ServersLock      *sync.Mutex

	isLoaded bool
}

func NewGuestManager(serversPath string) *SGuestManager {
	manager := &SGuestManager{}
	manager.ServersPath = serversPath
	manager.Servers = make(map[string]*SKVMGuestInstance, 0)
	manager.CandidateServers = make(map[string]*SKVMGuestInstance, 0)
	manager.ServersLock = &sync.Mutex{}
	manager.StartCpusetBalancer()
	manager.LoadExistingGuests()
	return manager
}

func (m *SGuestManager) Bootstrap() {
	if m.isLoaded || len(m.ServersPath) == 0 {
		log.Errorln("Guestman bootstrap has been called!!!!!")
	} else {
		m.isLoaded = true
		log.Infof("Loading existing guests ...")
		if len(m.CandidateServers) > 0 {
			m.VerifyExistingGuests(false)
		} else {
			m.OnLoadExistingGuestsComplete()
		}
	}
}

func (m *SGuestManager) VerifyExistingGuests(pendingDelete bool) {
	params := url.Values{
		"limit":          {"0"},
		"admin":          {"True"},
		"system":         {"True"},
		"pending_delete": {fmt.Sprintf("%s", pendingDelete)},
	}
	params.Set("filter.0", fmt.Sprintf("host_id.equals(%s)", "get host id //TODO"))
	if len(m.CandidateServers) > 0 {
		keys := make([]string, len(m.CandidateServers))
		var index = 0
		for k := range m.CandidateServers {
			keys[index] = k
			index++
		}
		params.Set("filter.1", strings.Join(keys, ","))
	}
	urlStr := fmt.Sprintf("/servers?%s", params.Encode())
	// TODO: get default context not use background context
	_, res, err := httpclients.GetDefaultComputeClient().Request(context.Background(), "GET", urlStr, nil, nil, false)
	if err != nil {
		m.OnVerifyExistingGuestsFail(err, pendingDelete)
	} else {
		m.OnVerifyExistingGuestsSucc(res, pendingDelete)
	}
}

func (m *SGuestManager) OnVerifyExistingGuestsFail(err error, pendingDelete bool) {
	log.Errorf("OnVerifyExistingGuestFail: %s, try again 30 seconds later", err.Error())
	timeutils2.AddTimeout(30*time.Second, func() { m.VerifyExistingGuests(false) })
}

func (m *SGuestManager) OnVerifyExistingGuestsSucc(res jsonutils.JSONObject, pendingDelete bool) {
	iServers, err := res.Get("servers")
	if err != nil {
		m.OnVerifyExistingGuestsFail(err, pendingDelete)
	} else {
		servers := iServers.(*jsonutils.JSONArray)
		for _, v := range servers.Value() {
			id, _ := v.GetString("id")
			server, ok := m.CandidateServers[id]
			if !ok {
				log.Errorf("verify_existing_guests return unknown server %s ???????", id)
			} else {
				server.ImportServer(pendingDelete)
			}
		}
		if !pendingDelete {
			m.VerifyExistingGuests(true)
		} else {
			var unknownServerrs = make([]*SKVMGuestInstance, 0)
			for _, server := range m.CandidateServers {
				log.Errorf("Server %s not found on this host", server.GetName())
				unknownServerrs = append(unknownServerrs, server)
			}
			for _, server := range unknownServerrs {
				m.RemoveCandidateServer(server)
			}
		}
	}
}

func (m *SGuestManager) RemoveCandidateServer(server *SKVMGuestInstance) {
	if _, ok := m.CandidateServers[server.GetId()]; ok {
		delete(m.CandidateServers, server.GetId())
		if len(m.CandidateServers) == 0 {
			m.OnLoadExistingGuestsComplete()
		}
	}
}

func (m *SGuestManager) OnLoadExistingGuestsComplete() {
	// TODO
}

func (m *SGuestManager) StartCpusetBalancer() {
	// TODO
}

func (m *SGuestManager) IsGuestDir(f os.FileInfo) bool {
	if !regutils.MatchUUID(f.Name()) {
		return false
	}
	if !f.Mode().IsDir() {
		return false
	}
	descFile := path.Join(m.ServersPath, f.Name(), "desc")
	if _, err := os.Stat(descFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (m *SGuestManager) LoadExistingGuests() {
	files, err := ioutil.ReadDir(m.ServersPath)
	if err != nil {
		log.Errorf("List servers path %s error %s", m.ServersPath, err)
	}
	for _, f := range files {
		if _, ok := m.Servers[f.Name()]; !ok && m.IsGuestDir(f) {
			log.Infof("Find existing guest %s", f.Name())
			m.LoadServer(f.Name())
		}
	}
}

func (m *SGuestManager) LoadServer(sid string) {
	guest := NewKVMGuestInstance(sid, m)
	err := guest.LoadDesc()
	if err != nil {
		log.Errorf("On load server error: %s", err)
		return
	}
	m.CandidateServers[sid] = guest
}

func (m *SGuestManager) PrepareCreate(sid string) error {
	m.ServersLock.Lock()
	defer m.ServersLock.Unlock()
	if _, ok := m.Servers[sid]; ok {
		return httperrors.NewBadRequestError("Guest %s exists", sid)
	}
	guest := NewKVMGuestInstance(sid, m)
	m.Servers[sid] = guest
	return guest.PrepareDir()
}

func (m *SGuestManager) PrepareDeploy(sid string) error {
	m.ServersLock.Lock()
	defer m.ServersLock.Unlock()
	if guest, ok := m.Servers[sid]; !ok {
		return httperrors.NewBadRequestError("Guest %s not exists", sid)
	} else {
		if guest.IsRunning() || guest.IsSuspend() {
			return httperrors.NewBadRequestError("Cannot deploy on running/suspend guest")
		}
	}
	return nil
}

func (m *SGuestManager) Monitor(sid, cmd string, callback func(string)) error {
	if guest, ok := m.Servers[sid]; ok {
		if guest.IsRunning() {
			guest.monitor.SimpleCommand(cmd, callback)
			return nil
		} else {
			return httperrors.NewBadRequestError("Server stopped??")
		}
	} else {
		return httperrors.NewNotFoundError("Not found")
	}
}

// Delay process
func (m *SGuestManager) DoDeploy(ctx context.Context, params interface{}) (jsonutils.JSONObject, error) {
	deployParams, ok := params.(*SGuestDeploy)
	if !ok {
		return nil, fmt.Errorf("Unknown params")
	}
	guest, ok := m.Servers[deployParams.sid]
	if ok {
		desc, _ := deployParams.body.Get("desc")
		if desc != nil {
			guest.SaveDesc(desc)
		}
		if jsonutils.QueryBoolean(deployParams.body, "k8s_pod", false) {
			return nil, nil
		}
		publicKey := sshkeys.GetKeys(deployParams.body)
		deploys, _ := deployParams.body.GetArray("deploys")
		password, _ := deployParams.body.GetString("password")
		resetPassword := jsonutils.QueryBoolean(deployParams.body, "reset_password", false)
		if resetPassword && len(password) == 0 {
			password = seclib.RandomPassword(12)
		}

		guestInfo, err := guest.DeployFs(&guestfs.SDeployInfo{
			publicKey, deploys, password, deployParams.isInit})
		if err != nil {
			log.Errorf("Deploy guest fs error: %s", err)
			return nil, err
		} else {
			return guestInfo, nil
		}
	} else {
		return nil, fmt.Errorf("Guest %s not found", sid)
	}
}

// delay cpuset balance
func (m *SGuestManager) CpusetBalance(ctx context.Context, params interface{}) (jsonutils.JSONObject, error) {
	// TODO
}

func (m *SGuestManager) Status(sid string) string {
	if guest, ok := m.Servers[sid]; ok {
		// TODO
		// if guest.IsMaster() && !guest.IsMirrorJobSucc() {
		// 	return "block_stream"
		// }
		if guest.IsRunning() {
			return "running"
		} else if guest.IsSuspend() {
			return "suspend"
		} else {
			return "stopped"
		}
	} else {
		return "notfound"
	}
}

func (m *SGuestManager) Delete(sid string) (*SKVMGuestInstance, error) {
	if guest, ok := m.Servers[sid]; ok {
		delete(m.Servers, sid)
		// 这里应该不需要append到deleted servers, 据观察 deleted servers没有用到
		return guest, nil
	} else {
		return nil, httperrors.NewNotFoundError("Not found")
	}
}

func (m *SGuestManager) GuestStart(ctx context.Context, sid string, body jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	if guest, ok := m.Servers[sid]; ok {
		if desc, err := body.Get("desc"); err != nil {
			guest.SaveDesc(desc)
		}
		if guest.IsStopped() {
			params, _ := body.Get("params")
			if err := guest.StartGuest(ctx, params); err != nil {
				return nil, httperrors.NewBadRequestError("Failed to start server")
			} else {
				return jsonutils.NewDict(jsonutils.JSONPair{"vnc_port", jsonutils.NewInt(0)}), nil
			}
		} else {
			vncPort := guest.GetVncPort()
			if vncPort > 0 {
				res := jsonutils.NewDict()
				res.Set("vnc_port", jsonutils.NewInt(int64(vncPort)))
				res.Set("is_running", jsonutils.JSONTrue)
				return res, nil
			} else {
				return nil, httperrors.NewBadRequestError("Seems started, but no VNC info")
			}
		}
	} else {
		return nil, httperrors.NewNotFoundError("Not found")
	}
}

func (m *SGuestManager) GuestStop(ctx context.Context, sid string, timeout int64) error {
	if guest, ok := m.Servers[sid]; !ok {
		guest.ExecStopTask(ctx, timeout)
		return nil
	} else {
		return httperrors.NewNotFoundError("Guest %s not found", sid)
	}
}

func (m *SGuestManager) GetFreeVncPort() int64 {
	vncPorts := make(map[int]struct{}, 0)
	for _, guest := range m.Servers {
		inUsePort := guest.GetVncPort()
		if inUsePort > 0 {
			vncPorts[inUsePort] = struct{}{}
		}
	}
	var port = 1
	for {
		// TODO: IsTcpPortUsed
		if _, ok := vncPorts[port]; !ok && !netutils2.IsTcpPortUsed("0.0.0.0", VNC_PORT_BASE+port) &&
			!netutils2.IsTcpPortUsed("0.0.0.0", MONITOR_PORT_BASE+port) {
			break
		} else {
			port += 1
		}
	}
	return port
}

func Stop() {
	// guestManger.ExitGuestCleanup()
}

func Init(serversPath string) {
	initGuestManager(serversPath)
}

var guestManger *SGuestManager
var wm *workmanager.SWorkManager

func GetGuestManager() *SGuestManager {
	return guestManger
}

func initGuestManager(serversPath string) {
	if guestManger == nil {
		guestManger = NewGuestManager(serversPath)
	}
}

func GetWorkManager() *workmanager.SWorkManager {
	return wm
}

func init() {
	wm = workmanager.NewWorkManger()
}
