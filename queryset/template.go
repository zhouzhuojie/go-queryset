package queryset

import (
	"text/template"

	"github.com/jirfag/go-queryset/queryset/methods"
)

var qsTmpl = template.Must(
	template.New("generator").
		Funcs(template.FuncMap{
			"lcf": methods.LowercaseFirstRune,
		}).
		Parse(qsCode),
)

const qsCode = `
// notest

// ===== BEGIN of all query sets

{{ range .Configs }}
  // ===== BEGIN of query set {{ .Name }}

	// {{ .Name }} is an queryset type for {{ .StructName }}
  type {{ .Name }} struct {
	  db *gorm.DB
  }

  // New{{ .Name }} constructs new {{ .Name }}
  func New{{ .Name }}(db *gorm.DB) {{ .Name }} {
	  return {{ .Name }}{
		  db: db.Model(&{{ .StructName }}{}),
	  }
  }

	func (qs {{ .Name }}) w(db *gorm.DB) {{ .Name }} {
	  return New{{ .Name }}(db)
  }

	{{ range .Methods }}
		{{ .GetDoc .GetMethodName }}
		func ({{ .GetReceiverDeclaration }}) {{ .GetMethodName }}({{ .GetArgsDeclaration }})
		{{- .GetReturnValuesDeclaration }} {
      {{ .GetBody }}
		}
	{{ end }}

  // ===== END of query set {{ .Name }}

	// ===== BEGIN of {{ .StructName }} modifiers

	{{ $ft := printf "%s%s" .StructName "DBSchemaField" | lcf }}
	type {{ $ft }} string

	func (f {{ $ft }}) String() string {
		return string(f)
	}

	// {{ .StructName }}DBSchema stores db field names of {{ .StructName }}
	var {{ .StructName }}DBSchema = struct {
		{{ range .Fields }}
			{{ .Name }} {{ $ft }}
		{{- end }}
	}{
		{{ range .Fields }}
			{{ .Name }}: {{ $ft }}("{{ .DBName }}"),
		{{- end }}
	}

	// Update updates {{ .StructName }} fields by primary key
	func (o *{{ .StructName }}) Update(db *gorm.DB, fields ...{{ $ft }}) error {
		dbNameToFieldName := map[string]interface{}{
			{{- range .Fields }}
				"{{ .DBName }}": o.{{ .Name }},
			{{- end }}
		}
		u := map[string]interface{}{}
		for _, f := range fields {
			fs := f.String()
			u[fs] = dbNameToFieldName[fs]
		}
		if err := db.Model(o).Updates(u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return err
			}

			return fmt.Errorf("can't update {{ .StructName }} %v fields %v: %s",
				o, fields, err)
		}

		return nil
	}

	// {{ .StructName }}Updater is an {{ .StructName }} updates manager
	type {{ .StructName }}Updater struct {
		fields map[string]interface{}
		db *gorm.DB
	}

	// New{{ .StructName }}Updater creates new {{ .StructName }} updater
	func New{{ .StructName }}Updater(db *gorm.DB) {{ .StructName }}Updater {
		return {{ .StructName }}Updater{
			fields: map[string]interface{}{},
			db: db.Model(&{{ .StructName }}{}),
		}
	}

	// ===== END of {{ .StructName }} modifiers
{{ end }}

// ===== END of all query sets
`
