package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/jessevdk/go-flags"
)

// nolint:gochecknoglobals
var (
	appName        = "puing"
	appUsage       = "[OPTIONS] HOST"
	appDescription = "`ping` command but with puing"
	appVersion     = "???"
	appRevision    = "???"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrPing
)

type options struct {
	Version bool `short:"V" long:"version" description:"Show version"`
}

var imageCnt int

func main() {
	imageCnt = 0
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(
			color.Error,
			"[ %v ] %s\n",
			color.New(color.FgRed, color.Bold).Sprint("ERROR"),
			err,
		)
	}

	os.Exit(int(code))
}

func run(cliArgs []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage
	parser.ShortDescription = appDescription
	parser.LongDescription = appDescription

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}

		return exitCodeErrArgs, fmt.Errorf("parse error: %w", err)
	}

	if opts.Version {
		// nolint:forbidigo
		fmt.Printf("%s: v%s-rev%s\n", appName, appVersion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("must requires an argument")
	}

	if 1 < len(args) {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("too many arguments")
	}

	pinger, err := initPinger(args[0])
	if err != nil {
		return exitCodeOK, fmt.Errorf("an error occurred while initializing pinger: %w", err)
	}

	if err := pinger.Run(); err != nil {
		return exitCodeErrPing, fmt.Errorf("an error occurred when running ping: %w", err)
	}

	return exitCodeOK, nil
}

func initPinger(host string) (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to init pinger %w", err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		pinger.Stop()
	}()

	color.New(color.FgHiWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	pinger.OnRecv = pingerOnrecv
	pinger.OnFinish = pingerOnFinish

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	return pinger, nil
}

// nolint:forbidigo
func pingerOnrecv(pkt *ping.Packet) {
	fmt.Fprintf(
		color.Output,
		"%s seq=%s %sbytes from %s: ttl=%s time=%s\n",
		renderASCIIArt(imageCnt),
		color.New(color.FgHiYellow, color.Bold).Sprintf("%d", pkt.Seq),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%d", pkt.Ttl),
		color.New(color.FgHiMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
	)
	imageCnt++
}

// nolint:forbidigo
func pingerOnFinish(stats *ping.Statistics) {
	color.New(color.FgWhite, color.Bold).Fprintf(
		color.Output,
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Fprintf(
		color.Output,
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgHiRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Fprintf(
		color.Output,
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
}

func renderASCIIArt(idx int) string {
	/**
	if len(puing) <= idx {
		return strings.Repeat(" ", len(puing[0]))
	}
	**/
	if len(puing) <= idx {
		imageCnt = 0
	}

	line := puing[imageCnt]

	line = colorize(line, 'R', color.New(color.BgRed, color.Bold))
	line = colorize(line, 'M', color.New(color.BgHiMagenta, color.Bold))
	line = colorize(line, 'Y', color.New(color.BgHiYellow, color.Bold))
	line = colorize(line, 'I', color.New(color.BgYellow, color.Bold))
	line = colorize(line, 'A', color.New(color.BgBlue, color.Bold))
	line = colorize(line, 'C', color.New(color.BgCyan, color.Bold))
	line = colorize(line, 'B', color.New(color.BgHiBlack, color.Bold))
	line = colorize(line, 'W', color.New(color.BgHiWhite, color.Bold))
	line = lastline(line, '-', color.New(color.FgHiWhite, color.Bold))
	return line
}

func colorize(text string, target rune, color *color.Color) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint(" "),
	)
}

func lastline(text string, target rune, color *color.Color) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint("-"),
	)
}
