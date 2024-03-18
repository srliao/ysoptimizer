package main

import "github.com/shizukayuki/ysoptimizer/pkg/good"

func init() {
	//override shizuka's defaults
	priority = []good.CharacterKey{}
	config = make(map[good.CharacterKey]*OptimizeTarget)
	raiden()
	navia()
	furina()
	chiori()
	nahida()
	yae()
	yelan(false)
	ayaka()
	ganyu()
	xiangling()
	xingqiu()

	//todo list (no particular order)
	// - ayaka
	// - ganyu
	// - shenhe
	// - nilou
	// - itto
	// - kazuha
	// - fischl

	// - bennett
	// - jean
	// - keqing
	// - tighnari
	// - kokomi
	// - ningguang
	// - sucrose
	// - layla
	// - venti
	// - mona
	// - noelle
	// - kuki (hb)
	// maybes?? i never use??
	// - beidou
	// - gorou
	// - rosaria
	// - chongyun
	// - sara

}

func raiden() {
	priority = append(priority, good.RaidenShogun)
	config[good.RaidenShogun] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ER, good.ATKP).
			Goblet(good.ElectroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP, good.EM).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.StaffOfTheScarletSands:
				s.Add(good.ATK, s.Get(good.EM)*3*.28)
			case good.EngulfingLightning:
				//+30% er after first
				s.Add(good.ER, 0.3) // always optimize burst dmg anyways
				//28% of er over 100%, to a max of 80%
				s.Add(good.ATK, min((s.Get(good.ER)-1)*0.28, 0.8))
			}
			// Skill: Eye of Stormy Judgment
			s.BurstDMG += .0030 * 90
			// A4: Enlightened One
			s.Add(good.ElectroP, (s.Get(good.ER)-1)*.4)
			return s.Get(good.ER) >= 2.50
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.ElectroP)
			dmg *= s.CritAverage(s.BurstCR, 0)
			// Musou no Hitotachi DMG (40 stacks)
			dmg *= 7.214 + .07*40
			return dmg
		},
	}
}

func navia() {
	priority = append(priority, good.Navia)
	config[good.Navia] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.GeoP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP, good.EM).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.Verdict:
				//surely i can pick up crystals
				s.SkillDMG += .18 * 2
			}
			switch s.SetBonus {
			case good.NighttimeWhispersInTheEchoingWoods:
				s.Add(good.GeoP, .20*(1+1.5))
			}
			// A1: Mutual Assistance Network
			s.Add(good.ATKP, .20*2)
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.SkillDMG + (.15 * 3) + s.Get(good.GeoP)
			dmg *= s.CritAverage(s.SkillCR, 0)
			// Rosula Shardshot Base DMG
			dmg *= 7.106 * 2
			return dmg
		},
	}
}

func furina() {
	priority = append(priority, good.Furina)
	config[good.Furina] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HydroP, good.HPP).
			Circlet(good.CR, good.CD).
			Skip(good.ATKP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Circlet, good.Goblet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.GoldenTroupe:
				s.SkillDMG += .25
			}
			return s.Get(good.ER) >= 1.60
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			hp := s.TotalHP()
			a4 := min(hp*.001*.007, .28)

			dmg := hp
			dmg *= 1 + s.AllDMG + s.SkillDMG + a4 + s.Get(good.HydroP)
			dmg *= s.CritAverage(s.SkillCR, 0)
			// Gentilhomme Usher DMG
			dmg *= .1013
			return dmg
		},
	}
}

func chiori() {
	priority = append(priority, good.Chiori)
	config[good.Chiori] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.DEFP).
			Goblet(good.GeoP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.EM, good.ER).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.UrakuMisugiri:
				s.NormalDMG += .16
				s.SkillDMG += .24
			}
			switch s.SetBonus {
			case good.GoldenTroupe:
				s.SkillDMG += .25
			}
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			// Tamoto DMG
			dmg := s.TotalATK()*1.48 + s.TotalDEF()*1.85
			dmg *= 1 + s.AllDMG + s.SkillDMG + s.Get(good.GeoP)
			if t.Weapon.Key == good.CinnabarSpindle {
				dmg += s.TotalDEF() * 0.8
			}
			dmg *= s.CritAverage(s.SkillCR, 0)
			return dmg
		},
	}
}

func nahida() {
	priority = append(priority, good.Nahida)
	config[good.Nahida] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.EM).
			Goblet(good.DendroP, good.EM).
			Circlet(good.CR, good.CD, good.EM).
			Skip(good.HPP, good.DEFP, good.ER).
			Max(1).SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.AThousandFloatingDreams:
				s.Add(good.DendroP, .10*3)
				// s.Add(good.EM, 32*3)
			}
			switch s.SetBonus {
			case good.GildedDreams:
				s.Add(good.EM, 50*3)
			}
			// A4: Awakening Elucidated
			if em := s.Get(good.EM); em > 200 {
				s.SkillDMG += min((em-200)*.001, .80)
				s.SkillCR += min((em-200)*.0003, .24)
			}
			// skip any vv pieces so not to steal good em piece
			for _, a := range s.Build {
				if a.SetKey == good.ViridescentVenerer {
					return false
				}
			}

			return true
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			em := s.Get(good.EM)
			spread := 1446.9 * (1 + (5*em)/(1200+em)) * 1.25

			dmg := s.TotalATK()*1.858 + em*3.715 + spread
			dmg *= 1 + s.AllDMG + s.SkillDMG + s.Get(good.DendroP)
			dmg *= s.CritAverage(s.SkillCR, 0)

			switch s.SetBonus {
			case good.DeepwoodMemories:
				dmg *= enemyMult(-.30, 0, 0)
			default:
				dmg *= enemyMult(0, 0, 0)
			}
			return dmg
		},
	}
}

