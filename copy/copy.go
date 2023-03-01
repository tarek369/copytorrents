package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	//"PTN"

	PTN "github.com/middelink/go-parse-torrent-name"
	cp "github.com/otiai10/copy"
	"golang.org/x/exp/slices"
)

var othertype = []string{".nfo", ".jpg", ".bif"}
var pics = []string{"screenshot", "screenshots", "screen"}
var source = "/downloads/torrents/rutorrent/completed/"
var destination_folder = "/mnt/unionfs/Media/Ramadan"
var destination_folder_movies = "/mnt/unionfs/Media/ArabicMovies"
var destination_folder_tveng = "/mnt/unionfs/Media/TV"
var destination_folder_anime = "/mnt/unionfs/Media/Anime"
var types = []string{".mkv", ".mp4", ".avi", ".ts", ".mpg"}

func main() {

	label := os.Args[2]
	downloaded_file := os.Args[1]
	source = source + downloaded_file
	file, err := os.OpenFile("/config/log/copy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		fmt.Println("Fucking error!", err)
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("====== STARTING (====== ")
	log.Println("command = ./copy ", downloaded_file)
	log.Println("Source episode directory = ", source)

	switch label {
	case "tvfullarab":
		_, to := copy(downloaded_file, destination_folder)
		chown(to)
	case "tveng":
		_, to := copy(downloaded_file, destination_folder_tveng)
		chown(to)
	case "animetv":
		_, to := copy(downloaded_file, destination_folder_anime)
		chown(to)
	case "arabmovie":
		_, to := copy(downloaded_file, destination_folder_movies)
		chown(to)
	case "tv":
		files, _ := ioutil.ReadDir(source)
		var episode_full_path string
		for _, file := range files {
			filename := file.Name()
			fileext := filepath.Ext(filename)
			if slices.Contains(types, fileext) {
				episode_full_path = filepath.Join(source, filename)
				break
			}
		}
		if episode_full_path == "" {
			episode_full_path = source

		}
		clean(episode_full_path)
		//copy(episode_full_path, label)
		//episode_full_path = filepath.Join(source, [f for f in os.listdir(
		//	source) if f.endswith(('.mkv', '.mp4', '.ts', '.mpg'))][0])
	}

}

func copy(downloaded_file, dest string) (bool, string) {

	to := filepath.Join(dest, downloaded_file)
	log.Println("Copying\n ", downloaded_file, to)
	err := cp.Copy(source, to)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Folder copied successfully")
	}

	return true, to
}

func chown(to string) {
	us, _ := user.Lookup("abc")
	uid, _ := strconv.Atoi(us.Uid)
	gid, _ := strconv.Atoi(us.Gid)
	err := os.Chown(to, uid, gid)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("File ownership changed successfully")
		log.Println(to)
	}
}

func clean(episode_full_path string) {
	// Try to get Folder name
	r1 := regexp.MustCompile(`.E(\d{1,2})`)
	episode_file := filepath.Base(episode_full_path)
	parse_results, _ := PTN.Parse(episode_file)
	epname := strings.Replace(parse_results.Title, " ", ".", -1)
	epname = strings.Split(epname, ".R23")[0]
	epname = strings.Split(epname, ".R023")[0]
	epname = strings.Split(epname, ".S1")[0]
	epname = strings.Split(epname, ".S0")[0]
	epname = strings.Split(epname, "_S1")[0]
	epname = strings.Split(epname, ".Ep")[0]
	epname = strings.Split(epname, ".EP")[0]
	epname = strings.Split(epname, ".SNone")[0]
	epname = r1.Split(epname, -1)[0]

	// Try to Clean Episode name
	r2 := regexp.MustCompile(`R\d\d.`)
	r3 := regexp.MustCompile("EP")
	episode_pure := r2.ReplaceAllString(episode_file, ".S01")
	episode_pure = r3.ReplaceAllString(episode_file, "E")
	episode_pure = strings.Replace(episode_pure, ".SHAHID.WEB-DL.AAC2.0.H.264.BY.RoMaNTiCPoET", "", -1)
	episode_pure = strings.Replace(episode_pure, "SHAHID", "", -1)
	episode_pure = strings.Replace(episode_pure, "H.264", "", -1)
	episode_pure = strings.Replace(episode_pure, "BY.RoMaNTiCPoET", "", -1)
	episode_pure = strings.Replace(episode_pure, " ", "", -1)
	episode_pure = strings.Replace(episode_pure, "..", ".", -1)
	episode_pure = strings.Replace(episode_pure, "_", ".", -1)
	episode_pure = strings.Replace(episode_pure, "..", ".", -1)
	episode_pure = strings.Replace(episode_pure, "._", "", -1)
	log.Println("File name", epname)
	dest_dir := filepath.Join(destination_folder, epname)
	dest_path := filepath.Join(dest_dir, episode_pure)

	log.Println("Copying\n ", episode_full_path, dest_path)

	_ = os.Mkdir(dest_dir, 0775)
	err := cp.Copy(episode_full_path, dest_path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Folder copied successfully")
		log.Println("Copied succsufully")
	}
	chown(dest_path)

}
