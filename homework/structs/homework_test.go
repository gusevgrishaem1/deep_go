package main

import (
	"encoding/json"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i := range name {
			person.PersonName[i] = name[i]
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.CoordX = int32(x)
		person.CoordY = int32(y)
		person.CoordZ = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonGold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.TypeAndFlagsAndMana |= uint16(mana)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b00000100) << 8
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b00001000) << 8
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b00010000) << 8
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		switch personType {
		case 0:
			person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b00000000) << 8
		case 1:
			person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b00100000) << 8
		case 2:
			person.TypeAndFlagsAndMana |= ((person.TypeAndFlagsAndMana >> 8) | 0b01000000) << 8
		}
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonHealth = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectAndStrength |= byte(respect) & 0b00001111
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.RespectAndStrength |= (byte(strength) & 0b00001111) << 4
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.ExperienceAndLevel |= byte(experience) & 0b00001111
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.ExperienceAndLevel |= (byte(level) & 0b00001111) << 4
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	PersonName          [42]byte `json:"PersonName"`
	RespectAndStrength  byte     `json:"RespectAndStrength"`
	ExperienceAndLevel  byte     `json:"ExperienceAndLevel"`
	CoordX              int32    `json:"CoordX"`
	CoordY              int32    `json:"CoordY"`
	CoordZ              int32    `json:"CoordZ"`
	PersonGold          uint32   `json:"PersonGold"`
	TypeAndFlagsAndMana uint16   `json:"TypeAndFlagsAndMana"`
	PersonHealth        uint16   `json:"PersonHealth"`
}

func NewGamePerson(options ...Option) GamePerson {
	var p GamePerson
	for _, opt := range options {
		opt(&p)
	}
	return p
}

func (p *GamePerson) Name() string {
	return unsafe.String(unsafe.SliceData(p.PersonName[:]), len(p.PersonName))
}

func (p *GamePerson) X() int {
	return int(p.CoordX)
}

func (p *GamePerson) Y() int {
	return int(p.CoordY)
}

func (p *GamePerson) Z() int {
	return int(p.CoordZ)
}

func (p *GamePerson) Gold() int {
	return int(p.PersonGold)
}

func (p *GamePerson) Mana() int {
	return int(p.TypeAndFlagsAndMana & 0b0000001111111111)
}

func (p *GamePerson) HasHouse() bool {
	return int(((p.TypeAndFlagsAndMana>>8)&0b00000100)>>2) == 1
}

func (p *GamePerson) HasGun() bool {
	return int(((p.TypeAndFlagsAndMana>>8)&0b00001000)>>3) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return int(((p.TypeAndFlagsAndMana>>8)&0b00010000)>>4) == 1
}

func (p *GamePerson) Type() int {
	return int(p.TypeAndFlagsAndMana >> 13)
}

func (p *GamePerson) Health() int {
	return int(p.PersonHealth)
}

func (p *GamePerson) Respect() int {
	return int(p.RespectAndStrength & 0b00001111)
}

func (p *GamePerson) Strength() int {
	return int((p.RespectAndStrength & 0b11110000) >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.ExperienceAndLevel & 0b00001111)
}

func (p *GamePerson) Level() int {
	return int((p.ExperienceAndLevel & 0b11110000) >> 4)
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
