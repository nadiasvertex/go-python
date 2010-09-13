#!/bin/bash

6g interp/*.go
6g pkg/core/parser/*.go
6g pkg/core/types/*.go
6g pkg/core/vm/*.go
6l -o gopy gopy.6
