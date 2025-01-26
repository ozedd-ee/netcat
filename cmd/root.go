package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ozedd-ee/netcat/src"
	"github.com/spf13/cobra"
)

var listen bool
var UDP bool
var zeroIO bool
var exec bool
var hex bool
var port string

var host = "127.0.0.1"

var rootCmd = &cobra.Command{
	Use:     `netcat`,
	Short:   `A lite clone of the Netcat command line utility`,
	Long:    `Add a longer description`,
	Example: `netcat -l -p 8080`,

	Run: func(cmd *cobra.Command, args []string) {
		if listen {
			if UDP {
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()
				go src.UDPListen(ctx, "udp", port)
	
				<-ctx.Done()
				time.Sleep(1 * time.Second) // Wait for other processes
				fmt.Println("\nShutting down server...")
			} else if exec {
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()
				if len(args) == 0 {
					log.Fatal("Shell not specified")
				}
				go src.Exec(port, args[0])
	
				<-ctx.Done()
				time.Sleep(1 * time.Second)
				fmt.Println("\nShutting down server...")
			} else if hex {
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()
				go src.ListenHex(port)

				<-ctx.Done()
				time.Sleep(1 * time.Second)
				fmt.Println("\nShutting down server...")
			} else {
				// Create a context that cancels on interrupt signals
				ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
				defer stop()
				go src.TCPListen("tcp", port)

				// Wait for a termination signal
				<-ctx.Done()
				time.Sleep(1 * time.Second) // Wait for other processes
				fmt.Println("\nShutting down server...")
			}
		} else if zeroIO {
			if len(args) != 0 {
				hostname := args[0]
				src.Ping(hostname, port)
			} else {
				src.Ping(host, port)
			}
		} else {
			if UDP {
				if len(args) == 0 {
					log.Fatal("Hostname not specified")
				}
				src.UdpConnect(args[0], port)
			} else {
				if len(args) == 0 {
					log.Fatal("Hostname not specified")
				}
				src.TcpConnect(args[0], port)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&listen, "listen", "l", false, "Run in listening mode for TCP connections")

	rootCmd.PersistentFlags().BoolVarP(&UDP, "udp", "u", false, "Run in listening mode for UDP connections")

	rootCmd.PersistentFlags().BoolVarP(&zeroIO, "zeroio", "z", false, "Zero-I/O mode, report connection status only")

	rootCmd.PersistentFlags().BoolVarP(&exec, "exec", "e", false, "Exec mode")

	rootCmd.PersistentFlags().BoolVarP(&hex, "hex", "x", false, "Dump transmitted data as hex to output")

	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Host port to use")
}
