package config

import (
	"github.com/kardianos/service"
)

type MinerService struct {
	cfg *Config
}

func NewService(cfg *Config) (*MinerService, error) {
	return &MinerService{cfg: cfg}, nil
}

func (m *MinerService) Install() error {
	svcConfig := &service.Config{
		Name:        AppName,
		DisplayName: "Miner Database Manager",
		Description: "Adminer database manager powered by FrankenPHP",
	}
	
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	
	return s.Install()
}

func (m *MinerService) Uninstall() error {
	svcConfig := &service.Config{
		Name:        AppName,
		DisplayName: "Miner Database Manager",
		Description: "Adminer database manager powered by FrankenPHP",
	}
	
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	
	return s.Uninstall()
}

func (m *MinerService) Status() (string, error) {
	svcConfig := &service.Config{
		Name:        AppName,
		DisplayName: "Miner Database Manager",
		Description: "Adminer database manager powered by FrankenPHP",
	}
	
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return "", err
	}
	
	status, err := s.Status()
	if err != nil {
		return "Unknown", err
	}
	
	switch status {
	case service.StatusRunning:
		return "Running", nil
	case service.StatusStopped:
		return "Stopped", nil
	default:
		return "Unknown", nil
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
