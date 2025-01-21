package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozedd-ee/netcat/src"
	"github.com/spf13/cobra"
)

var listen bool
var listenUDP bool
var port string

var rootCmd = &cobra.Command{
	Use:     `netcat`,
	Short:   `A lite clone of the Netcat command line utility`,
	Long:    `Add a longer description`,
	Example: `netcat -l -p 8080`,

	Run: func(cmd *cobra.Command, args []string) {
		if listen {
			// Create a context that cancels on interrupt signals
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			go src.TCPListen("tcp", port)

			// Wait for a termination signal
			<-ctx.Done()
			fmt.Println("\nShutting down server...")
		} else if listenUDP {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			go src.UDPListen(ctx, "udp", port)

			<-ctx.Done()
			fmt.Println("\nShutting down server...")
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

	rootCmd.PersistentFlags().BoolVarP(&listenUDP, "udp", "u", false, "Run in listening mode for UDP connections")

	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Host port to use")
}
