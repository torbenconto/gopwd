package api

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"

	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/pwgen"
	"github.com/torbenconto/gopwd/internal/ssl"
	"github.com/torbenconto/gopwd/internal/util"
)

func setupRouter(vaultPath string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust this to your needs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/list", func(c *gin.Context) {
		services, err := io.ListServices(vaultPath)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error listing services: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"services": services,
		})
	})

	r.POST("/get", func(c *gin.Context) {
		var req struct {
			Service     string `json:"service"`
			GpgPassword string `json:"gpg_password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request",
			})
			return
		}

		// Check if service exists
		if !io.Exists(filepath.Join(vaultPath, req.Service+".gpg")) {
			c.JSON(400, gin.H{
				"message": "service doesn't exist",
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
				"message": "error reading gpg-id: " + err.Error(),
			})
			return
		}

		args := []string{"--quiet", "--yes", "--compress-algo=none", "--no-encrypt-to", "--no-auto-check-trustdb", "--batch", "--pinentry-mode=loopback"}
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
			c.JSON(500, gin.H{
				"message": "error decrypting file: " + err.Error(),
			})
			return
		}

		// Successfully decrypted, return the content
		c.JSON(200, gin.H{
			"password": string(decrypted),
		})
	})
	r.POST("/update", func(c *gin.Context) {
		var req struct {
			Service    string `json:"service"`
			NewContent string `json:"new_content"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request",
			})
			return
		}

		// Check if service exists
		if !io.Exists(filepath.Join(vaultPath, req.Service+".gpg")) {
			c.JSON(400, gin.H{
				"message": "service doesn't exist",
			})
			return
		}

		// Encrypt the password
		gpgID, err := util.ReadGPGID(filepath.Join(vaultPath, ".gpg-id"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error reading gpg-id",
			})
			return
		}

		// Initialize GPG module with configuration
		gpgModule := gpg.NewGPG(gpgID, gpg.Config{})

		// Encrypt the password
		encrypted, err := gpgModule.Encrypt([]byte(req.NewContent))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error encrypting password",
			})
			return
		}

		// Write the encrypted password to the file
		err = io.WriteFile(filepath.Join(vaultPath, req.Service+".gpg"), encrypted)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error writing file",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "password updated",
		})
	})
	r.POST("/delete", func(c *gin.Context) {
		var req struct {
			Service string `json:"service"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request",
			})
			return
		}

		// Check if service exists
		if !io.Exists(filepath.Join(vaultPath, req.Service+".gpg")) {
			c.JSON(400, gin.H{
				"message": "service doesn't exist",
			})
			return
		}

		err := os.Remove(filepath.Join(vaultPath, req.Service+".gpg"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error deleting file: " + err.Error(),
			})
			return
		}

		dirPath := path.Dir(filepath.Join(vaultPath, req.Service+".gpg"))
		isEmpty, err := io.IsDirEmpty(dirPath)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error checking if directory is empty: " + err.Error(),
			})
			return
		}
		if isEmpty && dirPath != vaultPath {
			err := os.Remove(dirPath)
			if err != nil {
				c.JSON(500, gin.H{
					"message": "error removing directory: " + err.Error(),
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"message": "file deleted",
		})
	})
	r.POST("/insert", func(c *gin.Context) {
		var req struct {
			Service string `json:"service"`
			Content string `json:"content"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request: " + err.Error(),
			})
			return
		}

		// Check if service already exists
		if io.Exists(filepath.Join(vaultPath, req.Service+".gpg")) {
			c.JSON(400, gin.H{
				"message": "service already exists",
			})
			return
		}

		// Encrypt the password
		gpgID, err := util.ReadGPGID(filepath.Join(vaultPath, ".gpg-id"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error reading gpg-id: " + err.Error(),
			})
			return
		}

		// Initialize GPG module with configuration
		gpgModule := gpg.NewGPG(gpgID, gpg.Config{})

		// Encrypt the password
		encrypted, err := gpgModule.Encrypt([]byte(req.Content))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error encrypting password: " + err.Error(),
			})
			return
		}

		err = util.CreateStructureAndClean(req.Service, vaultPath, filepath.Join(vaultPath, req.Service+".gpg"), encrypted)

		c.JSON(200, gin.H{
			"message": "password inserted",
		})
	})
	r.POST("/generate", func(c *gin.Context) {
		req := &struct {
			Service   string `json:"service"`
			Length    int    `json:"length"`
			Humanized bool   `json:"humanized"`
			Symbols   bool   `json:"symbols"`
			Numbers   bool   `json:"numbers"`
			Lowercase bool   `json:"lowercase"`
			Uppercase bool   `json:"uppercase"`
		}{
			Length:    16,
			Symbols:   true,
			Numbers:   true,
			Lowercase: true,
			Uppercase: true,
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request",
			})
			return
		}

		servicePath := path.Join(vaultPath, req.Service) + ".gpg"

		if io.Exists(servicePath) {
			c.JSON(400, gin.H{
				"message": "service already exists",
			})
			return
		}

		// Generate password
		passwordGenerator := pwgen.NewPasswordGenerator(pwgen.PasswordGeneratorConfig{
			Length:    req.Length,
			Humanized: req.Humanized,
			Symbols:   req.Symbols,
			Numbers:   req.Numbers,
			Lowercase: req.Lowercase,
			Uppercase: req.Uppercase,
		})

		password, err := passwordGenerator.Generate()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error generating password: " + err.Error(),
			})
			return
		}

		gpgID, err := util.ReadGPGID(filepath.Join(vaultPath, ".gpg-id"))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error reading gpg-id: " + err.Error(),
			})
			return
		}

		// Initialize GPG module with configuration
		gpgModule := gpg.NewGPG(gpgID, gpg.Config{})

		// Encrypt the password
		encrypted, err := gpgModule.Encrypt([]byte(password))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error encrypting password: " + err.Error(),
			})
			return
		}

		err = util.CreateStructureAndClean(req.Service, vaultPath, servicePath, encrypted)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error creating structure: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message":  "password generated and inserted",
			"password": password,
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
	gin.SetMode(gin.ReleaseMode)
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

func Run(gopwdPath, vaultPath, addr, certPath, keyPath string) error {
	gin.SetMode(gin.ReleaseMode)
	// Check if SSL certificates exist, and generate them if they don't
	if !io.Exists(certPath) || !io.Exists(keyPath) {
		err := ssl.GenerateSSLCert(certPath, keyPath)
		if err != nil {
			return err
		}
	}

	r := setupRouter(vaultPath)

	// Capture shutdown signals to gracefully stop the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		fmt.Println("Shutting down server...")
		os.Exit(0)
	}()

	fmt.Printf("Starting API server on %s...\n", addr)
	err := r.RunTLS(addr, certPath, keyPath)
	if err != nil {
		return err
	}

	return nil
}