func yelan(mustNO bool) {
	priority = append(priority, good.Yelan)
	config[good.Yelan] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HydroP, good.HPP).
			Circlet(good.CR, good.CD).
			Skip(good.ATKP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			if mustNO && s.SetBonus == good.NoblesseOblige {
				return false
			}
			switch t.Weapon.Key {
			case good.AquaSimulacra:
				s.AllDMG += .20
			}
			//nob self buffs
			switch s.SetBonus {
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			// A1: Turn Control
			s.Add(good.HPP, .06)
			// A4: Adapt With Ease
			s.AllDMG += .01 + .035*4
			return s.Get(good.ER) >= 1.40
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalHP()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.HydroP)
			dmg *= s.CritAverage(0, 0)
			// Exquisite Throw DMG
			dmg *= .088
			return dmg
		},
	}
}

func xiangling() {
	priority = append(priority, good.Xiangling)
	config[good.Xiangling] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.ER).
			Goblet(good.PyroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.CrimsonWitchOfFlames:
				s.Add(good.PyroP, .15*.50*1)
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			return s.Get(good.ER) >= 1.80
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			em := s.Get(good.EM)
			amp := 1 + (2.778*em)/(1400+em)
			switch s.SetBonus {
			case good.CrimsonWitchOfFlames:
				amp += .15
			}

			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.PyroP)
			dmg *= s.CritAverage(s.BurstCR, 0)
			// Pyronado DMG
			novape := dmg * 2.24
			vape := novape * amp * 1.5
			dmg = vape + novape
			return dmg
		},
	}
}

func xingqiu() {
	priority = append(priority, good.Xingqiu)
	config[good.Xingqiu] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.ER).
			Goblet(good.HydroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			return s.Get(good.ER) >= 1.70
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.HydroP)
			dmg *= s.CritAverage(s.BurstCR, 0)

			//talent 10 rainsword
			return dmg * 1.031167984008789
		},
	}
}

func yae() {
	priority = append(priority, good.YaeMiko)
	config[good.YaeMiko] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.ER).
			Goblet(good.ElectroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.KagurasVerity:
				s.SkillDMG += .12 * 3
				s.Add(good.ElectroP, .12)
			}
			// return s.Get(good.ER) >= 1.40
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			em := s.Get(good.EM)
			agg := 1446.9 * (1 + (5*em)/(1200+em)) * 1.15

			atk := s.TotalATK()
			electro := 1 + s.AllDMG + s.Get(good.ElectroP)

			// combo
			skill := atk * 1.706
			skill = skill*2 + (skill + agg)
			skill *= electro + s.SkillDMG + em*.0015

			// burst := atk*(4.68+6.009*3) + agg*4
			// burst *= electro + s.BurstDMG

			dmg := skill
			dmg *= s.CritAverage(0, 0)
			return dmg
		},
	}
}

func ayaka() {
	priority = append(priority, good.KamisatoAyaka)
	config[good.KamisatoAyaka] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.CryoP).
			Circlet(good.CR, good.CD, good.ATKP).
			Skip(good.HPP, good.DEFP, good.EM).
			Max(1).SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch t.Weapon.Key {
			case good.MistsplitterReforged:
				s.Add(good.CryoP, .28)
			}
			switch s.SetBonus {
			case good.BlizzardStrayer:
				s.Add(good.CR, .20*2)
			}
			// Elemental Resonance: Shattering Ice
			s.Add(good.CR, .15)
			// return s.Get(good.ER) >= 1.40
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.CryoP)
			dmg *= s.CritAverage(0, 0)
			// Cutting DMG
			dmg *= 1.909
			return dmg
		},
	}
}

func ganyu() {
	priority = append(priority, good.Ganyu)
	config[good.Ganyu] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.CryoP).
			Circlet(good.CR, good.CD, good.ATKP).
			Skip(good.HPP, good.DEFP, good.EM).
			Max(1).SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.BlizzardStrayer:
				s.Add(good.CR, .20*2)
			}
			// A1: Undivided Heart
			// s.Add(good.CR, .20)
			// Elemental Resonance: Shattering Ice
			s.Add(good.CR, .15)
			return s.Get(good.ER) > 1.1
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.Get(good.CryoP)
			dmg *= s.CritAverage(0, 0)
			// Ice Shard DMG
			// dmg *= 1.265
			// aimed shot
			dmg *= 2.303999900817871 + 2.2320001125335693
			return dmg
		},
	}
}

func skeleton() {
	priority = append(priority, good.Yelan)
	config[good.Yelan] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HydroP, good.HPP).
			Circlet(good.CR, good.CD).
			Skip(good.ATKP, good.DEFP).
			Max(1).
			SlotMax(2, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return false
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return 0
		},
	}
}
