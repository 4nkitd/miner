package systray

import (
	"fmt"
	"os/exec"
	"runtime"
	
	"github.com/getlantern/systray"
)

type App struct {
	server      ServerInterface
	hosts       HostsInterface
	cli         CLIInterface
	service     ServiceInterface
	cfg         ConfigInterface
	menuItems   *MenuItems
}

type MenuItems struct {
	openAdminer *systray.MenuItem
	startStop   *systray.MenuItem
	autoStart   *systray.MenuItem
	uninstall   *systray.MenuItem
}

type ServerInterface interface {
	Start() error
	Stop() error
	IsRunning() bool
	URL() string
}

type HostsInterface interface {
	AddEntry(domain, ip string) error
	RemoveEntry(domain string) error
}

type CLIInterface interface {
	Register() error
	Unregister() error
}

type ServiceInterface interface {
	Install() error
	Uninstall() error
	Status() (string, error)
}

type ConfigInterface interface {
	URL() string
}

func NewApp(server ServerInterface, hosts HostsInterface, cli CLIInterface, service ServiceInterface, cfg ConfigInterface) *App {
	return &App{
		server:  server,
		hosts:   hosts,
		cli:     cli,
		service: service,
		cfg:     cfg,
	}
}

func (a *App) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetTitle("Miner")
	systray.SetTooltip("Miner - Database Manager")
	
	a.menuItems = &MenuItems{}
	
	a.menuItems.openAdminer = systray.AddMenuItem("Open Adminer", "Open Adminer in browser")
	systray.AddSeparator()
	a.menuItems.startStop = systray.AddMenuItem("Stop Server", "Stop the Adminer server")
	a.menuItems.autoStart = systray.AddMenuItemCheckbox("Auto-start on Boot", "Start Miner automatically", true)
	systray.AddSeparator()
	a.menuItems.uninstall = systray.AddMenuItem("Uninstall", "Remove Miner configuration")
	mQuit := systray.AddMenuItem("Quit", "Quit Miner")
	
	go a.handleEvents(mQuit)
}

func (a *App) handleEvents(mQuit *systray.MenuItem) {
	for {
		select {
		case <-a.menuItems.openAdminer.ClickedCh:
			a.openBrowser()
		case <-a.menuItems.startStop.ClickedCh:
			a.toggleServer()
		case <-a.menuItems.autoStart.ClickedCh:
			a.toggleAutoStart()
		case <-a.menuItems.uninstall.ClickedCh:
			a.uninstall()
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func (a *App) openBrowser() {
	url := a.cfg.URL()
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

func (a *App) toggleServer() {
	if a.server.IsRunning() {
		if err := a.server.Stop(); err != nil {
			fmt.Printf("Failed to stop server: %v\n", err)
			return
		}
		a.menuItems.startStop.SetTitle("Start Server")
		fmt.Println("Server stopped")
	} else {
		if err := a.server.Start(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			return
		}
		a.menuItems.startStop.SetTitle("Stop Server")
		fmt.Println("Server started")
	}
}

func (a *App) toggleAutoStart() {
	if a.menuItems.autoStart.Checked() {
		a.menuItems.autoStart.Uncheck()
		if err := a.service.Uninstall(); err != nil {
			fmt.Printf("Failed to disable auto-start: %v\n", err)
		} else {
			fmt.Println("Auto-start disabled")
		}
	} else {
		a.menuItems.autoStart.Check()
		if err := a.service.Install(); err != nil {
			fmt.Printf("Failed to enable auto-start: %v\n", err)
		} else {
			fmt.Println("Auto-start enabled")
		}
	}
}

func (a *App) uninstall() {
	fmt.Println("Uninstalling Miner...")
	
	if err := a.server.Stop(); err != nil {
		fmt.Printf("Warning: Failed to stop server: %v\n", err)
	}
	
	if err := a.hosts.RemoveEntry("miner.local"); err != nil {
		fmt.Printf("Warning: Failed to remove hosts entry: %v\n", err)
	}
	
	if err := a.cli.Unregister(); err != nil {
		fmt.Printf("Warning: Failed to unregister CLI commands: %v\n", err)
	}
	
	if err := a.service.Uninstall(); err != nil {
		fmt.Printf("Warning: Failed to uninstall service: %v\n", err)
	}
	
	fmt.Println("Miner uninstalled successfully")
	systray.Quit()
}

func (a *App) onExit() {
	if a.server != nil && a.server.IsRunning() {
		a.server.Stop()
	}
	fmt.Println("Miner exited")
}
