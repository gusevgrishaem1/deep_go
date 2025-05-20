package main

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

const (
	nameSize = 42

	manaMask       = 0x03FF // 10 бит: 00000011 11111111
	flagsShift     = 8
	typeShift      = 13
	houseFlag      = 1 << 2
	gunFlag        = 1 << 3
	familyFlag     = 1 << 4
	typeBuilder    = 0 << typeShift
	typeBlacksmith = 1 << typeShift
	typeWarrior    = 2 << typeShift
	typeMask       = 0xE000 // 3 бит: 11100000 00000000

	// Masks for respectAndStrength / experienceAndLevel
	low4BitsMask   = 0x0F
	high4BitsMask  = 0xF0
	high4BitsShift = 4
)

func WithName(s string) Option {
	return func(p *GamePerson) {
		copy(p.name[:], s)
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(p *GamePerson) {
		p.x, p.y, p.z = int32(x), int32(y), int32(z)
	}
}

func WithGold(gold int) Option {
	return func(p *GamePerson) {
		p.gold = uint32(gold)
	}
}

func WithMana(mana int) Option {
	return func(p *GamePerson) {
		p.typeAndFlagsAndMana |= uint16(mana) & manaMask
	}
}

func WithFlag(flag uint16) Option {
	return func(p *GamePerson) {
		p.typeAndFlagsAndMana |= flag << flagsShift
	}
}

func WithType(personType int) Option {
	return func(p *GamePerson) {
		switch personType {
		case BuilderGamePersonType:
			p.typeAndFlagsAndMana |= typeBuilder
		case BlacksmithGamePersonType:
			p.typeAndFlagsAndMana |= typeBlacksmith
		case WarriorGamePersonType:
			p.typeAndFlagsAndMana |= typeWarrior
		}
	}
}

func WithHealth(health int) Option {
	return func(p *GamePerson) {
		p.personHealth = uint16(health)
	}
}

func WithRespect(respect int) Option {
	return func(p *GamePerson) {
		p.respectAndStrength |= byte(respect) & low4BitsMask
	}
}

func WithStrength(strength int) Option {
	return func(p *GamePerson) {
		p.respectAndStrength |= (byte(strength) & low4BitsMask) << high4BitsShift
	}
}

func WithExperience(exp int) Option {
	return func(p *GamePerson) {
		p.experienceAndLevel |= byte(exp) & low4BitsMask
	}
}

func WithLevel(level int) Option {
	return func(p *GamePerson) {
		p.experienceAndLevel |= (byte(level) & low4BitsMask) << high4BitsShift
	}
}

func WithHouse() Option  { return WithFlag(houseFlag) }
func WithGun() Option    { return WithFlag(gunFlag) }
func WithFamily() Option { return WithFlag(familyFlag) }

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

// GamePerson представляет игрового персонажа с компактным(64 байта) представлением в памяти.
type GamePerson struct {
	name                [nameSize]byte // имя персонажа, зафиксировано на 42 байта
	respectAndStrength  byte           // 4 бита: уважение, 4 бита: сила
	experienceAndLevel  byte           // 4 бита: опыт, 4 бита: уровень
	x, y, z             int32          // координаты в 3D пространстве
	gold                uint32         // количество золота
	typeAndFlagsAndMana uint16         // 3 бита: тип, 3 бит: флаги, 10 бит: мана
	personHealth        uint16         // здоровье персонажа
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
		"hasFamily":  p.HasFamily(),
	}
	return json.Marshal(m)
}

func (p GamePerson) MarshalYAML() (interface{}, error) {
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
		"hasFamily":  p.HasFamily(),
	}
	return m, nil
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

func (p *GamePerson) UnmarshalYAML(value *yaml.Node) error {
	m := map[string]interface{}{}
	if err := value.Decode(&m); err != nil {
		return err
	}

	name := m["name"].(string)
	x := m["x"].(int)
	y := m["y"].(int)
	z := m["z"].(int)
	gold := m["gold"].(int)
	mana := m["mana"].(int)
	health := m["health"].(int)
	respect := m["respect"].(int)
	strength := m["strength"].(int)
	experience := m["experience"].(int)
	level := m["level"].(int)
	personType := m["type"].(int)
	hasHouse := m["hasHouse"].(bool)
	hasGun := m["hasGun"].(bool)
	hasFamily := m["hasFamily"].(bool)

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

func (p *GamePerson) X() int       { return int(p.x) }
func (p *GamePerson) Y() int       { return int(p.y) }
func (p *GamePerson) Z() int       { return int(p.z) }
func (p *GamePerson) Gold() int    { return int(p.gold) }
func (p *GamePerson) Mana() int    { return int(p.typeAndFlagsAndMana & manaMask) }
func (p *GamePerson) Health() int  { return int(p.personHealth) }
func (p *GamePerson) Respect() int { return int(p.respectAndStrength & low4BitsMask) }
func (p *GamePerson) Strength() int {
	return int((p.respectAndStrength & high4BitsMask) >> high4BitsShift)
}
func (p *GamePerson) Experience() int { return int(p.experienceAndLevel & low4BitsMask) }
func (p *GamePerson) Level() int {
	return int((p.experienceAndLevel & high4BitsMask) >> high4BitsShift)
}

func (p *GamePerson) HasHouse() bool  { return (p.typeAndFlagsAndMana>>flagsShift)&houseFlag != 0 }
func (p *GamePerson) HasGun() bool    { return (p.typeAndFlagsAndMana>>flagsShift)&gunFlag != 0 }
func (p *GamePerson) HasFamily() bool { return (p.typeAndFlagsAndMana>>flagsShift)&familyFlag != 0 }
func (p *GamePerson) Type() int       { return int((p.typeAndFlagsAndMana & typeMask) >> typeShift) }

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

	yamlData, err := yaml.Marshal(person)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(yamlData, &person)
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
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
