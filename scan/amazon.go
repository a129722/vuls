package scan

import (
	"github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/models"
	"github.com/future-architect/vuls/util"
	"golang.org/x/xerrors"
)

// inherit OsTypeInterface
type amazon struct {
	redhatBase
}

// NewAmazon is constructor
func newAmazon(c config.ServerInfo) *amazon {
	r := &amazon{
		redhatBase{
			base: base{
				osPackages: osPackages{
					Packages:  models.Packages{},
					VulnInfos: models.VulnInfos{},
				},
			},
			sudo: rootPrivAmazon{},
		},
	}
	r.log = util.NewCustomLogger(c)
	r.setServerInfo(c)
	return r
}

func (o *amazon) checkScanMode() error {
	if o.getServerInfo().Mode.IsOffline() {
		return xerrors.New("Remove offline scan mode, Amazon needs internet connection")
	}
	return nil
}

func (o *amazon) checkDeps() error {
	if o.getServerInfo().Mode.IsFast() {
		return o.execCheckDeps(o.depsFast())
	} else if o.getServerInfo().Mode.IsFastRoot() {
		return o.execCheckDeps(o.depsFastRoot())
	} else if o.getServerInfo().Mode.IsDeep() {
		return o.execCheckDeps(o.depsDeep())
	}
	return xerrors.New("Unknown scan mode")
}

func (o *amazon) depsFast() []string {
	if o.getServerInfo().Mode.IsOffline() {
		return []string{}
	}
	// repoquery
	return []string{"yum-utils"}
}

func (o *amazon) depsFastRoot() []string {
	return []string{
		"yum-utils",
		"yum-plugin-ps",
	}
}

func (o *amazon) depsDeep() []string {
	return o.depsFastRoot()
}

func (o *amazon) checkIfSudoNoPasswd() error {
	if o.getServerInfo().Mode.IsFast() {
		return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsFast())
	} else if o.getServerInfo().Mode.IsFastRoot() {
		return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsFastRoot())
	} else {
		return o.execCheckIfSudoNoPasswd(o.sudoNoPasswdCmdsDeep())
	}
}

func (o *amazon) sudoNoPasswdCmdsFast() []cmd {
	return []cmd{}
}

func (o *amazon) sudoNoPasswdCmdsFastRoot() []cmd {
	return []cmd{
		{"yum -q ps all --color=never", exitStatusZero},
		{"stat /proc/1/exe", exitStatusZero},
		{"needs-restarting", exitStatusZero},
		{"which which", exitStatusZero},
	}
}

func (o *amazon) sudoNoPasswdCmdsDeep() []cmd {
	return o.sudoNoPasswdCmdsFastRoot()
}

type rootPrivAmazon struct{}

func (o rootPrivAmazon) repoquery() bool {
	return false
}

func (o rootPrivAmazon) yumRepolist() bool {
	return false
}

func (o rootPrivAmazon) yumUpdateInfo() bool {
	return false
}

func (o rootPrivAmazon) yumChangelog() bool {
	return false
}

func (o rootPrivAmazon) yumMakeCache() bool {
	return false
}
