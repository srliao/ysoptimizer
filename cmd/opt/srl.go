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
	fischl()
	bennett(true)
	ningguang(true)

	//hp scalers
	kokomi()
	nilou()
	zhongli()
	layla()

	//supports
	kazuha()
	venti()
	jean()
	// don't have enough vv pieces to go around rip
	sucrose()
	shenhe(true) //also force NO cause she's used for others too sometimes

	itto()
	noelle()
	gorou()

	//todo list (no particular order)
	// - itto
	// - noelle

	// - keqing
	// - tighnari
	// - layla
	// - mona
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
			if mustNO && s.SetBonus != good.NoblesseOblige {
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

func shenhe(mustNO bool) {
	//ill just use whatever shizuka uses here; optimize ayaka's damage instead
	priority = append(priority, good.Shenhe)
	config[good.Shenhe] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ER, good.ATKP).
			Goblet(good.ATKP).
			Circlet(good.ATKP).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			if mustNO && s.SetBonus != good.NoblesseOblige {
				return false
			}
			//skip vv i already don't have enough
			for _, a := range s.Build {
				if a.SetKey == good.ViridescentVenerer {
					return false
				}
			}
			switch s.SetBonus {
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			return s.Get(good.ER) >= 2.20
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			ayaka, ok := optimized[good.KamisatoAyaka]
			if !ok {
				return 0
			}

			// A1: Deific Embrace
			ayaka.Add(good.CryoP, .15)
			// A4: Spirit Communion Seal
			ayaka.BurstDMG += .10

			if s.SetBonus == good.NoblesseOblige {
				ayaka.Add(good.ATKP, .20)
			}

			// Icy Quill
			flatdmg := s.TotalATK() * .776

			// Cutting DMG
			dmg := (ayaka.TotalATK() * 1.909) + flatdmg
			dmg *= 1 + ayaka.AllDMG + ayaka.BurstDMG + ayaka.Get(good.CryoP)
			dmg *= ayaka.CritAverage(0, 0)
			return dmg
		},
	}
}

func zhongli() {
	priority = append(priority, good.Zhongli)
	config[good.Zhongli] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HPP).
			Circlet(good.HPP).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			count := 0
			for _, v := range s.Build {
				if v.SetKey == good.TenacityOfTheMillelith {
					count++
				}
			}
			return count <= 4
			// return s.SetBonus == good.TenacityOfTheMillelith
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.TotalHP()
		},
	}
}

func layla() {
	priority = append(priority, good.Layla)
	config[good.Layla] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HPP).
			Circlet(good.HPP).
			Skip(good.DEFP).Max(1).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.ER) >= 1.60 && s.SetBonus == good.TenacityOfTheMillelith
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.TotalHP()
		},
	}
}

func nilou() {
	priority = append(priority, good.Nilou)
	config[good.Nilou] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HPP).
			Circlet(good.HPP).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			// A1: Court of Dancing Petals
			s.Add(good.EM, 100)
			// Elemental Resonance: Soothing Water
			s.Add(good.HPP, .25)
			switch t.Weapon.Key {
			case good.KeyOfKhajNisut:
				hp := s.TotalHP()
				s.Add(good.EM, hp*.0012*3)
				s.Add(good.EM, hp*.002)
			case good.IronSting:
				//actually all damage but this is fine
				s.Add(good.DendroP, 0.06*2) //rank 1
			}
			count := 0
			for _, v := range s.Build {
				if v.SetKey == good.TenacityOfTheMillelith {
					count++
				}
			}
			return count <= 2
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			// A4: Dreamy Dance of Aeons
			a4 := max(0, s.TotalHP()-30000)
			a4 = min(a4*.001*.09, 4)

			// Bloom DMG
			em := s.Get(good.EM)
			bloom := 1446.9 * (1 + (16*em)/(2000+em) + a4) * 2
			bloom *= enemyMult(0, 0, 1)
			return bloom
		},
	}
}

func kokomi() {
	priority = append(priority, good.SangonomiyaKokomi)
	config[good.SangonomiyaKokomi] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.HPP).
			Goblet(good.HPP).
			Circlet(good.HPP).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			s.Add(good.CR, -1)
			s.Add(good.Heal, .25)
			return s.SetBonus == good.FlowerOfParadiseLost
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			// Bloom DMG
			em := s.Get(good.EM)
			bloom := 1446.9 * (1 + (16*em)/(2000+em)) * 2
			bloom *= enemyMult(0, 0, 1)
			return bloom
		},
	}
}

func fischl() {
	priority = append(priority, good.Fischl)
	config[good.Fischl] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.ElectroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP, good.ER).
			Max(2).SlotMax(3, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.GoldenTroupe:
				s.SkillDMG += .25
			}
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			em := s.Get(good.EM)
			agg := 1446.9 * (1 + (5*em)/(1200+em)) * 1.15

			// Oz's ATK DMG
			dmg := s.TotalATK() * 1.776
			dmg = dmg*2 + (dmg + agg)
			dmg *= 1 + s.AllDMG + s.SkillDMG + s.Get(good.ElectroP)
			dmg *= s.CritAverage(s.SkillCR, 0)
			return dmg
		},
	}
}

func ningguang(mustNO bool) {
	priority = append(priority, good.Ningguang)
	config[good.Ningguang] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.GeoP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP).
			Max(2).SlotMax(3, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			if mustNO && s.SetBonus != good.NoblesseOblige {
				return false
			}
			switch s.SetBonus {
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.GeoP)
			dmg *= s.CritAverage(0, 0)
			dmg *= 1.4783200025558472 //t9 single shard
			return dmg
		},
	}
}

