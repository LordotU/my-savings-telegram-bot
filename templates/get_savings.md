Your total savings is *{{.TotalInUserBaseCurrency}} {{.UserBaseCurrency}}*.
{{ $userBaseCurrency := .UserBaseCurrency }}
Current rates are:
{{range $r := .Rates}}
*1 {{ $r.Currency }} = {{ $r.Rate }} {{ $userBaseCurrency }}*
{{end}}
Click to remove: