package src

import (
	"testing"
)

func TestLanguageNext(t *testing.T) {
	lang := LangEN
	if lang.Next() != LangES {
		t.Fatalf("expected LangES, got %s", lang.Next())
	}
	if LangES.Next() != LangEN {
		t.Fatalf("expected LangEN, got %s", LangES.Next())
	}
}

func TestDictionariesCompleteness(t *testing.T) {
	languages := []Language{LangEN, LangES}

	for _, lang := range languages {
		d := T(lang)

		if d.AboutTitle == "" {
			t.Errorf("[%s] AboutTitle is empty", lang)
		}
		if d.AboutSubtitle == "" {
			t.Errorf("[%s] AboutSubtitle is empty", lang)
		}
		if d.AboutVersion == "" {
			t.Errorf("[%s] AboutVersion is empty", lang)
		}
		if d.AboutCreator == "" {
			t.Errorf("[%s] AboutCreator is empty", lang)
		}
		if d.AboutRepo == "" {
			t.Errorf("[%s] AboutRepo is empty", lang)
		}
		if d.AboutLicense == "" {
			t.Errorf("[%s] AboutLicense is empty", lang)
		}
		if d.AboutLicenseVal == "" {
			t.Errorf("[%s] AboutLicenseVal is empty", lang)
		}
		if d.AboutYear == "" {
			t.Errorf("[%s] AboutYear is empty", lang)
		}
		if d.AboutBuiltWith == "" {
			t.Errorf("[%s] AboutBuiltWith is empty", lang)
		}
		if d.AboutReturnHint == "" {
			t.Errorf("[%s] AboutReturnHint is empty", lang)
		}
	}
}

func TestSpanishAboutContent(t *testing.T) {
	d := T(LangES)
	if d.AboutReturnHint != "Pulsa cualquier tecla para volver" {
		t.Errorf("expected 'Pulsa cualquier tecla para volver', got %q", d.AboutReturnHint)
	}
	if d.AboutVersion != "Versión:" {
		t.Errorf("expected 'Versión:', got %q", d.AboutVersion)
	}
}
