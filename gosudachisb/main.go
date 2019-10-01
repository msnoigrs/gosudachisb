package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ichiban/assets"
	"github.com/msnoigrs/gosudachi"
	"github.com/msnoigrs/gosudachisb"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
	%s [-m A|B|C] [-o file] [file ...]

Options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var (
		mode       string
		outputfile string
		printall   bool
		ignoreerr  bool
		debugmode  bool
	)
	flag.StringVar(&mode, "m", "C", "mode of splitting")
	flag.StringVar(&outputfile, "o", "", "output to file")
	flag.BoolVar(&printall, "a", false, "print all fields")
	flag.BoolVar(&ignoreerr, "f", false, "ignore error")
	flag.BoolVar(&debugmode, "d", false, "debug mode")

	flag.Parse()

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

	baseConfig := gosudachi.BaseConfig{
		SystemDict: sysdicpath,
	}

	dict, err := gosudachisb.NewDictionary(&baseConfig)
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