func bennett(mustNO bool) {
	priority = append(priority, good.Bennett)
	config[good.Bennett] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.PyroP, good.ATKP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.DEFP).
			Max(2).SlotMax(3, good.Sands, good.Goblet, good.Circlet).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			if mustNO && s.SetBonus != good.NoblesseOblige {
				return false
			}
			switch t.Weapon.Key {
			case good.SkyriderSword:
				s.Add(good.CR, 0.04)
			}
			switch s.SetBonus {
			case good.NoblesseOblige:
				s.Add(good.ATKP, .20)
			}
			return true
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.BurstDMG + s.Get(good.PyroP)
			dmg *= s.CritAverage(0, 0)
			dmg *= 2.339200019836426 //t9 skill press
			return dmg
		},
	}
}

func venti() {
	priority = append(priority, good.Venti)
	config[good.Venti] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.EM, good.ER).
			Goblet(good.EM).
			Circlet(good.EM, good.CR).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.ER) >= 2.50 && s.Get(good.CR) >= 0.3 && s.SetBonus == good.ViridescentVenerer
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.Get(good.EM)
		},
	}
}

func kazuha() {
	priority = append(priority, good.KaedeharaKazuha)
	config[good.KaedeharaKazuha] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.EM, good.ER).
			Goblet(good.EM).
			Circlet(good.EM).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.ER) >= 1.70 && s.Get(good.CR) >= 0.3 && s.SetBonus == good.ViridescentVenerer
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.Get(good.EM)
		},
	}
}

func sucrose() {
	priority = append(priority, good.Sucrose)
	config[good.Sucrose] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.EM, good.ER, good.ATKP).
			Goblet(good.EM, good.ATKP, good.AnemoP).
			Circlet(good.EM, good.CR, good.ATKP).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.ER) >= 1.30 && s.SetBonus == good.ViridescentVenerer
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.Get(good.EM)
		},
	}
}

func jean() {
	priority = append(priority, good.Jean)
	config[good.Jean] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP).
			Goblet(good.AnemoP).
			Circlet(good.CR, good.CD).
			Skip(good.DEFP).Max(1).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.ER) >= 1.60 && s.SetBonus == good.ViridescentVenerer
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			dmg := s.TotalATK()
			dmg *= 1 + s.AllDMG + s.SkillDMG + s.Get(good.AnemoP)
			dmg *= s.CritAverage(s.SkillCR, 0)
			// Skill DMG
			dmg *= 4.964
			return dmg
		},
	}
}

func itto() {
	priority = append(priority, good.AratakiItto)
	config[good.AratakiItto] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.DEFP, good.ER).
			Goblet(good.GeoP, good.DEFP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.EM).Max(1).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.GladiatorsFinale:
				s.NormalDMG += .35
			case good.HuskOfOpulentDreams:
				s.Add(good.GeoP, .06*4)
				s.Add(good.DEFP, .06*4)
			}
			return s.Get(good.ER) >= 1.30
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			def := s.TotalDEF()

			// Burst: ATK Bonus + C6
			bonusATK := def * (.979)
			bonusDmg := def * 0.35 //a4
			// 1-Hit DMG
			dmg := (s.TotalATK() + bonusATK) * 1.675

			if t.Weapon.Key == good.RedhornStonethresher {
				dmg += def * .40
			}

			dmg *= 1 + s.AllDMG + s.NormalDMG + s.Get(good.GeoP) + bonusDmg
			dmg *= s.CritAverage(0, 0)
			return dmg
		},
	}
}

func noelle() {
	priority = append(priority, good.Noelle)
	config[good.Noelle] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.ATKP, good.DEFP).
			Goblet(good.GeoP, good.DEFP).
			Circlet(good.CR, good.CD).
			Skip(good.HPP, good.EM).Max(1).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			switch s.SetBonus {
			case good.GladiatorsFinale:
				s.NormalDMG += .35
			case good.HuskOfOpulentDreams:
				s.Add(good.GeoP, .06*4)
				s.Add(good.DEFP, .06*4)
			}
			return s.Get(good.ER) >= 1.10
		},
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			def := s.TotalDEF()

			// Burst: ATK Bonus + C6
			bonusATK := def * (.85 + .50)
			// 1-Hit DMG
			dmg := (s.TotalATK() + bonusATK) * 1.564

			if t.Weapon.Key == good.RedhornStonethresher {
				dmg += def * .40
			}

			dmg *= 1 + s.AllDMG + s.NormalDMG + s.Get(good.GeoP)
			dmg *= s.CritAverage(0, 0)
			return dmg
		},
	}
}

func gorou() {
	priority = append(priority, good.Gorou)
	config[good.Gorou] = &OptimizeTarget{
		Filter: NewFilter().
			Sands(good.DEFP, good.ER).
			Goblet(good.DEFP, good.GeoP, good.ATKP).
			Circlet(good.DEFP, good.CR).
			Build(),
		Buffs: func(t *OptimizeTarget, s *OptimizeState) bool {
			return s.Get(good.CR) > 0.3 && s.Get(good.ER) >= 1.5
		},
		IgnoreEnemy: true,
		Target: func(t *OptimizeTarget, s *OptimizeState) float32 {
			return s.TotalDEF()
		},
	}
}
