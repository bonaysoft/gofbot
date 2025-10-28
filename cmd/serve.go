package cmd

/*
Copyright © 2024 Ambor <saltbo@foxmail.com>

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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bonaysoft/gofbot/pkg/adapters"
	"github.com/bonaysoft/gofbot/pkg/bot"
	"github.com/bonaysoft/gofbot/pkg/messenger"
	"github.com/bonaysoft/gofbot/pkg/storage"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, err := adapters.GetAdapter(viper.GetString("adapter"))
		if err != nil {
			return err
		}

		store, err := storage.New(viper.GetString("storage"))
		if err != nil {
			return err
		}

		if err := store.Start(cmd.Context()); err != nil {
			return err
		}

		s, err := bot.NewServer(adapter, messenger.NewDefaultManager(store, adapter.GetFunMap()))
		if err != nil {
			return err
		}
		return s.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().Int("port", 9613, "specify the port of the webhook server")

	serveCmd.PersistentFlags().String("webhook-scheme", "http", "specify the scheme of the webhook URL")
	serveCmd.PersistentFlags().String("webhook-host", "", "specify the host of the webhook URL")

	serveCmd.PersistentFlags().String("storage", "file", "specify the storage name")
	serveCmd.PersistentFlags().String("storage-file-location", "data/templates", "specify the file storage location")

	_ = viper.BindPFlags(serveCmd.PersistentFlags())
}
