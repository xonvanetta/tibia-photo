package tibia

type photoType string

const (
	SkillUp            photoType = "SkillUp"
	LevelUp            photoType = "LevelUp"
	LowHealth          photoType = "LowHealth"
	DeathPvE           photoType = "DeathPvE"
	DeathPvP           photoType = "DeathPvP"
	HighestDamageDealt photoType = "HighestDamageDealt"
	HighestHealingDone photoType = "HighestHealingDone"
	Hotkey             photoType = "Hotkey"
	GiftOfLife         photoType = "GiftOfLife"
)

func (pt photoType) Is(types ...photoType) bool {
	for _, p := range types {
		if pt == p {
			return true
		}
	}
	return false
}

func (pt photoType) String() string {
	return string(pt)
}

func (pt photoType) IsValid() bool {
	switch pt {
	case SkillUp, LevelUp, LowHealth, DeathPvE, DeathPvP, HighestDamageDealt, HighestHealingDone, Hotkey, GiftOfLife:
		return true
	}
	return false
}

func decodePhotoType(t string) photoType {
	switch photoType(t) {
	case LevelUp:
		return LevelUp
	case SkillUp:
		return SkillUp
	case LowHealth:
		return LowHealth
	case DeathPvE:
		return DeathPvE
	case DeathPvP:
		return DeathPvP
	case HighestDamageDealt:
		return HighestDamageDealt
	case HighestHealingDone:
		return HighestHealingDone
	case Hotkey:
		return Hotkey
	case GiftOfLife:
		return GiftOfLife
	default:
		return photoType(t)
	}
}
