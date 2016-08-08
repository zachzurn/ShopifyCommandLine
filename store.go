package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "strings"
    "time"
    "encoding/base64"
    "path/filepath"
    "github.com/zachzurn/go_shopify"
    "github.com/fsnotify/fsnotify"
    "github.com/mgutz/ansi"
)


type Store struct {
    Name string `json:"name"`

    Folder string `json:"folder"`

    ApiKey string `json:"api_key"`

    ApiPassword string `json:"api_password"`

    Shop string `json:"shop"`

    ThemeId int64 `json:"theme_id"`

    Api *shopify.API
}

func(s *Store) connect(){
    //Ready to go, settings validated
    if s.Api != nil {
        return
    }

    s.Api = &shopify.API{
            Shop:    s.Shop,
            Token:  s.ApiKey,
            Secret: s.ApiPassword,
    }
}

func (s *Store) download(){

    s.connect()

    localFolder := "./"+s.Folder+"/"

    assets, err := s.Api.Assets(s.ThemeId)


    if err != nil {
        fmt.Printf("%v Error fetching assets.\n%v %v\n", StrError, ErrorSpacer, err)
        os.Exit(1)
    }

    for _, asset := range assets {
        assetKey := strings.Replace(asset.Key, "\\/", "/", -1)
        localAssetPath := localFolder + assetKey

        fullAsset, err := s.Api.Asset(s.ThemeId,assetKey)
        
        if err != nil {
            fmt.Printf("Error fetching full asset: %v : %v\n", assetKey, err)
            continue
        }

        isBinary := fullAsset.Attachment != ""

        err = os.MkdirAll(filepath.Dir(localAssetPath), os.ModePerm)
        if err != nil {
                fmt.Printf("%s  %s \n%s %s\n",StrWarning, localAssetPath, WarningSpacer, err)
                return
        }   

        if isBinary {

            data, err := base64.StdEncoding.DecodeString(fullAsset.Attachment)
            if err != nil {
                fmt.Printf("%s  %s \n%s %s\n",StrWarning, localAssetPath, WarningSpacer, err)
                continue
            }

            err = ioutil.WriteFile(localAssetPath, data, 0644)

            if err != nil {
                fmt.Printf("%s  %s \n%s %s\n",StrWarning, localAssetPath, WarningSpacer, err)
                continue
            }

        } else {

            err = ioutil.WriteFile(localAssetPath, []byte(fullAsset.Value), 0644)

            if err != nil {
                fmt.Printf("%s  %s \n%s %s\n",StrWarning, localAssetPath, WarningSpacer, err)
                continue
            }

        }

        fmt.Printf("%s  %s\n",StrDownloaded, localAssetPath)

        time.Sleep(200 * time.Millisecond)
    }

}

func (s *Store) watch(){

    s.connect()

    fmt.Printf("%v %s\n",StrWatching,ansi.Color(s.Folder, "134"))

    watcher, err := NewBatcher(600 * time.Millisecond)
    if err != nil {
        fmt.Printf("%s  %s\n",StrError, err)
        os.Exit(1)
    }

    isDirectory := func(pth string) (bool, error) {
         fi, err := os.Stat(pth)
         if err != nil {
                 return false, err
         }

         return fi.IsDir(), nil
    }


    watch := func(folder string){
        err = watcher.Add(folder)
        if err != nil {
            fmt.Printf("%s  %s\n",StrError, err)
            os.Exit(1)
        }
    }


    watchWalk := func(folder string){
        filepath.Walk(folder,func(path string, info os.FileInfo, err error) error {

            if d, _ := isDirectory(path); d{
                watch(path)
            }

            if err != nil {
                return err
            }

            return nil

        })
    }

    unwatch := func(folder string){
        err = watcher.Remove(folder)
        if err != nil {
            fmt.Printf("%s  %s\n",StrError, err)
            os.Exit(1)
        }
    }

    defer watcher.Close()

    done := make(chan bool)
    go func() {
        for {
            select {
            case events := <-watcher.Events:
                
                event := events[len(events)-1]
                path := event.Name

                if event.Op&fsnotify.Create == fsnotify.Create {
                    if d, _ := isDirectory(event.Name); d{
                        //Directory was created, let's watch it
                        watch(path)
                    } else {
                        //File was created, lets upload it to Shopify
                        s.uploadLocalAsset(path);
                    }
                }


                if event.Op&fsnotify.Write == fsnotify.Write {
                    //File was written, let's upload the file to Shopify
                    s.uploadLocalAsset(path);
                }

                if event.Op&fsnotify.Remove == fsnotify.Remove {
                    if d, _ := isDirectory(path); d{
                        //Directory was deleted, lets unwatch the directory
                        unwatch(path)
                    } else {
                        //File was deleted, lets delete the file from Shopify
                        s.deleteRemoteAsset(path)
                    }
                }

            case err := <-watcher.Errors:
                fmt.Printf("%s  %s\n",StrWarning, err)
            }
        }
    }()

    watchWalk(s.Folder + "/" + "assets")
    watchWalk(s.Folder + "/" + "config") 
    watchWalk(s.Folder + "/" + "layout")
    watchWalk(s.Folder + "/" + "locales")
    watchWalk(s.Folder + "/" + "snippets")
    watchWalk(s.Folder + "/" + "templates")

    <-done

}

