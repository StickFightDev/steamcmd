package steamcmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

var (
	// ErrNotImplemented gets returned if something doesnt work yet
	ErrNotImplemented = errors.New("Not implemented")
)

// SteamCmd is a wrapper around the Steam CMD
type SteamCmd struct {
	sync.Mutex // mutex for operations

	SteamCmdDir string
	AppBasePath string

	LoginUser string
	LoginPass string

	Debug bool
}

// New creates a new steamcmd instance, if path is empty, a temporary
// path will be created. Otherwise an existing instance will be reused.
func New(user, pass) *SteamCmd {
	var err error

	// if user not specified correctly, default to anonymous login
	if user == "" || pass == "" {
		user = "anonymous"
		pass = ""
	}

	scmd := &SteamCmd{
		SteamCmdDir: path,
		AppBasePath: gamesPath,
		LoginUser:   user,
		LoginPass:   pass,
	}

	return scmd
}

// CheckLogin checks if the given login can authenticate with Steam
func (scmd SteamCmd) CheckLogin() error {
	return scmd.run("")
}

// EnsureInstalled checks if the SteamCmd is executable.
// Remember, steam needs curl, bzip2, tar and lib32gcc1
func (scmd SteamCmd) EnsureInstalled() error {
	_, err := exec.LookPath("steamcmd")
	return err
}

// GetAppPath returns the path where an app would be installed
func (scmd *SteamCmd) GetAppPath(id int) string {
	return filepath.Join(scmd.AppBasePath, strconv.Itoa(id))
}

// InstallUpdateApp installs and updates a given app
func (scmd *SteamCmd) InstallUpdateApp(id int) error {
	return scmd.run("+app_update", strconv.Itoa(id))
}

// AppInstalledVersion returns the Build ID of an Steam App
func (scmd *SteamCmd) AppInstalledVersion(id int) (int, error) {
	return 0, ErrNotImplemented
}

// AppAvailableVersion returns the latest Build ID for the Public branch
// of an Steam App
func (scmd *SteamCmd) AppAvailableVersion(id int) (int, error) {
	return 0, ErrNotImplemented
}

// DownloadWorkshopMod tries to download a mod from the workshop
func (scmd *SteamCmd) DownloadWorkshopMod(appid, id int) error {
	return scmd.run("+workshop_download_item", strconv.Itoa(appid), strconv.Itoa(id))
}

// run helper
// * exit status 8 - no subscription
func (scmd *SteamCmd) run(params ...string) error {
	if scmd.LoginUser == "" || scmd.LoginUser == "anonymous" || scmd.LoginPass == "" {
		loginParams := []string{"+login", "anonymous"}
		params = append(loginParams, params...)
	} else {
		loginParams := []string{"+login", scmd.LoginUser, scmd.LoginPass}
		params = append(loginParams, params...)
	}
	params = append(params, "+quit")

	task := exec.Command("steamcmd", params...)
	if scmd.Debug {
		task.Stdout = os.Stdout
		task.Stderr = os.Stderr
	}
	if err := task.Run(); err != nil {
		return errors.Wrap(err, "raw command failed")
	}

	return nil
}
