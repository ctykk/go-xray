package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ctykk/go-xray"
	"github.com/ctykk/go-xray/proxy/shadowsocks"
	"github.com/ctykk/go-xray/proxy/trojan"
	"github.com/ctykk/go-xray/proxy/vless"
	"github.com/ctykk/go-xray/proxy/vmess"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "xray",
		Short: "Xray proxy gateway - start HTTPS proxy through various proxy protocols",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(
		createShadowsocksCmd(),
		createTrojanCmd(),
		createVlessCmd(),
		createVmessCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// createShadowsocksCmd creates the shadowsocks subcommand
func createShadowsocksCmd() *cobra.Command {
	var (
		serverHost string
		serverPort uint16
		password   string
		cipher     string
		proxyPort  uint16
	)

	cmd := &cobra.Command{
		Use:   "shadowsocks",
		Short: "Start HTTP proxy through Shadowsocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate cipher
			cipherEnum, err := shadowsocks.ParseCipher(cipher)
			if err != nil {
				return fmt.Errorf("invalid cipher '%s': %w", cipher, err)
			}

			// Create proxy instance
			ss, err := shadowsocks.New(serverHost, serverPort, cipherEnum, password, "xray-ss")
			if err != nil {
				return fmt.Errorf("failed to create Shadowsocks proxy: %w", err)
			}

			return runProxy(ss, proxyPort)
		},
	}

	cmd.Flags().StringVar(&serverHost, "host", "", "Shadowsocks server host (required)")
	cmd.Flags().Uint16Var(&serverPort, "port", 0, "Shadowsocks server port (required)")
	cmd.Flags().StringVar(&password, "password", "", "Shadowsocks password (required)")
	cmd.Flags().StringVar(&cipher, "cipher", "AES-256-GCM", "Encryption cipher (AES-128-GCM, AES-256-GCM, CHACHA20-POLY1305, XCHACHA20-POLY1305)")
	cmd.Flags().Uint16Var(&proxyPort, "proxy-port", 0, "Local HTTP proxy port (required)")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("proxy-port")

	return cmd
}

// createTrojanCmd creates the trojan subcommand
func createTrojanCmd() *cobra.Command {
	var (
		serverHost string
		serverPort uint16
		password   string
		proxyPort  uint16
	)

	cmd := &cobra.Command{
		Use:   "trojan",
		Short: "Start HTTP proxy through Trojan",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create proxy instance
			tj, err := trojan.New(serverHost, serverPort, password, "xray-trojan")
			if err != nil {
				return fmt.Errorf("failed to create Trojan proxy: %w", err)
			}

			return runProxy(tj, proxyPort)
		},
	}

	cmd.Flags().StringVar(&serverHost, "host", "", "Trojan server host (required)")
	cmd.Flags().Uint16Var(&serverPort, "port", 0, "Trojan server port (required)")
	cmd.Flags().StringVar(&password, "password", "", "Trojan password (required)")
	cmd.Flags().Uint16Var(&proxyPort, "proxy-port", 0, "Local HTTP proxy port (required)")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("proxy-port")

	return cmd
}

// createVlessCmd creates the vless subcommand
func createVlessCmd() *cobra.Command {
	var (
		serverHost string
		serverPort uint16
		uuid       string
		encryption string
		proxyPort  uint16
	)

	cmd := &cobra.Command{
		Use:   "vless",
		Short: "Start HTTP proxy through VLESS",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create proxy instance
			vl, err := vless.New(serverHost, serverPort, uuid, encryption, "xray-vless")
			if err != nil {
				return fmt.Errorf("failed to create VLESS proxy: %w", err)
			}

			return runProxy(vl, proxyPort)
		},
	}

	cmd.Flags().StringVar(&serverHost, "host", "", "VLESS server host (required)")
	cmd.Flags().Uint16Var(&serverPort, "port", 0, "VLESS server port (required)")
	cmd.Flags().StringVar(&uuid, "uuid", "", "VLESS UUID (required)")
	cmd.Flags().StringVar(&encryption, "encryption", "none", "Encryption method (default: none)")
	cmd.Flags().Uint16Var(&proxyPort, "proxy-port", 0, "Local HTTP proxy port (required)")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("uuid")
	cmd.MarkFlagRequired("proxy-port")

	return cmd
}

// createVmessCmd creates the vmess subcommand
func createVmessCmd() *cobra.Command {
	var (
		serverHost string
		serverPort uint16
		uuid       string
		cipher     string
		proxyPort  uint16
	)

	cmd := &cobra.Command{
		Use:   "vmess",
		Short: "Start HTTP proxy through VMess",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate cipher
			cipherEnum, err := vmess.ParseCipher(cipher)
			if err != nil {
				return fmt.Errorf("invalid cipher '%s': %w", cipher, err)
			}

			// Create proxy instance
			vm, err := vmess.New(serverHost, serverPort, cipherEnum, uuid, "xray-vmess")
			if err != nil {
				return fmt.Errorf("failed to create VMess proxy: %w", err)
			}

			return runProxy(vm, proxyPort)
		},
	}

	cmd.Flags().StringVar(&serverHost, "host", "", "VMess server host (required)")
	cmd.Flags().Uint16Var(&serverPort, "port", 0, "VMess server port (required)")
	cmd.Flags().StringVar(&uuid, "uuid", "", "VMess UUID (required)")
	cmd.Flags().StringVar(&cipher, "cipher", "AES-128-GCM", "Encryption cipher (AES-128-GCM, CHACHA20-POLY1305, AUTO, NONE, ZERO)")
	cmd.Flags().Uint16Var(&proxyPort, "proxy-port", 0, "Local HTTP proxy port (required)")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("uuid")
	cmd.MarkFlagRequired("proxy-port")

	return cmd
}

// runProxy starts the HTTP proxy and handles graceful shutdown
func runProxy(p xray.Proxy, port uint16) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.HTTPProxy(ctx, port); err != nil {
		return fmt.Errorf("failed to start proxy: %w", err)
	}

	fmt.Printf("HTTP proxy listening on localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v, shutting down gracefully...\n", sig)
	cancel()

	fmt.Println("Proxy stopped successfully")
	return nil
}
