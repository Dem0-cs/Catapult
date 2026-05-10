package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"archive/zip"

    "github.com/schollz/progressbar/v3"
)

var url string = "https://github.com/Dem0-cs/Catapult-Data/releases/latest/download/meowcoin_bootstrap.zip"
var useCustomBoostrap bool = false
var customLink string

func main() {
    roaming, _ := os.UserConfigDir()
    meowFolder := filepath.Join(roaming, "Meowcoin")

    var serverChoice string
    fmt.Print("🛜 Use Catapult servers? (Y/n): ")
    fmt.Scanln(&serverChoice)
    
    if strings.ToLower(serverChoice) == "n" || strings.ToLower(serverChoice) == "no" {
        fmt.Print("🔗 Please paste the link to your custom .zip: ")
        fmt.Scanln(&customLink)
        useCustomBoostrap = true
    }

    fmt.Print("⚠️ This WIPES blocks, chainstate, assets, and indexes folders. Ready to Catapult? (y/N): ")
    var confirm string
    fmt.Scanln(&confirm)

    if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
        fmt.Println("Catapult aborted. No files were changed.")
        os.Exit(0)
    }

    fmt.Println("🚀 Initializing Download...")
    DownloadLatest(meowFolder, !useCustomBoostrap)
}


func DownloadLatest(targetFolder string, catapultServers bool) {
    zipPath := filepath.Join(targetFolder, "meowcoin_bootstrap.zip")

    out, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        log.Fatalf("Could not create file: %v", err)
    }
    defer out.Close()

    if !catapultServers {
        fmt.Println("📥 Downloading from custom link...")
        resp, err := http.Get(customLink)
        if err != nil { log.Fatal(err) }
        defer resp.Body.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "🚀 Downloading")
        io.Copy(io.MultiWriter(out, bar), resp.Body)
    } else {
        parts := []string{"001", "002"}
        for _, p := range parts {
            partUrl := fmt.Sprintf("https://github.com/Dem0-cs/Catapult-Data/releases/latest/download/meowcoin_bootstrap.zip.%s", p)
            fmt.Printf("📥 Catapulting Part %s...\n", p)
            
            resp, err := http.Get(partUrl)
            if err != nil { log.Fatal(err) }
            
            bar := progressbar.DefaultBytes(resp.ContentLength, "🚀 Part "+p)
            io.Copy(io.MultiWriter(out, bar), resp.Body)
            resp.Body.Close()
        }
    }

    fmt.Println("\n✅ All data received. Extracting...")
    ExtractFolder(targetFolder)
}


func ExtractFolder(targetFolder string) {
    zipFile := filepath.Join(targetFolder, "meowcoin_bootstrap.zip")


    foldersToRemove := []string{"blocks", "chainstate", "assets", "indexes"}
    fmt.Println("🧹 Wiping old data...")
    
    for _, folder := range foldersToRemove {
        fullPath := filepath.Join(targetFolder, folder)
        err := os.RemoveAll(fullPath)
        if err != nil {
            log.Printf("Warning: could not remove %s: %v (Is the wallet still open?)", folder, err)
        }
    }

    r, err := zip.OpenReader(zipFile)
    if err != nil {
        log.Fatalf("Failed to open zip: %v", err)
    }
    defer r.Close()

    fmt.Println("📦 Extracting blocks... please wait.")

    for _, f := range r.File {
        fpath := filepath.Join(targetFolder, f.Name)

        if !strings.HasPrefix(fpath, filepath.Clean(targetFolder)+string(os.PathSeparator)) {
            continue
        }


        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }


        if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            log.Fatal(err)
        }


        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            log.Fatal(err)
        }

        rc, err := f.Open()
        if err != nil {
            outFile.Close()
            log.Fatal(err)
        }

        _, err = io.Copy(outFile, rc)
        outFile.Close()
        rc.Close()

        if err != nil {
            log.Fatal(err)
        }

    }

    os.Remove(zipFile)
    fmt.Println("✨ Success! Meowcoin has been Catapulted!")
}