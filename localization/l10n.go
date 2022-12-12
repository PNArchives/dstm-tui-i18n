package localization

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

type localization struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

var (
	locale    = language.Japanese
	singleton = newLocalizer()
	//go:embed en
	//go:embed zh
	//go:embed ja
	localizedFS embed.FS
)

func Singleton() *localization {
	return &singleton
}

func newLocalizer() localization {
	// en, zh, ja
	b := i18n.NewBundle(locale)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// load all toml files in the localization directory
	err := fs.WalkDir(localizedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".toml") {
			return nil
		}
		bytes, err := localizedFS.ReadFile(path)
		if err != nil {
			return err
		}
		b.ParseMessageFileBytes(bytes, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	l := i18n.NewLocalizer(b, locale.String())
	return localization{
		bundle:    b,
		localizer: l,
	}
}

func SetLocale(str string) {
	var matcher = language.NewMatcher([]language.Tag{
		language.English,
		language.Chinese,
		language.Japanese,
	})
	tag, _ := language.MatchStrings(matcher, str)
	if locale == tag {
		return
	}
	locale = tag
	singleton.localizer = i18n.NewLocalizer(singleton.bundle, locale.String())
}

func (l localization) String(key string) string {
	return l.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
}

func (l localization) String4Plural(key string, pluralCount int) string {
	return l.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:   key,
		PluralCount: pluralCount,
	})
}

func (l localization) String4Data(key string, data map[string]interface{}) string {
	return l.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
}

func (l localization) String4PluralData(key string, pluralCount int, data map[string]interface{}) string {
	return l.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		PluralCount:  pluralCount,
		TemplateData: data,
	})
}
