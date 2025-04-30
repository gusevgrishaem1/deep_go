package main

import (
	"encoding/base64"
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
			person.data[i] = name[i]
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		idx := 42

		for _, b := range Get4Bytes(x) {
			person.data[idx] = b
			idx++
		}

		for _, b := range Get4Bytes(y) {
			person.data[idx] = b
			idx++
		}

		for _, b := range Get4Bytes(z) {
			person.data[idx] = b
			idx++
		}
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		idx := 54

		for _, b := range Get4Bytes(gold) {
			person.data[idx] = b
			idx++
		}
	}
}

func Get4Bytes(i int) []byte {
	return []byte{
		byte(i),
		byte(i >> 8),
		byte(i >> 16),
		byte(i >> 24),
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		idx := 58

		for _, b := range Get10Bits(mana) {
			person.data[idx] = b
			idx++
		}
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		idx := 60

		for _, b := range Get10Bits(health) {
			person.data[idx] = b
			idx++
		}
	}
}

func Get10Bits(i int) []byte {
	return []byte{
		byte(i),
		byte(i>>8) | 0b00000011,
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[62] |= byte(respect) & 0b00001111
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[62] |= (byte(strength) & 0b00001111) << 4
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[63] |= byte(experience) & 0b00001111
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[63] |= (byte(level) & 0b00001111) << 4
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[59] |= 0b00000100
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[59] |= 0b00001000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data[59] |= 0b00010000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		switch personType {
		case 0:
			person.data[59] |= 0b00000000
		case 1:
			person.data[59] |= 0b00100000
		case 2:
			person.data[59] |= 0b01000000
		}
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	data [64]byte
}

func (d GamePerson) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(d.data[:])

	return json.Marshal(encoded)
}

func (d *GamePerson) UnmarshalJSON(b []byte) error {
	var encoded string
	if err := json.Unmarshal(b, &encoded); err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	copy(d.data[:], decoded)
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
	return unsafe.String(unsafe.SliceData(p.data[:42]), 42)
}

func (p *GamePerson) X() int {
	return int(*(*int32)(unsafe.Add(unsafe.Pointer(&p.data), 42)))
}

func (p *GamePerson) Y() int {
	return int(*(*int32)(unsafe.Add(unsafe.Pointer(&p.data), 46)))
}

func (p *GamePerson) Z() int {
	return int(*(*int32)(unsafe.Add(unsafe.Pointer(&p.data), 50)))
}

func (p *GamePerson) Gold() int {
	return int(*(*int32)(unsafe.Add(unsafe.Pointer(&p.data), 54)))
}

func (p *GamePerson) Mana() int {
	return int((*(*int16)(unsafe.Add(unsafe.Pointer(&p.data), 58))) & 0b0000001111111111)
}

func (p *GamePerson) Health() int {
	return int((*(*int16)(unsafe.Add(unsafe.Pointer(&p.data), 60))) & 0b0000001111111111)
}

func (p *GamePerson) Respect() int {
	return int(p.data[62] & 0b00001111)
}

func (p *GamePerson) Strength() int {
	return int((p.data[62] & 0b11110000) >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.data[63] & 0b00001111)
}

func (p *GamePerson) Level() int {
	return int((p.data[63] & 0b11110000) >> 4)
}

func (p *GamePerson) HasHouse() bool {
	return int((p.data[59]&0b00000100)>>2) == 1
}

func (p *GamePerson) HasGun() bool {
	return int((p.data[59]&0b00001000)>>3) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return int((p.data[59]&0b00010000)>>4) == 1
}

func (p *GamePerson) Type() int {
	return int(p.data[59] >> 5)
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
