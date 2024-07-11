package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"github.com/torbenconto/gopwd/internal/io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func setupRouter(vaultPath string) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/list", func(c *gin.Context) {
		services, err := io.ListServices(vaultPath)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error listing services",
			})
			return
		}

		c.JSON(200, gin.H{
			"services": services,
		})
	})

	return r
}

func start(vaultPath string, addr string) {
	r := setupRouter(vaultPath)
	r.Run(addr)
}

func RunDaemon(gopwdPath, vaultPath, addr string, cmd []string) {
	daemonContext := &daemon.Context{
		PidFileName: filepath.Join(gopwdPath, "gopwd.pid"),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(gopwdPath, "gopwd.log"),
		LogFilePerm: 0640,
		WorkDir:     gopwdPath,
		Umask:       027,
		Args:        cmd,
	}

	d, err := daemonContext.Reborn()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if d != nil {
		return
	}
	defer daemonContext.Release()

	fmt.Println("Daemon started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		fmt.Println("Daemon terminated")
		os.Exit(0)
	}()

	start(vaultPath, addr)
}
