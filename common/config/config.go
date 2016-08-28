/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

// Package Config contains data structures that will hold the config data for the server
// and the utilities to serialize them
package config

import (
	"github.com/knowlet/kagami/common/consts"
	"github.com/knowlet/maplelib"
)

// A PacketEncodable is a generic interface for all data structures that
// are encodable in a packet
type PacketEncodable interface {
	// Encode serializes the object's data into a maplestory packet
	Encode(dst *maplelib.Packet)
}

// Ey b0ss, a B0ss holds the configuration for a single b0ss in a single world. B0ss.
type B0ss struct {
	attempts   int16
	channelIds []byte
}

// Ey b0ss, Encode serializes a B0ss object to a inter-server packet. B0ss.
func (b *B0ss) Encode(dst *maplelib.Packet) {
	dst.Encode2(uint16(b.attempts))
	dst.EncodeBuffer(b.channelIds)
}

// Ey b0ss, DecodeB0ss deserializes a B0ss object from a packet iterator. B0ss.
func DecodeB0ss(it *maplelib.PacketIterator) (b0ss *B0ss, err error) {
	b0ss = &B0ss{}
	tmp, err := it.Decode2()
	b0ss.attempts = int16(tmp)
	b0ss.channelIds, err = it.DecodeBuffer()
	return
}

// Ey b0ss, DefaultB0sses returns the default b0sses from consts.go. B0ss.
func DefaultB0sses() (pianus, pap, zakum, horntail [consts.WorldCount]*B0ss) {
	for i := 0; i < consts.WorldCount; i++ {
		pianus[i] = &B0ss{
			attempts:   consts.WorldMaxPianusAttempts[i],
			channelIds: consts.WorldPianusChannels[i],
		}
		pap[i] = &B0ss{
			attempts:   consts.WorldMaxPapAttempts[i],
			channelIds: consts.WorldPapChannels[i],
		}
		zakum[i] = &B0ss{
			attempts:   consts.WorldMaxZakumAttempts[i],
			channelIds: consts.WorldZakumChannels[i],
		}
		horntail[i] = &B0ss{
			attempts:   consts.WorldMaxHorntailAttempts[i],
			channelIds: consts.WorldHorntailChannels[i],
		}
	}
	return
}

// Attempts returns the maximum number of daily attempts for the boss, -1 = unlimited
func (b *B0ss) Attempts() int16 {
	return b.attempts
}

// ChannelIds returns a list of the channels where this boss spawns
// A single-element array containing 0xFF means that it will spawn on all channels
func (b *B0ss) ChannelIds() []byte {
	return b.channelIds
}

// SetAttempts sets the maximum number of daily attempts for the boss, -1 = unlimited
func (b *B0ss) SetAttempts(attempts int16) {
	b.attempts = attempts
}

// SetChannelIds sets the list of the channels where this boss spawns
// A single-element array containing 0xFF means that it will spawn on all channels
func (b *B0ss) SetChannelIds(channelIds []byte) {
	b.channelIds = channelIds
}

// A Rates object holds the configuration for exp / drop rates for a single world
type Rates struct {
	mobExp   int32
	questExp int32
	mobMeso  int32
	mobDrop  int32
}

// Encode serializes a Rates object to a inter-server packet
func (r *Rates) Encode(dst *maplelib.Packet) {
	dst.Encode4(uint32(r.mobExp))
	dst.Encode4(uint32(r.questExp))
	dst.Encode4(uint32(r.mobMeso))
	dst.Encode4(uint32(r.mobDrop))
}

// DecodeRates deserializes a Rates object from a packet iterator
func DecodeRates(it *maplelib.PacketIterator) (rates *Rates, err error) {
	rates = &Rates{}
	tmp1, err := it.Decode4()
	tmp2, err := it.Decode4()
	tmp3, err := it.Decode4()
	tmp4, err := it.Decode4()
	if err != nil {
		return
	}

	rates.mobExp = int32(tmp1)
	rates.questExp = int32(tmp2)
	rates.mobMeso = int32(tmp3)
	rates.mobDrop = int32(tmp4)
	return
}

// DefaultRates returns the default rates from consts.go
func DefaultRates() (rates [consts.WorldCount]*Rates) {
	for i := 0; i < consts.WorldCount; i++ {
		rates[i] = &Rates{
			mobExp:   consts.WorldMobExp[i],
			questExp: consts.WorldQuestExp[i],
			mobMeso:  consts.WorldMeso[i],
			mobDrop:  consts.WorldDrop[i],
		}
	}
	return
}

func (r *Rates) MobExp() int32   { return r.mobExp }
func (r *Rates) QuestExp() int32 { return r.questExp }
func (r *Rates) MobMeso() int32  { return r.mobMeso }
func (r *Rates) MobDrop() int32  { return r.mobDrop }

