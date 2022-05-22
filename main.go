package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/wailorman/tplinkexporter/clients"
	"github.com/wailorman/tplinkexporter/collectors"
)

func main() {
	var (
		host      = kingpin.Flag("host", "Host of target tplink easysmart switch.").Required().String()
		username  = kingpin.Flag("username", "Username for switch GUI login").Default("admin").String()
		password  = kingpin.Flag("password", "Password for switch GUI login").Required().String()
		hostname  = kingpin.Flag("hostname", "Name of the switch").String()
		portnames = kingpin.Flag("portnames", "Port descriptions, Format: 1-Server,2-Lab").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	portNamesMap := make(map[int]string)
	if *portnames != "" {
		for i, pair := range strings.Split(*portnames, ",") {
			pairData := strings.SplitN(pair, "-", 2)
			portNum, err := strconv.Atoi(pairData[0])

			if err != nil {
				log.Fatalf("Failed to parse port names. Pair #%d: %s", i, err)
			}

			portNamesMap[portNum] = pairData[1]
		}
	}

	tplinkSwitch := clients.NewTPLinkSwitch(*host, *hostname, *username, *password, portNamesMap)
	trafficCollector := collectors.NewTrafficCollector("tplink_exporter", tplinkSwitch)
	prometheus.MustRegister(trafficCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Beginning to serve on port :9717")
	log.Fatal(http.ListenAndServe(":9717", nil))
}
