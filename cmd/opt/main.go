package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shizukayuki/excel-hk4e"
	"github.com/shizukayuki/ysoptimizer/pkg/good"
)

var (
	optimized = map[good.CharacterKey]OptimizeState{}
	slow      = flag.Bool("slow", false, "Run in slow mode. Don't filter artifacts that have dead stats")
	path      = flag.String("good-file", "./GOOD.json", "Location of GOOD.json")
	excelPath = flag.String("genshin-data", "", "Specify GenshinData path. Defaults to $HOME/git/GenshinData or $GENSHIN_DATA_REPO if set")
)

func main() {
	flag.Parse()

	repo := os.Getenv("GENSHIN_DATA_REPO")
	if *excelPath != "" {
		repo = *excelPath
	}
	if repo == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			check(err)
		}
		repo = filepath.Join(home, "git", "GenshinData")
	}
	err := excel.LoadResources(func(name string, v any) error {
		d, err := os.ReadFile(filepath.Join(repo, name))
		if err != nil {
			return err
		}
		return json.Unmarshal(d, v)
	})
	check(err)

	var prio []good.CharacterKey
	for _, s := range flag.Args() {
		char, err := good.CharacterKeyString(s)
		if err != nil {
			fmt.Println(priority)
			panic(err)
		}
		prio = append(prio, char)
	}
	if len(prio) == 0 {
		prio = priority
	}

	db := ParseGODatebase(*path)
	for char, t := range config {
		t.Datebase = db

		for _, v := range db.Characters {
			if v.Key == char {
				t.Character = v
			}
		}
		for _, v := range db.Weapons {
			if v.Location == char {
				t.Weapon = v
			}
		}
		for _, v := range db.Artifacts {
			if v.Location == char {
				t.CurArtifacts[v.SlotKey] = v
			}
		}
	}

	for _, v := range db.Artifacts {
		v.Location = good.UnknownCharacterKey
		err := v.Process()
		check(err)
	}

	var total time.Duration
	for _, char := range prio {
		now := time.Now()

		cfg := config[char]
		cur := cfg.CurrentState()
		opt := cfg.Optimize()
		optimized[char] = opt

		dur := time.Since(now)
		total += dur

		change := false
		var a, b strings.Builder
		for i, x := range cur.Build {
			y := opt.Build[i]
			if y == nil {
				log.Fatalf("no build found for %s", cfg.Character.Key)
			}
			if x != y {
				change = true
				slot := strings.ToUpper(string(y.SlotKey.String()[0]))
				if x != nil {
					fmt.Fprintf(&a, "- [%s] %s\n", slot, x)
				}
				fmt.Fprintf(&b, "+ [%s] %s\n", slot, y)
			}
		}
		for _, a := range opt.Build {
			a.Location = char
		}

		result := fmt.Sprintf("%.0f->%.0f (%0.1f%%)", cur.Result, opt.Result, 100*(opt.Result-cur.Result)/cur.Result)
		if int(opt.Result) > int(cur.Result) {
			result = G(result)
		} else if int(opt.Result) < int(cur.Result) {
			result = R(result)
		}

		fmt.Printf("> %s / %s @ %s ; %.2f\n", cfg.Character.Key, cfg.Weapon.Key, result, dur.Seconds())
		if change {
			fmt.Print(R(a.String()))
			fmt.Print(G(b.String()))
			fmt.Println(opt)
		}

		// export build as gcsim config
		goodToSRL := map[string]string{
			"hp":            "hp",
			"hp_":           "hp%",
			"atk":           "atk",
			"atk_":          "atk%",
			"def":           "def",
			"def_":          "def%",
			"eleMas":        "em",
			"enerRech_":     "er",
			"heal_":         "heal",
			"critRate_":     "cr",
			"critDMG_":      "cd",
			"physical_dmg_": "phys%",
			"anemo_dmg_":    "anemo%",
			"geo_dmg_":      "geo%",
			"electro_dmg_":  "electro%",
			"hydro_dmg_":    "hydro%",
			"pyro_dmg_":     "pyro%",
			"cryo_dmg_":     "cryo%",
			"dendro_dmg_":   "dendro%",
		}
		ascToLevel := map[uint32]uint32{
			0: 20,
			1: 40,
			2: 50,
			3: 60,
			4: 70,
			5: 80,
			6: 90,
		}

		setHist := make(map[good.ArtifactSetKey]int)
		subs := good.Stats{}
		for _, a := range opt.Build {
			setHist[a.SetKey]++
			subs.Merge(a.Subs())
		}

		var conf strings.Builder
		fmt.Fprintf(&conf, "%s char lvl=%d/%d cons=%d talent=%d,%d,%d;\n",
			strings.ToLower(cfg.Character.Key.String()),
			cfg.Character.Level,
			ascToLevel[cfg.Character.Ascension],
			cfg.Character.Constellation,
			cfg.Character.Talent.Auto,
			cfg.Character.Talent.Skill,
			cfg.Character.Talent.Burst,
		)
		fmt.Fprintf(&conf, "%s add weapon=\"%s\" lvl=%d/%d refine=%d;\n",
			strings.ToLower(cfg.Character.Key.String()),
			strings.ToLower(cfg.Weapon.Key.String()),
			cfg.Weapon.Level,
			ascToLevel[cfg.Weapon.Ascension],
			cfg.Weapon.Refinement,
		)

		for k, count := range setHist {
			if count < 2 {
				continue
			}
			fmt.Fprintf(&conf, "%s add set=\"%s\" count=%d;\n",
				strings.ToLower(cfg.Character.Key.String()),
				strings.ToLower(k.String()),
				count,
			)
		}

		fmt.Fprintf(&conf, "%s add stats", strings.ToLower(cfg.Character.Key.String()))
		for _, a := range opt.Build {
			fmt.Fprintf(&conf, " %v=%v", goodToSRL[a.MainStatKey.String()], a.Stats.Get(a.MainStatKey))
		}
		fmt.Fprintf(&conf, ";\n")

		fmt.Fprintf(&conf, "%s add stats", strings.ToLower(cfg.Character.Key.String()))
		for k := good.UnknownStatKey + 1; k < good.EndStatType; k++ {
			if subs[k] <= 0 {
				continue
			}
			fmt.Fprintf(&conf, " %v=%v", goodToSRL[good.StatKey(k).String()], subs[k])
		}
		fmt.Fprintf(&conf, ";\n\n")

		_ = os.WriteFile(fmt.Sprintf("./gcsim/%s", strings.ToLower(cfg.Character.Key.String())), []byte(conf.String()), 0o644)

	}
	fmt.Printf("total %.2f\n", total.Seconds())
}

func R(s string) string {
	return "\x1b[31m" + s + "\x1b[1;0m"
}

func G(s string) string {
	return "\x1b[32m" + s + "\x1b[1;0m"
}

func ParseGODatebase(filename string) *good.Datebase {
	data, err := os.ReadFile(filename)
	check(err)

	data, err = handleTraveler(data)
	check(err)

	var db good.Datebase
	err = json.Unmarshal(data, &db)
	check(err)

	return &db
}

func handleTraveler(data []byte) ([]byte, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	var chars []map[string]json.RawMessage
	if err := json.Unmarshal(obj["characters"], &chars); err != nil {
		return nil, err
	}

	for _, c := range chars {
		if !strings.HasPrefix(string(c["key"]), "\"Traveler") {
			continue
		}
		c["key"] = json.RawMessage("\"Traveler\"")
	}

	var err error
	obj["characters"], err = json.Marshal(chars)
	if err != nil {
		return nil, err
	}

	return json.Marshal(obj)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
