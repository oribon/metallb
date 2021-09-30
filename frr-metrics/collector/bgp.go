package collector

import (
	"os/exec"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	bgpstats "go.universe.tf/metallb/internal/bgp"
	"go.universe.tf/metallb/internal/bgp/bgpfrr"
)

var (
	sessionUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(bgpstats.Namespace, bgpstats.Subsystem, bgpstats.SessionUp.Name),
		bgpstats.SessionUp.Help,
		bgpstats.Labels,
		nil,
	)

	updatesSentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(bgpstats.Namespace, bgpstats.Subsystem, bgpstats.UpdatesSent.Name),
		bgpstats.UpdatesSent.Help,
		bgpstats.Labels,
		nil,
	)

	prefixesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(bgpstats.Namespace, bgpstats.Subsystem, bgpstats.Prefixes.Name),
		bgpstats.Prefixes.Help,
		bgpstats.Labels,
		nil,
	)
)

type bgp struct {
	Log          log.Logger
	getNeighbors func() (string, error)
}

func NewBGP(l log.Logger) *bgp {
	log := log.With(l, "collector", bgpstats.Subsystem)
	return &bgp{Log: log, getNeighbors: getBGPNeighbors}
}

func (c *bgp) Describe(ch chan<- *prometheus.Desc) {
	ch <- sessionUpDesc
	ch <- updatesSentDesc
	ch <- prefixesDesc
}

func (c *bgp) Collect(ch chan<- prometheus.Metric) {
	res, err := c.getNeighbors()
	if err != nil {
		level.Error(c.Log).Log("error", err, "msg", "failed to fetch BGP neighbors from FRR")
		return
	}

	neighbors, err := bgpfrr.ParseNeighbours(res)
	if err != nil {
		level.Error(c.Log).Log("error", err, "msg", "failed to parse BGP neighbors from output")
		return
	}

	updateNeighborsMetrics(ch, neighbors)
}

func updateNeighborsMetrics(ch chan<- prometheus.Metric, neighbors []*bgpfrr.Neighbor) {
	for _, n := range neighbors {
		sessionUp := 1
		if !n.Connected {
			sessionUp = 0
		}
		ip := n.Ip.String()

		ch <- prometheus.MustNewConstMetric(sessionUpDesc, prometheus.GaugeValue, float64(sessionUp), ip)
		ch <- prometheus.MustNewConstMetric(updatesSentDesc, prometheus.CounterValue, float64(n.UpdatesSent), ip)
		ch <- prometheus.MustNewConstMetric(prefixesDesc, prometheus.GaugeValue, float64(n.PfxSnt), ip)
	}
}

func getBGPNeighbors() (string, error) {
	res, err := runVtysh("show bgp neighbors json")
	if err != nil {
		return "", err
	}

	return res, err
}

func runVtysh(args ...string) (string, error) {
	newArgs := append([]string{"-c"}, args...)
	out, err := exec.Command("/usr/bin/vtysh", newArgs...).CombinedOutput()
	return string(out), err
}
