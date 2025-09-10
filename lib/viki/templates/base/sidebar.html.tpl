<ul>
{{- range .Nodes }}
    <li>
        {{- if .IsDir }}
            <span>{{ .Name }}</span>
            {{- if .Children }}
                {{ template "inner" .Children }}
            {{- end }}
        {{- else }}
            <a href="{{ .URL }}">{{ .Name }}</a>
        {{- end }}
    </li>
{{- end }}
</ul>

{{/* Define the recursive template */}}
{{ define "inner" }}
<ul>
{{- range . }}
    <li>
        {{- if .IsDir }}
            <span>{{ .Name }}</span>
            {{- if .Children }}
                {{ template "inner" .Children }}
            {{- end }}
        {{- else }}
            <a href="/{{ .URL }}">{{ .Name }}</a>
        {{- end }}
    </li>
{{- end }}
</ul>
{{ end }}