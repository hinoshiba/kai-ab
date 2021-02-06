package main

const (
	VERSION_TEMPLATE string = "v0.0.1"

	PATH_ETC_TEMPLATE string = "./etc/template.yaml"
	PATH_ETC_FILTERS string = "./etc/auto_filters.yaml"
	PATH_ENV string = ".kai-ab.env"
	PATH_REPORT string = "./var/report/"
	PATH_CSV_BOTH string = "./var/csv/"
	PATH_CSV_IN string = "./in/"
	PATH_CSV_OUT string = "./out/"

	PATH_FMT_DATE string = "200601" //yyyymm

	STR_HELP string = VERSION_TEMPLATE + `
Usage: kai-ab <sub command> [<subcommand option>]
Support of categorize and calc account book.

Subcommand
  init <path>        : Initialize path if does not initialized.
  template <yyyymm>  : Create directory of month '<yyyymm>' and Create Template file from 'etc/template.yaml'.
  autofil [<path>]   : Auto categorize row at csv. rule of 'etc/autofil.yaml'.
  mfil [<path>]      : Start the manual filter mode.
  check              : Check the all csv files.
  calc               : Generate report from csv files.
`

)
