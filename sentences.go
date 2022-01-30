package nmea0183

import (
	"fmt"

	"github.com/spf13/viper"
)


// Methods on Sentences are used to define sentence parsing definitions and create handlers
type Sentences struct {
	formats map[string][]string
	variables map[string]string
}

// Pass a map containing a list of variable names for each sentence definition
// and a map associating the variable with an internal format definition string
// returns a sentence pointer to the made structure
func MakeSentences(formats  map[string][]string, variables map[string]string) *Sentences{
	sent := Sentences{formats: formats, variables: variables}
	return &sent
}

// This method on a sentence definition created a "Handle" struct and returns a pointer
// This pointer is used as a handle to the definition and data and methods such as Parse
// An advanced application might have many "Handlers" to deal with different data sources
func (sent *Sentences) MakeHandle() *Handle {
	h := setUp()
	h.sentences = sent
	return h
}

// Adds a sentence format to sentences.  Give the sentence name eg RMC followed by a
// list of variable names you wish to use.
// This would be used instead of an external definitions file
func (sent *Sentences) AddFormat(key string, form []string){
	sent.formats[key] = form
}

// Defines how each variable will be parsed from or written to sentences
// give the sentence name folowed by the internal format definition
// This would be used instead of an external definitions file
func (sent *Sentences) AddVariable(key string, varFormat string){
	sent.variables[key] = varFormat
}

// Loads sentence definitions from a file
// If no parameters uses defaults.
// 1st parameter is the path followed by the file name and format
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
	sent.variables = viper.GetStringMapString("variables")

	return err
}

// Loads a default definitions if the definition files does not exist
// and then writes the file.
// This is intended to help write definition files by producing a copy based on 
// in default examples. Edit it or use as a template for other definition files
func (sent *Sentences) SaveLoadDefault(setting ...string){
	configSet := []string{".", "nmea_sentences", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	viper.SetDefault("formats", GetDefaultFormats())
	viper.SetDefault("variables", GetDefaultVars())
	viper.ReadInConfig() // Find and read the config file

	//Don't overwrite if file exists
	//if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	viper.SafeWriteConfig()
	//}
	sent.formats = viper.GetStringMapStringSlice("formats")
	sent.variables = viper.GetStringMapString("variables")

}
