package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/ssl"
	"github.com/torbenconto/gopwd/internal/util"
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

	type getServiceRequest struct {
		Service     string `json:"service"`
		GpgPassword string `json:"gpg_password"`
	}
	r.POST("/get", func(c *gin.Context) {
		var req getServiceRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request",
			})
			return
		}

		file, err := io.ReadFile(filepath.Join(vaultPath, req.Service+".gpg"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error reading file",
			})
			return
		}

		// Decrypt file
		gpgID, err := util.ReadGPGID(filepath.Join(vaultPath, ".gpg-id"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error reading gpg-id",
			})
			return
		}

		args := []string{"--yes", "--compress-algo=none", "--no-encrypt-to", "--no-auto-check-trustdb", "--batch"}
		if req.GpgPassword != "" {
			args = append(args, "--passphrase", req.GpgPassword)
		}

		// Initialize GPG module with configuration
		gpgModule := gpg.NewGPG(gpgID, gpg.Config{
			Args: args,
		})

		// Attempt to decrypt the file
		decrypted, err := gpgModule.Decrypt(file)
		if err != nil {
			// Log the error for diagnostics
			fmt.Printf("Error decrypting file: %v\n", err)
			c.JSON(500, gin.H{
				"message": "error decrypting file",
			})
			return
		}

		// Successfully decrypted, return the content
		c.JSON(200, gin.H{
			"password": string(decrypted),
		})
	})

	return r
}

func start(vaultPath, addr, certPath, keyPath string) {
	// Check if cert exists
	if !io.Exists(certPath) || !io.Exists(keyPath) {
		// Generate cert
		err := ssl.GenerateSSLCert(certPath, keyPath)
		if err != nil {
			fmt.Println("Error generating SSL cert:", err)
			return
		}
	}

	r := setupRouter(vaultPath)
	err := r.RunTLS(addr, certPath, keyPath)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}

func RunDaemon(gopwdPath, vaultPath, addr string, cmd []string, certPath, keyPath string) {
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

	start(vaultPath, addr, certPath, keyPath)
}
