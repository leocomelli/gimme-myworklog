# Gimme my worklog

This project will create a csv file with all worklog registered at Jira (Tempo plugin).

### Compile (cross compilation)
* [Installing Go from source](http://golang.org/doc/install/source)
* [An introduction to cross compilation with Go 1.1](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1)

### Usage:

	main -url=<JIRA_URL> \
	     -startDate=<YYYY-MM-DD> \
	     -endDate=<YYYY-MM-DD> \
	     -tempoApiToken=<TEMPO_API_TOKEN> \
	     -username=<USERNAME> \
	     -password=<PASSWORD> \
	     -output=<OUTPUT_FILE>