func (r *Rates) SetMobExp(mobExp int32)     { r.mobExp = mobExp }
func (r *Rates) SetQuestExp(questExp int32) { r.questExp = questExp }
func (r *Rates) SetMobMeso(mobMeso int32)   { r.mobMeso = mobMeso }
func (r *Rates) SetMobDrop(mobDrop int32)   { r.mobDrop = mobDrop }

// A WorldConf holds the configuration for a single world
type WorldConf struct {
	defaultGmChatMode   bool
	ribbon              byte
	maxMultiLevel       byte
	defaultStorageSlots byte
	maxStat             uint16
	defaultCharSlots    byte
	maxCharSlots        byte
	maxPlayerLoad       int32
	fameTime            int64
	fameResetTime       int64
	mapUnloadTime       int64
	maxChannels         byte
	eventMsg            string
	scrollingHeader     string
	name                string
	rates               *Rates
	pianus              *B0ss
	pap                 *B0ss
	zakum               *B0ss
	horntail            *B0ss
}

// Encode serializes a WorldConf object to a inter-server packet
func (wc *WorldConf) Encode(dst *maplelib.Packet) {
	var tmp byte = 0

	if wc.defaultGmChatMode {
		tmp = 1
	}

	dst.Encode1(tmp)
	dst.Encode1(wc.ribbon)
	dst.Encode1(wc.maxMultiLevel)
	dst.Encode1(wc.defaultStorageSlots)
	dst.Encode2(wc.maxStat)
	dst.Encode1(wc.defaultCharSlots)
	dst.Encode1(wc.maxCharSlots)
	dst.Encode4(uint32(wc.maxPlayerLoad))
	dst.Encode8(uint64(wc.fameTime))
	dst.Encode8(uint64(wc.fameResetTime))
	dst.Encode8(uint64(wc.mapUnloadTime))
	dst.Encode1(wc.maxChannels)
	dst.EncodeString(wc.eventMsg)
	dst.EncodeString(wc.scrollingHeader)
	dst.EncodeString(wc.name)
	wc.rates.Encode(dst)
	wc.pianus.Encode(dst)
	wc.pap.Encode(dst)
	wc.zakum.Encode(dst)
	wc.horntail.Encode(dst)
}

// DecodeWorldConf deserializes a WorldConf object from an packet iterator
func DecodeWorldConf(it *maplelib.PacketIterator) (wc *WorldConf, err error) {
	wc = &WorldConf{}
	tmp1, err := it.Decode1()
	wc.defaultGmChatMode = tmp1 > 0
	wc.ribbon, err = it.Decode1()
	wc.maxMultiLevel, err = it.Decode1()
	wc.defaultStorageSlots, err = it.Decode1()
	wc.maxStat, err = it.Decode2()
	wc.defaultCharSlots, err = it.Decode1()
	wc.maxCharSlots, err = it.Decode1()
	tmp2, err := it.Decode4()
	tmp3, err := it.Decode8()
	tmp4, err := it.Decode8()
	tmp5, err := it.Decode8()
	wc.maxPlayerLoad = int32(tmp2)
	wc.fameTime = int64(tmp3)
	wc.fameResetTime = int64(tmp4)
	wc.mapUnloadTime = int64(tmp5)
	wc.maxChannels, err = it.Decode1()
	wc.eventMsg, err = it.DecodeString()
	wc.scrollingHeader, err = it.DecodeString()
	wc.name, err = it.DecodeString()
	wc.rates, err = DecodeRates(it)
	wc.pianus, err = DecodeB0ss(it)
	wc.pap, err = DecodeB0ss(it)
	wc.zakum, err = DecodeB0ss(it)
	wc.horntail, err = DecodeB0ss(it)
	return
}

// DefaultWorldConf returns the default world config from consts.go
func DefaultWorldConf() (configs [consts.WorldCount]*WorldConf) {
	defaultrates := DefaultRates()
	defaultpianus, defaultpap, defaultzakum, defaulthorntail := DefaultB0sses()

	for i := 0; i < consts.WorldCount; i++ {
		configs[i] = &WorldConf{
			defaultGmChatMode:   consts.WorldDefaultGMChat[i],
			ribbon:              consts.WorldRibbon[i],
			maxMultiLevel:       consts.WorldMaxMultiLevel[i],
			defaultStorageSlots: consts.WorldDefaultStorageSlots[i],
			maxStat:             consts.WorldMaxStats[i],
			defaultCharSlots:    consts.WorldDefaultCharSlots[i],
			maxCharSlots:        consts.WorldMaxCharSlots[i],
			maxPlayerLoad:       consts.WorldMaxPlayerLoad[i],
			fameTime:            consts.WorldFameDelay[i],
			fameResetTime:       consts.WorldFameResetTime[i],
			mapUnloadTime:       consts.WorldMapUnloadTime[i],
			maxChannels:         consts.WorldChannelCount[i],
			eventMsg:            consts.WorldEventMessage[i],
			scrollingHeader:     consts.WorldScrollingHeader[i],
			name:                consts.WorldName[i],
			rates:               defaultrates[i],
			pianus:              defaultpianus[i],
			pap:                 defaultpap[i],
			zakum:               defaultzakum[i],
			horntail:            defaulthorntail[i],
		}
	}
	return
}

