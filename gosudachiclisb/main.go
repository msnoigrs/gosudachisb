package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ichiban/assets"
	"github.com/msnoigrs/gosudachi"
	"github.com/msnoigrs/gosudachi/data"
	"github.com/msnoigrs/gosudachisb"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	plugins = make(map[string]reflect.Type)

	pkgstr                            = "github.com/msnoigrs/gosudachi"
	defaultInputTextPlugin            = pkgstr + ".DefaultInputTextPlugin"
	prolongedSoundMarkInputTextPlugin = pkgstr + ".ProlongedSoundMarkInputTextPlugin"
	inhibitConnectionPlugin           = pkgstr + ".InhibitConnectionPlugin"
	meCabOovProviderPlugin            = pkgstr + ".MeCabOovProviderPlugin"
	simpleOovProviderPlugin           = pkgstr + ".SimpleOovProviderPlugin"
	joinNumericPlugin                 = pkgstr + ".JoinNumericPlugin"
	joinKatakanaOovPlugin             = pkgstr + ".JoinKatakanaOovPlugin"
)

func init() {
	register(gosudachi.DefaultInputTextPlugin{})
	register(gosudachi.ProlongedSoundMarkInputTextPlugin{})
	register(gosudachi.MeCabOovProviderPlugin{})
	register(gosudachi.SimpleOovProviderPlugin{})
	register(gosudachi.JoinNumericPlugin{})
	register(gosudachi.JoinKatakanaOovPlugin{})
	register(gosudachi.InhibitConnectionPlugin{})
}

func register(x interface{}) {
	t := reflect.TypeOf(x)
	n := t.PkgPath() + "." + t.Name()
	plugins[n] = t
}

