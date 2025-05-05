package main

import (
	"encoding/json"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(s string) Option {
	return func(p *GamePerson) {
		copy(p.name[:], s)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndFlagsAndMana |= uint16(mana)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b00000100) << 8
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b00001000) << 8
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b00010000) << 8
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		switch personType {
		case 0:
			person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b00000000) << 8
		case 1:
			person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b00100000) << 8
		case 2:
			person.typeAndFlagsAndMana |= ((person.typeAndFlagsAndMana >> 8) | 0b01000000) << 8
		}
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.personHealth = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectAndStrength |= byte(respect) & 0b00001111
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectAndStrength |= (byte(strength) & 0b00001111) << 4
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experienceAndLevel |= byte(experience) & 0b00001111
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experienceAndLevel |= (byte(level) & 0b00001111) << 4
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name                [42]byte
	respectAndStrength  byte
	experienceAndLevel  byte
	x                   int32
	y                   int32
	z                   int32
	gold                uint32
	typeAndFlagsAndMana uint16
	personHealth        uint16
}

func (p GamePerson) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"name":       p.Name(),
		"respect":    p.Respect(),
		"strength":   p.Strength(),
		"level":      p.Level(),
		"experience": p.Experience(),
		"x":          p.X(),
		"y":          p.Y(),
		"z":          p.Z(),
		"gold":       p.Gold(),
		"type":       p.Type(),
		"health":     p.Health(),
		"mana":       p.Mana(),
		"hasHouse":   p.HasHouse(),
		"hasGun":     p.HasGun(),
		"hasFamily":  p.HasFamilty(),
	}
	return json.Marshal(m)
}

func (p *GamePerson) UnmarshalJSON(data []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	x := int(m["x"].(float64))
	y := int(m["y"].(float64))
	z := int(m["z"].(float64))
	gold := int(m["gold"].(float64))
	mana := int(m["mana"].(float64))
	health := int(m["health"].(float64))
	respect := int(m["respect"].(float64))
	strength := int(m["strength"].(float64))
	experience := int(m["experience"].(float64))
	level := int(m["level"].(float64))
	personType := int(m["type"].(float64))
	hasHouse := m["hasHouse"].(bool)
	hasGun := m["hasGun"].(bool)
	hasFamily := m["hasFamily"].(bool)

	options := []Option{
		WithName(m["name"].(string)),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithType(personType),
	}

	if hasHouse {
		options = append(options, WithHouse())
	}
	if hasGun {
		options = append(options, WithGun())
	}
	if hasFamily {
		options = append(options, WithFamily())
	}

	*p = NewGamePerson(options...)

	return nil
}

func NewGamePerson(options ...Option) GamePerson {
	var p GamePerson
	for _, opt := range options {
		opt(&p)
	}
	return p
}

func (p *GamePerson) Name() string {
	return unsafe.String(unsafe.SliceData(p.name[:]), len(p.name))
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.typeAndFlagsAndMana & 0b0000001111111111)
}

func (p *GamePerson) HasHouse() bool {
	return int(((p.typeAndFlagsAndMana>>8)&0b00000100)>>2) == 1
}

func (p *GamePerson) HasGun() bool {
	return int(((p.typeAndFlagsAndMana>>8)&0b00001000)>>3) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return int(((p.typeAndFlagsAndMana>>8)&0b00010000)>>4) == 1
}

func (p *GamePerson) Type() int {
	return int(p.typeAndFlagsAndMana >> 13)
}

func (p *GamePerson) Health() int {
	return int(p.personHealth)
}

func (p *GamePerson) Respect() int {
	return int(p.respectAndStrength & 0b00001111)
}

func (p *GamePerson) Strength() int {
	return int((p.respectAndStrength & 0b11110000) >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.experienceAndLevel & 0b00001111)
}

func (p *GamePerson) Level() int {
	return int((p.experienceAndLevel & 0b11110000) >> 4)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	jsonData, err := json.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(jsonData, &person)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
