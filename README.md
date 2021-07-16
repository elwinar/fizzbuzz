# FizzBuzz API

Basic implementation of a fizz-buzz REST server for the technical screening of
LeBonCoin.

Implementation is kept as simple as possible and without external library. The
language used is Golang, and a makefile with `build` and `test` targets is
available for ease of use. The provided tests aren't complete to keep them
short but demonstrate the nominal case of the API and check compliance with the
instruction.

## Ready for production

Instructions are very vague, so this implementation considers that "ready for
production" means:
- "well-behvaed" in the unix sense:
	- logs in stdout
	- handling of signals
	- correct exit code on error
	- make commands for building
	- command-line parameters for chosing the listening interface
	- help message with the available options
- user-friendly
	- clear interface with a documentation
	- consistent output
	- correct error messages & status codes

No particular effort has been made towards instrumentation of the code for
monitoring, like logging all input requests or exporting metrics, as it is very
platform-specific.

Likewise, not effort towards specific deployment infrastructure has been made,
like providing a dockerfile for building containers, for the same reasons.