func newPlugin(name string) (interface{}, bool) {
	t, ok := plugins[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}

func makeInputTextPlugin(k string) gosudachi.InputTextPlugin {
	switch k {
	case "DefaultInputTextPlugin", "com.worksap.nlp.sudachi.DefaultInputTextPlugin", defaultInputTextPlugin:
		plugin, ok := newPlugin(defaultInputTextPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.InputTextPlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "ProlongedSoundMarkInputTextPlugin", "com.worksap.nlp.sudachi.ProlongedSoundMarkInputTextPlugin", prolongedSoundMarkInputTextPlugin:
		plugin, ok := newPlugin(prolongedSoundMarkInputTextPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.InputTextPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makeOovProviderPlugin(k string) gosudachi.OovProviderPlugin {
	switch k {
	case "MeCabOovProviderPlugin", "com.worksap.nlp.sudachi.MeCabOovProviderPlugin", meCabOovProviderPlugin:
		plugin, ok := newPlugin(meCabOovProviderPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.OovProviderPlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "SimpleOovProviderPlugin", "com.worksap.nlp.sudachi.SimpleOovProviderPlugin", simpleOovProviderPlugin:
		plugin, ok := newPlugin(simpleOovProviderPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.OovProviderPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makePathRewritePlugin(k string) gosudachi.PathRewritePlugin {
	switch k {
	case "JoinNumericPlugin", "com.worksap.nlp.sudachi.JoinNumericPlugin", joinNumericPlugin:
		plugin, ok := newPlugin(joinNumericPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.PathRewritePlugin)
		if !ok {
			return nil
		}
		return rplugin
	case "JoinKatakanaOovPlugin", "com.worksap.nlp.sudachi.JoinKatakanaOovPlugin", joinKatakanaOovPlugin:
		plugin, ok := newPlugin(joinKatakanaOovPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.PathRewritePlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func makeEditConnectionCostPlugin(k string) gosudachi.EditConnectionCostPlugin {
	switch k {
	case "InhibitConnectionPlugin", "com.worksap.nlp.sudachi.InhibitConnectionPlugin", inhibitConnectionPlugin:
		plugin, ok := newPlugin(inhibitConnectionPlugin)
		if !ok {
			return nil
		}
		rplugin, ok := plugin.(gosudachi.EditConnectionCostPlugin)
		if !ok {
			return nil
		}
		return rplugin
	}
	return nil
}

func runFromReader(tokenizer *gosudachi.JapaneseTokenizer, mode string, input io.Reader, output io.Writer, printAll bool, ignoreError bool) error {
	s := gosudachisb.NewLineScanner(input)
	for s.Scan() {
		err := run(tokenizer, mode, s.Text(), output, printAll)
		if err != nil {
			if ignoreError {
				fmt.Fprintln(os.Stderr, err)
			} else {
				return err
			}
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}

func run(tokenizer *gosudachi.JapaneseTokenizer, mode string, text string, output io.Writer, printAll bool) error {
	ms, err := tokenizer.Tokenize(mode, text)
	if err != nil {
		return err
	}
	for i := 0; i < ms.Length(); i++ {
		m := ms.Get(i)

		fmt.Fprintf(output, "%s\t%s\t%s",
			m.Surface(),
			strings.Join(m.PartOfSpeech(), ","),
			m.NormalizedForm())
		if printAll {
			fmt.Fprintf(output, "\t%s\t%s\t%d",
				m.DictionaryForm(),
				m.ReadingForm(),
				m.GetDictionaryId())
			if m.IsOOV() {
				fmt.Fprintf(output, "\t(OOV)")
			}
		}
		fmt.Fprintf(output, "\n")
	}
	fmt.Fprintln(output, "EOS")
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
	%s [-r file|-s jsonstring] [-m A|B|C] [-o file] [-p dir] [-j] [file ...]

Options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var (
		settingfile   string
		mergesettings string
		mode          string
		resourcesdir  string
		outputfile    string
		printall      bool
		ignoreerr     bool
		debugmode     bool
		utf16string   bool
	)
	flag.StringVar(&settingfile, "r", "", "read settings from file (overrides -s)")
	flag.StringVar(&mergesettings, "s", "", "additional settings (overrides -r)")
	flag.StringVar(&mode, "m", "C", "mode of splitting")
	flag.StringVar(&resourcesdir, "p", "", "root directory of resources")
	flag.StringVar(&outputfile, "o", "", "output to file")
	flag.BoolVar(&printall, "a", false, "print all fields")
	flag.BoolVar(&ignoreerr, "f", false, "ignore error")
	flag.BoolVar(&debugmode, "d", false, "debug mode")
	flag.BoolVar(&utf16string, "j", false, "use UTF-16 string")

	flag.Parse()

	if resourcesdir == "" {
		ex, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		resourcesdir = filepath.Dir(ex)
	}

	var output io.Writer
	if outputfile != "" {
		if !filepath.IsAbs(outputfile) {
			var err error
			outputfile, err = filepath.Abs(outputfile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		outputfd, err := os.OpenFile(outputfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", outputfile, err)
			os.Exit(1)
		}
		defer outputfd.Close()
		bufiooutput := bufio.NewWriter(outputfd)
		defer bufiooutput.Flush()
		output = bufiooutput
	} else {
		output = os.Stdout
	}

	l, err := assets.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\v", err)
		os.Exit(1)
	}
	defer l.Close()

	sysdicpath := filepath.Join(l.Path, "dict", "system.dic")

	settings, pluginmaker, err := parseSettings(resourcesdir, settingfile, mergesettings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to parse settings: %s\n", err)
		os.Exit(1)
	}

	settings.GetBaseConfig().SystemDict = sysdicpath

	if utf16string {
		settings.GetBaseConfig().Utf16String = utf16string
	}

	inputTextPlugins, err := pluginmaker.GetInputTextPluginArray(makeInputTextPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any InputTextPlugin: %s\n", err)
		os.Exit(1)
	}
	oovProviderPlugins, err := pluginmaker.GetOovProviderPluginArray(makeOovProviderPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any OovProviderPlugin: %s\n", err)
		os.Exit(1)
	}
	pathRewritePlugins, err := pluginmaker.GetPathRewritePluginArray(makePathRewritePlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any PathRewritePlugin: %s\n", err)
		os.Exit(1)
	}
	editConnectionCostPlugins, err := pluginmaker.GetEditConnectionCostPluginArray(makeEditConnectionCostPlugin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to cleate any ConnectionCostPlugin: %s\n", err)
		os.Exit(1)
	}

	dict, err := gosudachi.NewJapaneseDictionary(
		settings.GetBaseConfig(),
		inputTextPlugins,
		oovProviderPlugins,
		pathRewritePlugins,
		editConnectionCostPlugins,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dict.Close()

	tokenizer := dict.Create()
	if debugmode {
		tokenizer.DumpOutput = output
	}

	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			input, err := os.OpenFile(arg, os.O_RDONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s", arg, err)
				os.Exit(1)
			}
			err = runFromReader(tokenizer, mode, input, output, printall, ignoreerr)
			input.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	} else {
		err = runFromReader(tokenizer, mode, os.Stdin, output, printall, ignoreerr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func parseSettings(basePath string, settingfile string, mergeString string) (gosudachi.Settings, gosudachi.PluginMaker, error) {
	settings := gosudachi.NewSettingsJSON()

	var settingsreader io.Reader

	if settingfile != "" {
		var err error
		if !filepath.IsAbs(settingfile) {
			settingfile, err = filepath.Abs(settingfile)
			if err != nil {
				return nil, nil, err
			}
		}
		settingsfd, err := os.OpenFile(settingfile, os.O_RDONLY, 0644)
		if err != nil {
			return nil, nil, err
		}
		defer settingsfd.Close()
		settingsreader = settingsfd
	} else {
		settingsf, err := data.Assets.Open("sudachi.json")
		if err != nil {
			return nil, nil, err
		}
		defer settingsf.Close()
		settingsreader = settingsf
	}

	err := settings.ParseSettingsJSON(basePath, settingsreader)
	if err != nil {
		return nil, nil, err
	}

	if settingfile == "" && mergeString != "" {
		err = settings.ParseSettingsJSON(basePath, strings.NewReader(mergeString))
		if err != nil {
			return nil, nil, err
		}
	}
	return settings, settings, nil
}
