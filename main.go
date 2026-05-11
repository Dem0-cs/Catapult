package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

var version string = "1.1.0"
var catapultInfo string = "https://catapult.if0n.xyz/info.json"
var url string = "https://github.com/Dem0-cs/Catapult-Data/releases/latest/download/meowcoin_bootstrap.zip"
var useCustomBoostrap bool = false
var customLink string


type CatapultData struct {
    Version string `json:"version"`
    Hash1   string `json:"part1_hash"`
    Hash2   string `json:"part2_hash"`
}

var info CatapultData

func main() {
    client := http.Client{
        Timeout: 5 * time.Second,
    }

    req, err := client.Get(catapultInfo)
    if err != nil {
        fmt.Println("\n[!] WARNING: Could not connect to the Catapult API.")
        fmt.Println("    Reason:", err)
        fmt.Println("-------------------------------------------------------")
        fmt.Println("We cannot verify if your Meowcoin bootstrap is authentic.")
        fmt.Print("Do you want to proceed at your own risk? (y/N): ")

        var response string
        fmt.Scanln(&response)
        if strings.ToLower(response) != "y" {
            fmt.Println("Please check your internet or try again later.")
            return
        }
    }

    defer req.Body.Close()
    bodyBytes, err := io.ReadAll(req.Body)
    if err != nil {
        fmt.Println("Error reading response:", err)
        return
    }


    err = json.Unmarshal(bodyBytes, &info)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    if(version != info.Version) {
        fmt.Print("⚠️ A new version of Catapult has been detected! It is highly recommended you update.")
        fmt.Print("Continue anyways? (y/N): ")

        var confirmVersion string
        fmt.Scanln(&confirmVersion)

        if strings.ToLower(confirmVersion) != "y" {
            fmt.Print("Download the latest update from our Github https://github.com/Dem0-cs/Catapult")
            fmt.Println("Press 'Enter' to exit...")
                    
            fmt.Scanln()
            return
        }
    }

    home, _ := os.UserHomeDir()
    var meowFolder string

    if runtime.GOOS == "windows" {
        roaming, _ := os.UserConfigDir()
        meowFolder = filepath.Join(roaming, "Meowcoin")
    } else {
        meowFolder = filepath.Join(home, ".meowcoin")
    }

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
        fmt.Println("Press 'Enter' to exit...")
                    
        fmt.Scanln() 
        return
    }

    fmt.Println("🚀 Initializing Download...")
    DownloadLatest(meowFolder, !useCustomBoostrap, info)
}


func calculateHash(body io.Reader) (string, error) {
    h := sha256.New()
    if _, err := io.Copy(h, body); err != nil {
        return "", err
    }
    return hex.EncodeToString(h.Sum(nil)), nil
}



func DownloadLatest(targetFolder string, catapultServers bool, info CatapultData) {
    zipPath := filepath.Join(targetFolder, "meowcoin_bootstrap.zip")
    out, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil { log.Fatal(err) }
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
        for i, p := range parts {
            partUrl := fmt.Sprintf("https://github.com/Dem0-cs/Catapult-Data/releases/latest/download/meowcoin_bootstrap.zip.%s", p)
            fmt.Printf("📥 Catapulting Part %s...\n", p)
            
            resp, err := http.Get(partUrl)
            if err != nil { log.Fatal(err) }
            
            h := sha256.New()
            bar := progressbar.DefaultBytes(resp.ContentLength, "🚀 Part "+p)
            
            multiWriter := io.MultiWriter(out, bar, h)
            _, err = io.Copy(multiWriter, resp.Body)
            resp.Body.Close()

            if err != nil { log.Fatal(err) }

            localHash := hex.EncodeToString(h.Sum(nil))
            remoteHash := info.Hash1
            if i == 1 { remoteHash = info.Hash2 } 

            if localHash != remoteHash {
                fmt.Printf("\n❌ HASH MISMATCH on Part %s!\n", p)
                fmt.Printf("Expected: %s\nActual:   %s\n", remoteHash, localHash)
                
                out.Close() 
                os.Remove(zipPath)

                fmt.Println("-------------------------------------------------------")
                fmt.Println("The downloaded data is corrupted or has been tampered with.")
                fmt.Println("Press 'Enter' to exit...")
                
                fmt.Scanln() 
                return 
            }
            fmt.Printf("✅ Part %s Verified\n", p)
        }
    }

    fmt.Println("\n✅ All data received and verified. Extracting...")
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