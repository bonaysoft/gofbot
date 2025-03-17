package cmd

/*
Copyright © 2025 Ambor <saltbo@foxmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bonaysoft/gofbot/pkg/adapters"
	"github.com/bonaysoft/gofbot/pkg/messenger"
	"github.com/bonaysoft/gofbot/pkg/storage"
)

var debug bool

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:          "template [TEMPLATES DIRECTORY]",
	Short:        "Render bot templates locally and display the output.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("storage-file-location", args[0])

		values, err := loadValues()
		if err != nil {
			return err
		}

		adapter, err := adapters.GetAdapter(viper.GetString("adapter"))
		if err != nil {
			return err
		}

		store, err := storage.New("file")
		if err != nil {
			return err
		}
		if err := store.Start(cmd.Context()); err != nil {
			return err
		}

		mm := messenger.NewDefaultManager(store, adapter.GetFunMap())
		values["chatProvider"] = adapter.Name()
		slog.Debug("matching", slog.Any("values", values))
		msg, err := mm.Match(values)
		if err != nil {
			return err
		}

		out, err := mm.BuildReply(msg, values)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

// loadValues loads values from a file.
// If it's the default file, it's allowed to not exist, otherwise return an error.
// Supports multiple file formats, including YAML, JSON, and TOML.
func loadValues() (map[string]any, error) {
	valuesFile := viper.GetString("values")
	values := make(map[string]any)

	// Non-default file must exist
	if _, err := os.Stat(valuesFile); os.IsNotExist(err) && !strings.HasPrefix(valuesFile, "values.") {
		return nil, fmt.Errorf("values file %q not found", valuesFile)
	}

	// Read and parse the file
	v := viper.New()
	v.SetConfigFile(valuesFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read values file: %w", err)
	}

	if err := v.Unmarshal(&values); err != nil {
		return nil, fmt.Errorf("failed to unmarshal values: %w", err)
	}

	// Handle command line value overrides with --set
	sets := viper.GetStringSlice("set")
	for _, s := range sets {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid set value: %s, must be in format key=value", s)
		}
		values[parts[0]] = parts[1]
	}

	return values, nil
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().String("adapter", "terminal", "specify the adapter name")
	templateCmd.Flags().StringArray("set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	templateCmd.Flags().StringP("values", "f", "values.json", "specify values in a YAML file or a URL (can specify multiple)")
	templateCmd.Flags().BoolVar(&debug, "debug", debug, "enable verbose output")

	_ = viper.BindPFlags(templateCmd.Flags())
}
