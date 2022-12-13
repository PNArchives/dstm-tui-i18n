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

	"github.com/PNCommand/dstm/dst/env"
	l10n "github.com/PNCommand/dstm/localization"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "dstm",
	Version: "v0.0.1",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		str := l10n.Singleton().String("_long_des")
		println(str)

		env.CheckSystem()
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

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate("(*•ᴗ•*) " + rootCmd.Use + " " + rootCmd.Version + "\n")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.dstm.toml)")
	rootCmd.Flags().StringP("lang", "l", "en", "specify language")

	viper.BindPFlag("lang", rootCmd.Flags().Lookup("lang"))
}

func initDefaultConfig() {
	viper.SetDefault("lang", "en")

	viper.SetDefault("separator", "-")

	viper.SetDefault("dstRootDir", os.ExpandEnv("$HOME/Server"))
	viper.SetDefault("ugcDir", os.ExpandEnv("$HOME/Server/ugc_mods"))
	viper.SetDefault("v1ModDir", os.ExpandEnv("$HOME/Server/mods"))
	viper.SetDefault("v2ModDir", os.ExpandEnv("$HOME/Server/ugc_mods/content"))

	viper.SetDefault("kleiRootDir", os.ExpandEnv("$HOME/Klei"))
	viper.SetDefault("worldsDirName", "worlds")
}

// initConfig reads in config file and ENV variables if set.
// 引数 > 環境変数 > 設定ファイル
func initConfig() {
	initDefaultConfig()

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigName(".dstm")
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix("dstm")
	viper.BindEnv("lang")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			err := viper.SafeWriteConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			// Config file was found but another error was produced
			fmt.Println("Config file was found but another error was produced")
			os.Exit(1)
		}
	}
	l10n.SetLocale(viper.GetString("lang"))
}
