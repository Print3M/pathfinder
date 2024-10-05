package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Flags struct {
	Url          string
	NoRecursion  bool
	NoSubdomains bool
	NoExternals  bool
	WithAssets   bool
	Headers      []string
	Threads      uint
	Output       string
	Quiet        bool
	Rate         uint
}

func InitCli(main func(flags *Flags)) {
	var flags Flags

	var cmd = &cobra.Command{
		Use:   "PathFinder",
		Short: "PathFinder is a crawler script for concurrent and recursive scraping of URLs from any website.",
		Long:  "PathFinder is a crawler script for concurrent and recursive scraping of URLs from any website.",
		Run: func(cmd *cobra.Command, args []string) {
			main(&flags)
		},
	}

	cmd.Flags().StringVarP(&flags.Url, "url", "u", "", "URL to start")
	cmd.MarkFlagRequired("url")

	cmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Output file")

	cmd.Flags().StringArrayVarP(&flags.Headers, "headers", "H", []string{}, "Add HTTP header (one -H must contain only one header)")

	cmd.Flags().BoolVarP(&flags.Quiet, "quiet", "q", false, "Disable printing scraped URLs on the screen")

	cmd.Flags().UintVarP(&flags.Threads, "threads", "t", 10, "Number of concurrent threads")

	cmd.Flags().UintVarP(&flags.Rate, "rate", "r", 0, "Number of requests per second")

	cmd.Flags().BoolVar(&flags.NoRecursion, "no-recursion", false, "Disable recursive scraping")

	cmd.Flags().BoolVar(&flags.NoSubdomains, "no-subdomains", false, "Disable subdomains scraping")

	cmd.Flags().BoolVar(&flags.NoExternals, "no-externals", false, "Disable external URLs scraping")

	cmd.Flags().BoolVar(&flags.WithAssets, "with-assets", false, "Enable asset URLs scraping (images, CSS, JS etc.)")

	cmd.Execute()
}

func PrepareOutputFile(name string) *os.File {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Open file error: %v\n", err)
		os.Exit(1)
	}

	return file
}
