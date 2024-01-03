package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tiiuae/nollakonffi"
)

var (
	service  = flag.String("service", "_swarm._tcp", "Set the service category to look for devices.")
	domain   = flag.String("domain", "local", "Set the search domain. For local networks, default is fine.")
	waitTime = flag.Int("wait", 10, "Duration in [s] to run discovery.")
)

func main() {
	flag.Parse()

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := nollakonffi.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *nollakonffi.ServiceEntry)
	go func(results <-chan *nollakonffi.ServiceEntry) {
		for entry := range results {
			fmt.Printf("Found service: %s with IP %s and port %d\n", entry.ServiceName(), entry.AddrIPv4, entry.Port)
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*waitTime))
	defer cancel()
	err = resolver.Browse(ctx, *service, *domain, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}
