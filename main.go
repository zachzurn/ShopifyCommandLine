package main

import (
    "fmt"
    "os"
    "encoding/json"
    "io/ioutil"
)

var textExtensions = []string{".liquid", ".js", ".css", ".json", ".html", ".txt", ".md", ".svg"}
var allowedExtensions = []string{".liquid", ".js", ".css", ".json", ".html", ".txt", ".png", ".jpg", ".jpeg", ".gif", ".eot", ".woff", ".svg"}

func (s *Store) applyAction(action string){
    switch action{
        case "download": 
            s.download()
        case "deploy": 
            s.deploy()
        case "watch": 
            s.watch()
        default:
            fmt.Printf("The action '%v' is not allowed. The possible actions are 'download' , 'deploy' and 'watch'\n",action)
            os.Exit(1)
    }
}

func (s *Store) applyActionTo(action string, dest *Store){

}

func main(){

	file, e := ioutil.ReadFile("./shopify.json")
    if e != nil {
        fmt.Printf("%v No shopify.json file found. Please create one.\n",StrError)
        os.Exit(1)
    }

    var settings Settings
    json.Unmarshal(file, &settings)

    for _, ext := range settings.AllowedExtensions {
        allowedExtensions = append(allowedExtensions, ext)
    }

    action, sourceTarget, destTarget := "","",""
    args := os.Args[1:]

    switch len(args) {
        case 2:
            action, sourceTarget = args[0], args[1]
            settings.validate([]string{sourceTarget})
            store := settings.Stores[sourceTarget]
            store.applyAction(action)
        case 3:
            action, sourceTarget, destTarget = args[0], args[1], args[2]
            settings.validate([]string{sourceTarget,destTarget})
            store := settings.Stores[sourceTarget]
            destStore := settings.Stores[destTarget]
            store.applyActionTo(action,&destStore)
        default:
            fmt.Printf("%v Please type a target and action. Example 'download production'.\n",StrError)
            os.Exit(1)            
    }

    
}