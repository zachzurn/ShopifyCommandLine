package main

import (
    "fmt"
    "os"
    "github.com/mgutz/ansi"
)

type Settings struct {

    AllowedExtensions []string `json:"allowed_extensions"`

    Stores map[string]Store `json:"stores"`

}


/*
Ensure the settings are valid and that the required stores are set up
*/
func (s *Settings) validate(requiredStores []string){
    
    stores := s.Stores

    for key, store := range stores {
        store.validate(key)
    }

    for _, v := range requiredStores {
        if _, ok := stores[v]; ok == false {
            fmt.Printf("%v Missing target '%v'. Please add this target under 'stores' in your shopify.json file.\n",StrError,ansi.Color(v, "134"))
            os.Exit(1)
        }
    }

}