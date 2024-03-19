package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shizukayuki/ysoptimizer/pkg/good"
)

// export build as gcsim config
var goodToSRL = map[string]string{
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
var ascToLevel = map[uint32]uint32{
	0: 20,
	1: 40,
	2: 50,
	3: 60,
	4: 70,
	5: 80,
	6: 90,
}

func writeConfigToFile(cfg *OptimizeTarget, opt *OptimizeState) {
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
