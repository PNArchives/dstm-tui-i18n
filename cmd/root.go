/*
Copyright © 2022 yechentide

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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "dstm",
	Version: "v0.0.1",
	Short:   "A brief description of your application",
	Long:    "A longer description.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello DSTM!")
		fmt.Println("Language:", viper.GetString("lang"))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// cobra.Command 実行前の初期化処理を定義する
	// rootCmd.Execute > コマンドライン引数の処理 > cobra.OnInitialize > rootCmd.Run という順に実行される
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("(*•ᴗ•*) " + rootCmd.Use + " " + rootCmd.Version + "\n")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.dstm.toml)")
	rootCmd.Flags().StringP("lang", "l", "en", "specify language")

	viper.BindPFlag("lang", rootCmd.Flags().Lookup("lang"))
	viper.SetDefault("lang", "en")

	// l10n.Locale = matchLangTag(appConf.Common.Lang)
	// rootCmd.Short = local.String("_short_des", l10n.MsgOnly, 0, nil)
	// rootCmd.Long = local.String("_long_des", l10n.MsgOnly, 0, nil)
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
// 引数 > 設定ファイル > 環境変数
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".dstm")
		viper.SetConfigType("toml")
	}

	viper.SetEnvPrefix("dstm")
	viper.BindEnv("lang")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("========== ========== ========== ========== ==========")
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		for k, v := range viper.AllSettings() {
			fmt.Printf("    %s: %s\n", k, v)
		}
		fmt.Println("========== ========== ========== ========== ==========")
	}
}
