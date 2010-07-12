#!/bin/bash

6g interp/*.go
6g pkg/core/parser/*.go
6l -o gopy *.6
