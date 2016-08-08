package main

import (
    "github.com/mgutz/ansi"
)

var StrWatching =    ansi.Color("--- WATCHING", "cyan")
var StrError =       ansi.Color("------ ERROR", "red")
var StrWarning =     ansi.Color("---- WARNING", "214")
var StrRemoved =     ansi.Color("---- REMOVED", "magenta")
var StrAdded =       ansi.Color("------ ADDED", "green")
var StrDownloaded =  ansi.Color("- DOWNLOADED", "green")
var StrUploaded =    ansi.Color("--- UPLOADED", "green")
var StrUpdated =     ansi.Color("---- UPDATED", "cyan")
var ErrorSpacer =    ansi.Color("------------", "red")
var WarningSpacer =  ansi.Color("------------", "214")
var StrSpacer =      "------------"