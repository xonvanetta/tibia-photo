package tibia

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodePhotoType(t *testing.T) {
	tests := []struct {
		input    string
		expected photoType
	}{
		{
			input:    LevelUp.String(),
			expected: LevelUp,
		},
		{
			input:    SkillUp.String(),
			expected: SkillUp,
		},
		{
			input:    LowHealth.String(),
			expected: LowHealth,
		},
		{
			input:    DeathPvE.String(),
			expected: DeathPvE,
		},
		{
			input:    DeathPvP.String(),
			expected: DeathPvP,
		},
		{
			input:    HighestDamageDealt.String(),
			expected: HighestDamageDealt,
		},
		{
			input:    HighestHealingDone.String(),
			expected: HighestHealingDone,
		},
		{
			input:    Hotkey.String(),
			expected: Hotkey,
		},
		{
			input:    GiftOfLife.String(),
			expected: GiftOfLife,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := decodePhotoType(test.input)
			assert.True(t, actual.IsValid())
			assert.Equal(t, test.expected, actual)
		})
	}

	invalid := decodePhotoType("invalid")
	assert.False(t, invalid.IsValid())
	assert.Equal(t, "invalid", invalid.String())
}
