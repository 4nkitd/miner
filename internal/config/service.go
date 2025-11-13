package config

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/4nkitd/miner/internal/assets"
	"github.com/4nkitd/miner/internal/hosts"
	"github.com/4nkitd/miner/internal/server"
	"github.com/kardianos/service"
)

type MinerService struct {
	cfg *Config
}

func NewService(cfg *Config) (*MinerService, error) {
	return &MinerService{cfg: cfg}, nil
}

func (m *MinerService) baseConfig() *service.Config {
	return &service.Config{
		Name:        AppName,
		DisplayName: "Miner Database Manager",
		Description: "Adminer database manager powered by FrankenPHP",
	}
}

func (m *MinerService) Install() error {
	svcConfig := m.baseConfig()
	prg := &program{cfg: m.cfg}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	return s.Install()
}

func (m *MinerService) Uninstall() error {
	svcConfig := m.baseConfig()
	prg := &program{cfg: m.cfg}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	return s.Uninstall()
}

func (m *MinerService) Start() error {
	svcConfig := m.baseConfig()
	prg := &program{cfg: m.cfg}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	return s.Start()
}

func (m *MinerService) Status() (string, error) {
	svcConfig := m.baseConfig()
	prg := &program{cfg: m.cfg}
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

type program struct {
	cfg     *Config
	srv     *server.Server
	tempDir string
}

func (p *program) Start(s service.Service) error {
	// Load fresh config (ensures embedded asset extraction if needed)
	cfg, err := New()
	if err != nil {
		return err
	}
	p.cfg = cfg
	p.tempDir = cfg.TempAssets

	// Ensure hosts entry exists (service may start before install finished)
	hm := hosts.NewManager(cfg.HostsPath)
	if has, _ := hm.HasEntry(cfg.Domain); !has {
		_ = hm.AddEntry(cfg.Domain, cfg.Host)
	}

	p.srv = server.NewServer(cfg.Port, cfg.Domain, cfg.AssetsDir)
	if err := p.srv.Start(); err != nil {
		return fmt.Errorf("service server start failed: %w", err)
	}

	go p.waitSignals()
	return nil
}

func (p *program) waitSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	_ = p.Stop(nil)
}

func (p *program) Stop(s service.Service) error {
	if p.srv != nil && p.srv.IsRunning() {
		_ = p.srv.Stop()
	}
	if p.tempDir != "" {
		_ = assets.Cleanup(p.tempDir)
	}
	return nil
}
