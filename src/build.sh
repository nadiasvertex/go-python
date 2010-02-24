#!/bin/bash

~/bin/8g interp/*.go
~/bin/8g pkg/core/parser/*.go
~/bin/8l -o gopy *.8
