package routes

import (
	"net"

	"github.com/pkg/errors"
	"go.universe.tf/metallb/e2etest/pkg/executor"
)

func Add(exec executor.Executor, target, via string) error {
	_, dst, err := net.ParseCIDR(target)
	if err != nil {
		return err
	}

	gw := net.ParseIP(via)

	cmd := "ip"
	args := []string{"route", "add", dst.String(), "via", gw.String()}
	if dst.IP.To4() == nil {
		args = []string{"-6", "route", "add", dst.String(), "via", gw.String()}
	}
	out, err := exec.Exec(cmd, args...)
	if err != nil {
		return errors.Wrapf(err, "Failed to add route %s %s %s", cmd, args, out)
	}

	return nil
}

func Delete(exec executor.Executor, target, via string) error {
	_, dst, err := net.ParseCIDR(target)
	if err != nil {
		return err
	}

	gw := net.ParseIP(via)

	cmd := "ip"
	args := []string{"route", "del", dst.String(), "via", gw.String()}
	if dst.IP.To4() == nil {
		args = []string{"-6", "route", "del", dst.String(), "via", gw.String()}
	}
	out, err := exec.Exec(cmd, args...)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete route %s", out)
	}

	return nil
}
