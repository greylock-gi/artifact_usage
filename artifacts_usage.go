package main

import (
	"encoding/csv"
	"os"
	"io"
	"fmt"
	"strings"
	"sort"
)

type User struct {		// possible user of artifact
	name string			// name (i.e. "Zhongli")
	build string		// build (i.e. "Burst support")
	substats []string	// desired substats (i.e. {"CRIT", "ER%", "HP%"})
}

var sands_mainstats = []string {
	"ATK%", "HP%", "DEF%", "ER%", "EM",
}

var goblet_mainstats = []string {
	"ATK%", "HP%", "DEF%", "EM",
	"Anemo DMG%", "Geo DMG%", "Electro DMG%", "Dendro DMG%", "Hydro DMG%", "Pyro DMG%", "Cryo DMG%", "Physical DMG%",
}

var circlet_mainstats = []string {
	"ATK%", "HP%", "DEF%", "EM",
	"CRIT", "Crit Rate", "Healing Bonus",
}

var all_mainstats = []string {
	"ATK%", "HP%", "DEF%", "ER%", "EM",
	"Anemo DMG%", "Geo DMG%", "Electro DMG%", "Dendro DMG%", "Hydro DMG%", "Pyro DMG%", "Cryo DMG%", "Physical DMG%",
	"CRIT", "Crit Rate", "Healing Bonus",
}

var all_substats = []string {
	"ATK%", "HP%", "DEF%", "ER%", "EM",
	"Flat ATK", "Flat HP", "Flat DEF",
	"CRIT", "Crit Rate",
}	

func main() {
	artifacts, err := loadCSV("artifacts.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	builds, err := loadCSV("builds.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	
	fmt.Printf("<html>\n" +
		"\t<head>\n" +
		"\t\t<link rel=\"stylesheet\" href=\"artifact_usage.css\" />\n" +
		"\t</head>\n" +
		"\t<body>\n")
	for _, artifact := range artifacts {
		fmt.Printf("\t\t<h2>%v</h2>\n" +
			"\t\t<table>\n" +
			"\t\t\t<tr>\n" +
			"\t\t\t\t<td>2pc</td>\n" +
			"\t\t\t\t<td>%v</td>\n" +
			"\t\t\t</tr>\n" +
			"\t\t\t<tr>\n" +
			"\t\t\t\t<td>4pc</td>\n" + 
			"\t\t\t\t<td>%v</td>\n" +
			"\t\t\t</tr>\n" +
			"\t\t</table>\n", artifact[0], artifact[3], artifact[4])
		flower := make([]User, 0)
		sands := make(map[string][]User)
		goblet := make(map[string][]User)
		circlet := make(map[string][]User)
		for _, build := range builds {
			if  strings.Contains(build[6], artifact[1]) || strings.Contains(build[6], artifact[2]) {
				var user User;
				
				user.name = build[0]
				user.build = build[1]
				substats, err := split_substats(build[5])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unknown substat in \"%v\" (%v %v)\n", build[5], build[0], build[1]);
				}
				user.substats = substats

				flower = add_user(flower, "", user)

				mainstats, err := split_mainstats(build[2], sands_mainstats)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unknown sands mainstat in \"%v\" (%v %v)\n", build[2], build[0], build[1]);
					continue
				}
				for _, mainstat := range mainstats {
					sands[mainstat] = add_user(sands[mainstat], mainstat, user)
				}

				mainstats, err = split_mainstats(build[3], goblet_mainstats)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unknown goblet mainstat in \"%v\" (%v %v)\n", build[3], build[0], build[1]);
					continue
				}
				for _, mainstat := range mainstats {
					goblet[mainstat] = add_user(goblet[mainstat], mainstat, user)
				}

				mainstats, err = split_mainstats(build[4], circlet_mainstats)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unknown circlet mainstat in \"%v\" (%v %v)\n", build[4], build[0], build[1]);
					continue
				}
				for _, mainstat := range mainstats {
					circlet[mainstat] = add_user(circlet[mainstat], mainstat, user)
				}
			}
		}
		list_users(flower, get_icon_tag(artifact[5], 4) + "<br>Flower<br />" + get_icon_tag(artifact[5], 2) + "<br>Feather")
		list_users2(sands, "Sands", get_icon_tag(artifact[5], 5))
		list_users2(goblet, "Goblet", get_icon_tag(artifact[5], 1))
		list_users2(circlet, "Circlet", get_icon_tag(artifact[5], 3))
	}

	// Off-set pieces
	flower := make([]User, 0)
	sands := make(map[string][]User)
	goblet := make(map[string][]User)
	circlet := make(map[string][]User)
	for _, build := range builds {
		var user User;
			
		user.name = build[0]
		user.build = build[1]
		substats, err := split_substats(build[5])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown substat in \"%v\" (%v %v)\n", build[5], build[0], build[1]);
			continue
		}
		user.substats = substats

		flower = add_user(flower, "", user)

		mainstats, err := split_mainstats(build[2], sands_mainstats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown sands mainstat in \"%v\" (%v %v)\n", build[2], build[0], build[1]);
			continue
		}
		for _, mainstat := range mainstats {
			sands[mainstat] = add_user(sands[mainstat], mainstat, user)
		}

		mainstats, err = split_mainstats(build[3], goblet_mainstats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown goblet mainstat in \"%v\" (%v %v)\n", build[3], build[0], build[1]);
			continue
		}
		for _, mainstat := range mainstats {
			goblet[mainstat] = add_user(goblet[mainstat], mainstat, user)
		}

		mainstats, err = split_mainstats(build[4], circlet_mainstats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown circlet mainstat in \"%v\" (%v %v)\n", build[4], build[0], build[1]);
			continue
		}
		for _, mainstat := range mainstats {
			circlet[mainstat] = add_user(circlet[mainstat], mainstat, user)
		}
	}
	fmt.Printf("\t\t<h2>As an off-set piece</h2>\n")
	list_users(flower, "Flower<br />Feather")
	list_users2(sands, "Sands", "")
	list_users2(goblet, "Goblet", "")
	list_users2(circlet, "Circlet", "")

	fmt.Printf("\t</body>\n" +
		"</html>\n")
}