func (s *Store) deleteRemoteAsset(pth string) {
    pth = strings.Replace(pth,`\`,`/`,-1)
    assetKey, err := filepath.Rel(s.Folder, pth)
    assetName := filepath.Base(assetKey)
    
    if err != nil{
        
    } 

    err = s.Api.Delete(s.ThemeId,strings.Replace(assetKey,`\`,`/`,-1))

    if err != nil {
        fmt.Printf("%v Couldn't remove asset in Shopify. %v\n%v %v\n",StrWarning,pth,WarningSpacer,err)
        return
    }

    fmt.Printf("%v %v\n",StrRemoved,assetName)

}

func (s *Store) uploadLocalAsset(pth string) {

    pth = strings.Replace(pth,`\`,`/`,-1)
    assetKey, err := filepath.Rel(s.Folder, pth)
    assetName := filepath.Base(assetKey)
    assetExt := filepath.Ext(assetName)

    if err != nil{
        
    } 

    if !extensionAllowed(assetExt) {
        fmt.Printf("%v Extension %v ignored. Add in shopify.json as 'allowed_extensions'.\n%v Allowed -> %v\n",StrWarning,ansi.Color(assetExt, "134"),WarningSpacer,allowedExtensions)
        return;
    }

    asset := s.Api.NewAssetUpload()

    asset.Key = strings.Replace(assetKey,"\\","/",-1)

    file, e := ioutil.ReadFile(pth)
    
    if e != nil {
        fmt.Printf("%v Couldn't read file. %v\n",StrWarning,pth)
        return
    }

    if isTextFile(assetExt) {

        asset.Value = string(file[:])

    } else {

        asset.Attachment = base64.StdEncoding.EncodeToString(file)

    }


    err = asset.Upload(s.ThemeId)

    if err != nil {
        fmt.Printf("%v Couldn't upload file. %v\n%v %v\n",StrWarning,assetKey,WarningSpacer,err)
        return
    }

    fmt.Printf("%v %v\n",StrUploaded,assetName)
}



func (s *Store) deploy(){

    fmt.Printf("%v Deploy functionality is not ready yet.\n",StrWarning)
    //Get a list of all local assets

    //Get a list of remote assets

    //Upload all local assets

    //Remove and remote assets not in the local list
}

/*
Ensure the store is valid by checking that the correct parameters exist
*/
func (s *Store) validate(storeKey string){
    
    if s.Name == "" {
        fmt.Printf("%v Store '%s' is missing the 'name' field. Enter in the name for the target.\n",StrError,StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

    if s.Folder == "" {
        fmt.Printf("%v Store '%s' 'is missing the 'folder' field. Enter in the name of the folder that the theme contents will go.\n",StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

    if s.ApiKey == "" {
        fmt.Printf("%v Store '%s' is missing the 'api_key' field.\n",StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

    if s.ApiPassword == "" {
        fmt.Printf("%v Store '%s' is missing the 'api_password' field.\n",StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

    if s.Shop == "" {
        fmt.Printf("%v Store '%s' is missing the 'shop' field. This should be the beginning part of the shopify url.\n",StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

    if s.ThemeId == 0 {
        fmt.Printf("%v Store '%s' is missing the 'theme_id' field.\n",StrError,ansi.Color(storeKey, "134"))
        os.Exit(1)
    }

}



/* UTILITY */
func extensionAllowed(ext string) bool {
    for _, a := range allowedExtensions {
        if a == ext {
            return true
        }
    }
    return false
}

func isTextFile(ext string) bool {
    for _, a := range textExtensions {
        if a == ext {
            return true
        }
    }
    return false
}