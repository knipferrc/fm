package config

import (
	"log"
	"os"

	"github.com/knipferrc/fm/internal/constants"

	"github.com/spf13/viper"
)

type SettingsConfig struct {
	StartDir         string `mapstructure:"start_dir"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	RoundedPanes     bool   `mapstructure:"rounded_panes"`
	EnableLogging    bool   `mapstructure:"enable_logging"`
	EnableMouseWheel bool   `mapstructure:"enable_mousewheel"`
	PrettyMarkdown   bool   `mapstructure:"pretty_markdown"`
}

type DirTreeColors struct {
	SelectedItem   string `mapstructure:"selected_item"`
	UnselectedItem string `mapstructure:"unselected_item"`
}

type PaneColors struct {
	ActiveBorderColor   string `mapstructure:"active_border_color"`
	InactiveBorderColor string `mapstructure:"inactive_border_color"`
}

type ColorVariant struct {
	Foreground string `mapstructure:"foreground"`
	Background string `mapstructure:"background"`
}

type StatusBarColors struct {
	SelectedFile ColorVariant `mapstructure:"selected_file"`
	Bar          ColorVariant `mapstructure:"bar"`
	TotalFiles   ColorVariant `mapstructure:"total_files"`
	Logo         ColorVariant `mapstructure:"logo"`
}

type ColorsConfig struct {
	DirTree   DirTreeColors   `mapstructure:"dir_tree"`
	Pane      PaneColors      `mapstructure:"pane"`
	Spinner   string          `mapstructure:"spinner"`
	StatusBar StatusBarColors `mapstructure:"status_bar"`
}

// Main app config
type Config struct {
	Settings SettingsConfig `mapstructure:"settings"`
	Colors   ColorsConfig   `mapstructure:"colors"`
}

// Load users config and create the config if it does not exist
func LoadConfig() {
	// Create an fm directory in the .config dir in the users home
	// directory if it does not already exist
	if _, err := os.Stat(os.ExpandEnv("$HOME/.config/fm")); os.IsNotExist(err) {
		err := os.Mkdir(os.ExpandEnv("$HOME/.config/fm"), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Set config file inside .config folder under an fm folder
	viper.SetConfigFile(os.ExpandEnv("$HOME/.config/fm/config.yml"))

	// Set viper config name, type and path
	viper.SafeWriteConfigAs(os.ExpandEnv("$HOME/.config/fm/config.yml"))

	// Read the current config file in and catch an error if it occurs
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// Get the users config and return it for use
func GetConfig() (config Config) {
	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatal("Error parsing config", err)
	}

	return
}

// Setup viper defaults
func SetDefaults() {
	// App Settings
	viper.SetDefault("settings.start_dir", ".")
	viper.SetDefault("settings.show_icons", true)
	viper.SetDefault("settings.rounded_panes", false)
	viper.SetDefault("settings.enable_logging", false)
	viper.SetDefault("settings.enable_mousewheel", true)
	viper.SetDefault("settings.pretty_markdown", true)

	// DirTree colors
	viper.SetDefault("colors.dir_tree.selected_item", constants.Pink)
	viper.SetDefault("colors.dir_tree.unselected_item", constants.White)

	// Pane colors
	viper.SetDefault("colors.pane.active_border_color", constants.Pink)
	viper.SetDefault("colors.pane.inactive_border_color", constants.White)

	// Spinner colors
	viper.SetDefault("colors.spinner", constants.Pink)

	// StatusBar colors
	viper.SetDefault("colors.status_bar.selected_file.foreground", constants.White)
	viper.SetDefault("colors.status_bar.selected_file.background", constants.Pink)
	viper.SetDefault("colors.status_bar.bar.foreground", constants.White)
	viper.SetDefault("colors.status_bar.bar.background", constants.DarkGray)
	viper.SetDefault("colors.status_bar.total_files.foreground", constants.White)
	viper.SetDefault("colors.status_bar.total_files.background", constants.LightPurple)
	viper.SetDefault("colors.status_bar.logo.foreground", constants.White)
	viper.SetDefault("colors.status_bar.logo.background", constants.DarkPurple)
}