package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/skygenesisenterprise/aether-registry/server/src/config"
	"github.com/skygenesisenterprise/aether-registry/server/src/middleware"
	"github.com/skygenesisenterprise/aether-registry/server/src/routes"
	"github.com/skygenesisenterprise/aether-registry/server/src/services"
)

func displayBanner() {
	fmt.Printf("\n")
	fmt.Printf("\033[1;36m    ██╗    ██╗██╗  ██╗ █████╗ ████████╗██╗  ██╗███████╗████████╗\n")
	fmt.Printf("\033[1;36m    ██║    ██║██║  ██║██╔══██╗╚══██╔══╝██║  ██║██╔════╝╚══██╔══╝\n")
	fmt.Printf("\033[1;36m    ██║ █╗ ██║███████║███████║   ██║   ███████║█████╗     ██║   \n")
	fmt.Printf("\033[1;36m    ██║███╗██║██╔══██║██╔══██║   ██║   ██╔══██║██╔══╝     ██║   \n")
	fmt.Printf("\033[1;36m    ╚███╔███╔╝██║  ██║██║  ██║   ██║   ██║  ██║███████╗   ██║   \n")
	fmt.Printf("\033[1;36m     ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚══════╝   ╚═╝   \n")
	fmt.Printf("\033[0;37m")
	fmt.Printf("\n")
	fmt.Printf("\033[1;33m    ╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("\033[1;33m    ║                    Aether Registry System                    ║\n")
	fmt.Printf("\033[1;33m    ║               Enterprise Package Management                  ║\n")
	fmt.Printf("\033[1;33m    ║                   Version 1.0.0-alpha                        ║\n")
	fmt.Printf("\033[1;33m    ╚══════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("\033[0;37m")
	fmt.Printf("\n")
	fmt.Printf("\033[1;32m[✓] System Architecture: %s\033[0m\n", runtime.GOARCH)
	fmt.Printf("\033[1;32m[✓] Operating System: %s\033[0m\n", runtime.GOOS)
	fmt.Printf("\033[1;32m[✓] Go Version: %s\033[0m\n", runtime.Version())
	fmt.Printf("\033[1;32m[✓] CPU Cores: %d\033[0m\n", runtime.NumCPU())
	fmt.Printf("\033[1;32m[✓] Process ID: %d\033[0m\n", os.Getpid())
	fmt.Printf("\n")
}

func main() {
	displayBanner()

	cfg := config.Load()

	prismaService, err := services.NewPrismaService(cfg)
	if err != nil {
		fmt.Printf("\033[1;33m[!] Warning: Database connection failed: %v\033[0m\n", err)
		fmt.Printf("\033[1;33m[!] Running in mock mode\033[0m\n")
	} else {
		fmt.Printf("\033[1;32m[✓] Database connected\033[0m\n")
		defer prismaService.Close()
	}

	jwtService := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.Issuer)
	serviceKeyService := services.NewServiceKeyService()

	fmt.Printf("\033[1;32m[✓] Service keys initialized\033[0m\n")
	fmt.Printf("\033[1;33m[!] Default service keys generated (check logs for sk_ keys)\033[0m\n")

	router := gin.Default()

	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	routes.SetupRoutesWithServiceKey(router, jwtService, serviceKeyService)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("\033[1;32m[✓] Server starting on %s\033[0m\n", addr)
	fmt.Printf("\033[1;36m[✓] API available at http://localhost%s/api/v1\033[0m\n", addr)
	fmt.Printf("\033[1;33m[!] Service key auth enabled\033[0m\n")
	fmt.Printf("\n")

	if err := router.Run(addr); err != nil {
		fmt.Printf("\033[1;31m[✗] Failed to start server: %v\033[0m\n", err)
		os.Exit(1)
	}
}