func add_user(users []User, mainstat string, user User) []User {
	filtered_substats := []string {}
	crit_is_present := false
	for _, s := range user.substats {
		if s == "CRIT" {
			crit_is_present = true
			break
		}
	}
	for _, s := range user.substats {
		if s == "Crit Rate" && crit_is_present {
			continue
		}
		if s == mainstat && mainstat != "CRIT" {
			continue
		}
		filtered_substats = append(filtered_substats, s)
	}
	u := user
	u.substats = filtered_substats
	return append(users, u)
}


func get_icon_tag(set_id string, typ int) string {
	return "<img src=\"https://enka.network/ui/UI_RelicIcon_" +  set_id + "_" + fmt.Sprintf("%d", typ) + ".png\" width=\"64\" height=\"64\" />"
}

func list_users2(artifacts map[string][]User, title string, icon_tag string) {
	for _, mainstat := range all_mainstats {
		artifact, ok := artifacts[mainstat]
		if ok {
			list_users(artifact, icon_tag + "<br />" + mainstat + " " + title)
		}
	}
}

func list_users(artifact []User, first_cell_text string) {
	sort.Slice(artifact, func(i, j int) bool {
		r := compare_slices_of_strings(artifact[i].substats, artifact[j].substats)
		if r == 0 {
			return compare_strings(artifact[i].name, artifact[j].name) == -1
		}
		return r == -1
	})

	cnt := 0
	for i, user := range artifact {
		if i != len(artifact) - 1 && compare_slices_of_strings(user.substats, artifact[i+1].substats) == 0 {
			continue
		}
		cnt++;
	}

	fmt.Printf("\t\t<table class=\"usage\">\n")
	fmt.Printf("\t\t\t<tr>\n")
	fmt.Printf("\t\t\t\t<td class=\"usage-type\" rowspan=\"%v\">%v</td>\n", cnt, first_cell_text)
	fmt.Printf("\t\t\t\t<td class=\"usage-users\">")
	for i, user := range artifact {
		if user.build != "" {
			fmt.Printf("<b>%v</b> (%v)", user.name, user.build)
		} else {
			fmt.Printf("<b>%v</b>", user.name)
		}
		if i != len(artifact) - 1 && compare_slices_of_strings(user.substats, artifact[i+1].substats) == 0 {
			fmt.Printf(", ")
		} else {
			fmt.Printf("</td>\n")
			fmt.Printf("\t\t\t\t<td class=\"usage-substats\">")
			for j, substat := range user.substats {
				fmt.Printf("%v", substat)
				if j != len(user.substats) - 1 {
					fmt.Printf(", ")
				}
			}
			fmt.Printf("</td>\n")
			fmt.Printf("\t\t\t</tr>\n")
			if i != len(artifact) - 1 {
				fmt.Printf("\t\t\t<tr>\n")
				fmt.Printf("\t\t\t\t<td class=\"usage-users\">")
			}
		}
	}
	fmt.Printf("\t\t</table>\n")
}


func loadCSV(filename string) ([][]string, error) {
	result := make([][]string, 0)
	fi, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	reader := csv.NewReader(fi)
	// skip first line
	_, err = reader.Read()
	if err == io.EOF {
		return result, nil
	}
	if err != nil  {
		return result, err
	}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil  {
			return result, err
		}
		// skip empty and commented-out lines
		if len(line) == 0 {
			continue
		}
		if len(line[0]) == 0 || line[0][0] == '#' {
			continue
		}
		result = append(result, line)
	}
	fi.Close()
	return result, nil
}

func split_mainstats(list string, stats []string) ([]string, error) {
	result, err := split_stats(list, stats)
	if err == nil && len(result) == 0 {	// empty mainstat means "any mainstat"
		return stats, nil
	}
	return result, err
}

func split_substats(list string) ([]string, error) {
	return split_stats(list, all_substats)
}

func split_stats(list string, stats []string) ([]string, error) {
	var result []string

	for {
		list = strings.TrimLeft(list, " /|,;>≥=~")
		if list == "" {
			break
		}
		found := false
		for _, stat := range stats {
			if strings.HasPrefix(list, stat) {
				result = append(result, stat)
				list = strings.TrimPrefix(list, stat)
				found = true
			}
		}
		if !found {
			return []string {}, fmt.Errorf("Unknown stat")
		}
	}
	sort.Strings(result)
	return result, nil
}

func same_substats(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return false
}

func compare_strings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compare_slices_of_strings(a, b []string) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		r := compare_strings(a[i], b[i])
		if r != 0 {
			return r
		}
	}
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}
