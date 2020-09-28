Your total savings is *{{.TotalInUserBaseCurrency}} {{.UserBaseCurrency}}*.
{{ $userBaseCurrency := .UserBaseCurrency }}
Current rates are:
{{range $currency, $rate := .SavingsRates}}
*1 {{ $currency }} = {{ $rate }} {{ $userBaseCurrency }}*
{{end}}
Click to remove: