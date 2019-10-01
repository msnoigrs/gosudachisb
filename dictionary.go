package gosudachisb

import (
	"github.com/msnoigrs/gosudachi"
)

func NewInputTextPlugins() []gosudachi.InputTextPlugin {
	plugins := make([]gosudachi.InputTextPlugin, 2, 2)
	plugins[0] = gosudachi.NewDefaultInputTextPlugin(nil)
	repSymbol := "ー"
	plugins[1] = gosudachi.NewProlongedSoundMarkInputTextPlugin(
		&gosudachi.ProlongedSoundMarkInputTextPluginConfig{
			ProlongedSoundMarks: &[]string{"ー", "-", "⁓", "〜", "〰"},
			ReplacementSymbol:   &repSymbol,
		},
	)
	return plugins
}

func NewOovProviderPlugins() []gosudachi.OovProviderPlugin {
	plugins := make([]gosudachi.OovProviderPlugin, 2, 2)
	plugins[0] = gosudachi.NewMeCabOovProviderPlugin(nil)
	var (
		leftId  int16 = 5968
		rightId int16 = 5968
		cost    int16 = 3857
	)
	plugins[1] = gosudachi.NewSimpleOovProviderPlugin(
		&gosudachi.SimpleOovProviderPluginConfig{
			OovPos:  &[]string{"補助記号", "一般", "*", "*", "*", "*"},
			LeftId:  &leftId,
			RightId: &rightId,
			Cost:    &cost,
		},
	)
	return plugins
}

func NewPathRewritePlugins() []gosudachi.PathRewritePlugin {
	plugins := make([]gosudachi.PathRewritePlugin, 2, 2)
	enableNormalize := true
	plugins[0] = gosudachi.NewJoinNumericPlugin(
		&gosudachi.JoinNumericPluginConfig{
			EnableNormalize: &enableNormalize,
		},
	)
	minLength := 3
	plugins[1] = gosudachi.NewJoinKatakanaOovPlugin(
		&gosudachi.JoinKatakanaOovPluginConfig{
			OovPOS:    &[]string{"名詞", "普通名詞", "一般", "*", "*", "*"},
			MinLength: &minLength,
		},
	)
	return plugins
}

func NewDictionary(config *gosudachi.BaseConfig) (*gosudachi.JapaneseDictionary, error) {
	inputTextPlugins := NewInputTextPlugins()
	oovProviderPlugins := NewOovProviderPlugins()
	pathRewritePlugins := NewPathRewritePlugins()

	return gosudachi.NewJapaneseDictionary(
		config,
		inputTextPlugins,
		oovProviderPlugins,
		pathRewritePlugins,
		nil,
	)
}
