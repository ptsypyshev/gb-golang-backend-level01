package conf

import (
	"flag"
	"fmt"
	"os"
)

// GetConfig reads arguments from flags
func GetConfig() map[string]string {
	host := flag.String("host", "localhost", "Specify server IP address or domain name")
	port := flag.String("port", "8000", "Specify server TCP port")
	proto := flag.String("proto", "tcp", "Specify transport protocol - tcp or udp")
	isHelpFlagged := flag.Bool("help", false, "Prints help page")
	isShortHelpFlagged := flag.Bool("h", false, "Alias for help")
	flag.Parse()

	if *isHelpFlagged || *isShortHelpFlagged {
		fmt.Println("Usage: main [arguments]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	config := make(map[string]string)
	config["host"] = *host
	config["port"] = *port
	config["proto"] = *proto
	return config
}
