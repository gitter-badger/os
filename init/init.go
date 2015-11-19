// +build linux

package init

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
	"github.com/rancher/os/util"
)

const (
	STATE string = "/state"
)

func loadModules(cfg *config.CloudConfig) (*config.CloudConfig, error) {
	mounted := map[string]bool{}

	f, err := os.Open("/proc/modules")
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	reader := bufio.NewScanner(f)
	for reader.Scan() {
		mounted[strings.SplitN(reader.Text(), " ", 2)[0]] = true
	}

	for _, module := range cfg.Rancher.Modules {
		if mounted[module] {
			continue
		}

		log.Debugf("Loading module %s", module)
		if err := exec.Command("modprobe", module).Run(); err != nil {
			log.Errorf("Could not load module %s, err %v", module, err)
		}
	}

	return cfg, nil
}

func sysInit(c *config.CloudConfig) (*config.CloudConfig, error) {
	args := append([]string{config.SYSINIT_BIN}, os.Args[1:]...)

	cmd := &exec.Cmd{
		Path: config.ROS_BIN,
		Args: args,
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return c, err
	}

	return c, os.Stdin.Close()
}

func MainInit() {
	if err := RunInit(); err != nil {
		log.Fatal(err)
	}
}

func mountState(cfg *config.CloudConfig) error {
	var err error

	if cfg.Rancher.State.Dev == "" {
		return nil
	}

	dev := util.ResolveDevice(cfg.Rancher.State.Dev)
	if dev == "" {
		return fmt.Errorf("Could not resolve device %q", cfg.Rancher.State.Dev)
	}
	fsType := cfg.Rancher.State.FsType
	if fsType == "auto" {
		fsType, err = util.GetFsType(dev)
	}

	if err != nil {
		return err
	}

	log.Debugf("FsType has been set to %s", fsType)
	log.Infof("Mounting state device %s to %s", dev, STATE)
	return util.Mount(dev, STATE, fsType, "")
}

func tryMountState(cfg *config.CloudConfig) error {
	if mountState(cfg) == nil {
		return nil
	}

	// If we failed to mount lets run bootstrap and try again
	if err := bootstrap(cfg); err != nil {
		return err
	}

	return mountState(cfg)
}

func tryMountAndBootstrap(cfg *config.CloudConfig) (*config.CloudConfig, error) {
	if err := tryMountState(cfg); !cfg.Rancher.State.Required && err != nil {
		return cfg, nil
	} else if err != nil {
		return cfg, err
	}

	log.Debugf("Switching to new root at %s", STATE)
	return cfg, switchRoot(STATE, cfg.Rancher.RmUsr)
}

func getCgroupArgs(cgroupHierarchy map[string]string) []string {
	args := []string{}

	for k, v := range cgroupHierarchy {
		args = append(args, "--dfs-cgroup="+k+":"+v)
	}

	return args
}

func getLaunchArgs(cfg *config.CloudConfig, dockerCfg *config.DockerConfig) []string {

	args := []string{config.DOCKER_BIN}

	if len(cfg.Rancher.Network.Dns.Nameservers) > 0 {
		args = append(args, "--dfs-dns-nameservers="+strings.Join(cfg.Rancher.Network.Dns.Nameservers, ","))
	}
	if len(cfg.Rancher.Network.Dns.Search) > 0 {
		args = append(args, "--dfs-dns-search="+strings.Join(cfg.Rancher.Network.Dns.Search, ","))
	}
	args = append(args, "--dfs-emulate-systemd")

	if !cfg.Rancher.Debug {
		args = append(args, "--dfs-logfile="+config.SYSTEM_DOCKER_LOG)
	}

	args = append(args, dockerCfg.Args...)
	args = append(args, dockerCfg.ExtraArgs...)

	return args
}

func RunInit() error {
	os.Setenv("PATH", "/sbin:/usr/sbin:/usr/bin")
	// Magic setting to tell Docker to do switch_root and not pivot_root
	os.Setenv("DOCKER_RAMDISK", "true")

	initFuncs := []config.CfgFunc{
		func(_ *config.CloudConfig) (*config.CloudConfig, error) {
			cfg, err := config.LoadConfig()
			if err != nil {
				return cfg, err
			}

			if cfg.Rancher.Debug {
				cfgString, err := config.Dump(false, false, true)
				if err != nil {
					log.WithFields(log.Fields{"err": err}).Error("Error serializing config")
				} else {
					log.Debugf("Config: %s", cfgString)
				}
			}

			return cfg, nil
		},
		loadModules,
		tryMountAndBootstrap,
		func(_ *config.CloudConfig) (*config.CloudConfig, error) {
			return config.LoadConfig()
		},
		loadModules,
		sysInit,
	}

	cgroupHierarchy := map[string]string{
		"cpu":      "cpu",
		"cpuacct":  "cpu",
		"net_cls":  "net_cls",
		"net_prio": "net_cls",
	}
	cmd := exec.Command("/usr/bin/true", getCgroupArgs(cgroupHierarchy)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cfg, err := config.ChainCfgFuncs(nil, initFuncs...)
	if err != nil {
		return err
	}

	args := getLaunchArgs(cfg, &cfg.Rancher.SystemDocker)
	log.Info("Launching System Docker")
	return syscall.Exec(config.DOCKERLAUNCH_BIN, args, cfg.Rancher.SystemDocker.Environment)
}
