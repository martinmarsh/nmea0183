package nmea0183

import (
	"fmt"

	"github.com/spf13/viper"
)


type Sentences struct {
	formats, variables map[string][]string
}


func MakeSentences(formats, variables map[string][]string) *Sentences{
	sent := Sentences{formats: formats, variables: variables}
	return &sent
}

func (sent *Sentences) MakeHandle() *Handle {
	h := setUp()
	h.sentences = sent
	return h
}

func (sent *Sentences) AddFormat(key string, form []string){
	sent.formats[key] = form
}

func (sent *Sentences) AddVariable(key string, varFormat []string){
	sent.variables[key] = varFormat
}

func (sent *Sentences) Load(setting ...string) error {
	configSet := []string{".", "nmea_sentences", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = fmt.Errorf("sentence file was not found. Use Create or download nmea_config.yaml: %w", err)
			return err
		} else {
			// Handle file was found but another error was produced
			err = fmt.Errorf("fatal error in config file: %w", err)
			return err
		}
	}

	sent.formats = viper.GetStringMapStringSlice("formats")
	sent.variables = viper.GetStringMapStringSlice("variables")

	return err
}

func (sent *Sentences) SaveDefault(setting ...string){
	configSet := []string{".", "nmea_sentences", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	viper.SetDefault("formats", GetDefaultFormats())
	viper.SetDefault("variables", GetDefaultVars())
	err := viper.ReadInConfig() // Find and read the config file

    viper.GetViper()
	//Don't overwrite if file exists
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		viper.SafeWriteConfig()
	}
}
