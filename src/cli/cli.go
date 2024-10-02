package cli

import "github.com/spf13/cobra"

type Flags struct {
	Url          string
	NoRecursion  bool
	NoSubdomains bool
	NoExternals  bool
	WithAssets   bool
	Threads      uint64
	Output       string
	Quiet        bool
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

	cmd.Flags().BoolVarP(&flags.Quiet, "quiet", "q", false, "Disable printing scraped URLs on the screen")

	cmd.Flags().Uint64VarP(&flags.Threads, "threads", "t", 10, "Number of concurrent threads")

	cmd.Flags().BoolVar(&flags.NoRecursion, "no-recursion", false, "Disable recursive scraping")

	cmd.Flags().BoolVar(&flags.NoSubdomains, "no-subdomains", false, "Disable subdomains scraping")

	cmd.Flags().BoolVar(&flags.NoExternals, "no-externals", false, "Disable external URLs scraping")

	cmd.Flags().BoolVar(&flags.WithAssets, "with-assets", false, "Enable asset URLs scraping (images, CSS, JS etc.)")

	cmd.Execute()
}
