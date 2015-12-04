package compose

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/cli/logger"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/rancher/os/config"
	rosDocker "github.com/rancher/os/docker"
	"github.com/rancher/os/util"
)

func CreateService(cfg *config.CloudConfig, name string, serviceConfig *project.ServiceConfig) (project.Service, error) {
	if cfg == nil {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return nil, err
		}
	}

	p, err := RunServiceSet("once", cfg, map[string]*project.ServiceConfig{
		name: serviceConfig,
	})
	if err != nil {
		return nil, err
	}

	return p.CreateService(name)
}

func RunServiceSet(name string, cfg *config.CloudConfig, configs map[string]*project.ServiceConfig) (*project.Project, error) {
	p, err := newProject(name, cfg)
	if err != nil {
		return nil, err
	}

	addServices(p, map[interface{}]interface{}{}, configs)

	return p, p.Up()
}

func GetProject(cfg *config.CloudConfig, networkingAvailable bool) (*project.Project, error) {
	return newCoreServiceProject(cfg, networkingAvailable)
}

func newProject(name string, cfg *config.CloudConfig) (*project.Project, error) {
	clientFactory, err := rosDocker.NewClientFactory(docker.ClientOpts{})
	if err != nil {
		return nil, err
	}

	serviceFactory := &rosDocker.ServiceFactory{
		Deps: map[string][]string{},
	}
	context := &docker.Context{
		ClientFactory: clientFactory,
		Context: project.Context{
			ProjectName:       name,
			NoRecreate:        true, // for libcompose to not recreate on project reload, looping up the boot :)
			EnvironmentLookup: rosDocker.NewConfigEnvironment(cfg),
			ServiceFactory:    serviceFactory,
			Log:               cfg.Rancher.Log,
			LoggerFactory:     logger.NewColorLoggerFactory(),
		},
	}
	serviceFactory.Context = context

	return docker.NewProject(context)
}

func addServices(p *project.Project, enabled map[interface{}]interface{}, configs map[string]*project.ServiceConfig) map[interface{}]interface{} {
	// Note: we ignore errors while loading services
	unchanged := true
	for name, serviceConfig := range configs {
		hash := project.GetServiceHash(name, serviceConfig)

		if enabled[name] == hash {
			continue
		}

		if err := p.AddConfig(name, serviceConfig); err != nil {
			log.Infof("Failed loading service %s", name)
			continue
		}

		if unchanged {
			enabled = util.MapCopy(enabled)
			unchanged = false
		}
		enabled[name] = hash
	}
	return enabled
}

func newCoreServiceProject(cfg *config.CloudConfig, network bool) (*project.Project, error) {
	projectEvents := make(chan project.Event)
	enabled := map[interface{}]interface{}{}

	p, err := newProject("os", cfg)
	if err != nil {
		return nil, err
	}

	p.AddListener(project.NewDefaultListener(p))
	p.AddListener(projectEvents)

	p.ReloadCallback = func() error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return err
		}

		enabled = addServices(p, enabled, cfg.Rancher.Services)

		for service, serviceEnabled := range cfg.Rancher.ServicesInclude {
			if _, ok := enabled[service]; ok || !serviceEnabled {
				continue
			}

			bytes, err := LoadServiceResource(service, network, cfg)
			if err != nil {
				if err == util.ErrNoNetwork {
					log.Debugf("Can not load %s, networking not enabled", service)
				} else {
					log.Errorf("Failed to load %s : %v", service, err)
				}
				continue
			}

			err = p.Load(bytes)
			if err != nil {
				log.Errorf("Failed to load %s : %v", service, err)
				continue
			}

			enabled[service] = service
		}

		return nil
	}

	go func() {
		for event := range projectEvents {
			if event.EventType == project.EventContainerStarted && event.ServiceName == "network" {
				network = true
			}
		}
	}()

	err = p.ReloadCallback()
	if err != nil {
		log.Errorf("Failed to reload os: %v", err)
		return nil, err
	}

	return p, nil
}

func LoadServiceResource(name string, network bool, cfg *config.CloudConfig) ([]byte, error) {
	return util.LoadResource(name, network, cfg.Rancher.Repositories.ToArray())
}