func (wc *WorldConf) DefaultGmChatMode() bool   { return wc.defaultGmChatMode }
func (wc *WorldConf) Ribbon() byte              { return wc.ribbon }
func (wc *WorldConf) MaxMultiLevel() byte       { return wc.maxMultiLevel }
func (wc *WorldConf) DefaultStorageSlots() byte { return wc.defaultStorageSlots }
func (wc *WorldConf) MaxStat() uint16           { return wc.maxStat }
func (wc *WorldConf) DefaultCharSlots() byte    { return wc.defaultCharSlots }
func (wc *WorldConf) MaxCharSlots() byte        { return wc.maxCharSlots }
func (wc *WorldConf) MaxPlayerLoad() int32      { return wc.maxPlayerLoad }
func (wc *WorldConf) FameTime() int64           { return wc.fameTime }
func (wc *WorldConf) FameResetTime() int64      { return wc.fameResetTime }
func (wc *WorldConf) MapUnloadTime() int64      { return wc.mapUnloadTime }
func (wc *WorldConf) MaxChannels() byte         { return wc.maxChannels }
func (wc *WorldConf) EventMsg() string          { return wc.eventMsg }
func (wc *WorldConf) ScrollingHeader() string   { return wc.scrollingHeader }
func (wc *WorldConf) Name() string              { return wc.name }
func (wc *WorldConf) Rates() *Rates             { return wc.rates }
func (wc *WorldConf) Pianus() *B0ss             { return wc.pianus }
func (wc *WorldConf) Pap() *B0ss                { return wc.pap }
func (wc *WorldConf) Zakum() *B0ss              { return wc.zakum }
func (wc *WorldConf) Horntail() *B0ss           { return wc.horntail }

// SetDefaultGmChatMode sets wether the GM chat is enabled by default in the world
func (wc *WorldConf) SetDefaultGmChatMode(defaultGmChatMode bool) {
	wc.defaultGmChatMode = defaultGmChatMode
}

// SetRibbon sets the world's ribbon on world selection 0 = None, 1 = E, 2 = N, 3 = H
func (wc *WorldConf) SetRibbon(ribbon byte)               { wc.ribbon = ribbon }
func (wc *WorldConf) SetMaxMultiLevel(maxMultiLevel byte) { wc.maxMultiLevel = maxMultiLevel }

func (wc *WorldConf) SetDefaultStorageSlots(defaultStorageSlots byte) {
	wc.defaultStorageSlots = defaultStorageSlots
}

func (wc *WorldConf) SetMaxStat(maxStat uint16) { wc.maxStat = maxStat }

func (wc *WorldConf) SetDefaultCharSlots(defaultCharSlots byte) {
	wc.defaultCharSlots = defaultCharSlots
}

func (wc *WorldConf) SetMaxCharSlots(maxCharSlots byte)    { wc.maxCharSlots = maxCharSlots }
func (wc *WorldConf) SetMaxPlayerLoad(maxPlayerLoad int32) { wc.maxPlayerLoad = maxPlayerLoad }

// SetFameTime sets the cooldown in seconds before you can fame someone again
func (wc *WorldConf) SetFameTime(fameTime int64) { wc.fameTime = fameTime }

// SetFameResetTime sets the cooldown in seconds before you can fame the same player again
func (wc *WorldConf) SetFameResetTime(fameResetTime int64)      { wc.fameResetTime = fameResetTime }
func (wc *WorldConf) SetMapUnloadTime(mapUnloadTime int64)      { wc.mapUnloadTime = mapUnloadTime }
func (wc *WorldConf) SetMaxChannels(maxChannels byte)           { wc.maxChannels = maxChannels }
func (wc *WorldConf) SetEventMsg(eventMsg string)               { wc.eventMsg = eventMsg }
func (wc *WorldConf) SetScrollingHeader(scrollingHeader string) { wc.scrollingHeader = scrollingHeader }
func (wc *WorldConf) SetName(name string)                       { wc.name = name }
func (wc *WorldConf) SetRates(rates *Rates)                     { wc.rates = rates }
func (wc *WorldConf) SetPianus(pianus *B0ss)                    { wc.pianus = pianus }
func (wc *WorldConf) SetPap(pap *B0ss)                          { wc.pap = pap }
func (wc *WorldConf) SetZakum(zakum *B0ss)                      { wc.zakum = zakum }
func (wc *WorldConf) SetHorntail(horntail *B0ss)                { wc.horntail = horntail }
