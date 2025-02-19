package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/go-ping/ping"
	"github.com/gookit/color"
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
	Count     int  `short:"c" long:"count" default:"20" description:"Stop after <count> replies"`
	Privilege bool `short:"P" long:"privilege" description:"Enable privileged mode"`
	Version   bool `short:"V" long:"version" description:"Show version"`
}

var imageCnt int

func main() {
	imageCnt = 0
	code, err := run(os.Args[1:])
	if err != nil {
		color.Error.Println("ERROR")
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

	pinger, err := initPinger(args[0], opts)
	if err != nil {
		return exitCodeOK, fmt.Errorf("an error occurred while initializing pinger: %w", err)
	}

	if err := pinger.Run(); err != nil {
		return exitCodeErrPing, fmt.Errorf("an error occurred when running ping: %w", err)
	}

	return exitCodeOK, nil
}

func initPinger(host string, opts options) (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to init pinger %w", err)
	}

	pinger.Count = opts.Count

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		pinger.Stop()
	}()

	color.New(color.FgWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	pinger.OnRecv = pingerOnrecv
	pinger.OnFinish = pingerOnFinish

	if opts.Privilege || runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	return pinger, nil
}

func pingerOnrecv(pkt *ping.Packet) {
	fmt.Printf(
		"%s seq=%s %sbytes from %s: ttl=%s time=%s\n",
		renderASCIIArt(imageCnt),
		color.New(color.FgYellow, color.Bold).Sprintf("%d", pkt.Seq),
		color.New(color.FgBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
		color.New(color.FgCyan, color.Bold).Sprintf("%d", pkt.Ttl),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
	)
	imageCnt++
}

func pingerOnFinish(stats *ping.Statistics) {
	color.New(color.FgWhite, color.Bold).Printf(
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Printf(
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Printf(
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
}

func renderASCIIArt(idx int) string {

	if len(puing) <= idx {
		imageCnt = 0
	}

	line := puing[imageCnt]

	//Base colors
	line = colorize(line, 'B', color.New(color.BgBlack, color.Bold))
	line = colorize(line, 'W', color.New(color.BgHiWhite, color.Bold))
	line = colorize(line, 'C', color.New(color.BgHiCyan, color.Bold))
	line = colorize(line, 'A', color.New(color.BgCyan, color.Bold))

	//RGB colors
	line = colorizeRGB(line, 'Y', color.NewRGBStyle(color.RGB(255, 255, 105), color.RGB(255, 255, 105))) // hair color
	line = colorizeRGB(line, 'K', color.NewRGBStyle(color.RGB(207, 197, 95), color.RGB(207, 197, 95)))   // Shadow 1 of Yellow
	line = colorizeRGB(line, 'M', color.NewRGBStyle(color.RGB(241, 182, 178), color.RGB(241, 182, 178))) // Mouse color
	line = colorizeRGB(line, 'N', color.NewRGBStyle(color.RGB(192, 120, 121), color.RGB(192, 120, 121))) // Shadow 1 of Mouse color
	line = colorizeRGB(line, 'S', color.NewRGBStyle(color.RGB(255, 242, 236), color.RGB(255, 242, 236))) // Skin color

	line = lastline(line, '-', color.New(color.FgWhite, color.Bold))
	return line
}

func colorize(text string, target rune, color color.Style) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint(" "),
	)
}
func colorizeRGB(text string, target rune, color *color.RGBStyle) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint(" "),
	)
}
func lastline(text string, target rune, color color.Style) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint("-"),
	)
}
