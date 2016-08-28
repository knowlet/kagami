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

// Package packets contains all packet headers GMS v62 and packet builders
// that do not require server-specific packages
package packets

// Send packet headers
const (
	// common
	OPing = 0x0009

	// login server
	OLoginStatus        = 0x0000
	OServerList         = 0x0002
	OCharList           = 0x0003
	OServerIP           = 0x0004
	OCharNameResponse   = 0x0005
	OAddNewCharEntry    = 0x0006
	ODeleteCharResponse = 0x0007
	OChangeChannel 		= 0x0008
	// ping
	OCashShopUse		= 0x000A
	// 0x0B - sub_49D849 [4][1]
	OSecurityClient		= 0x000C
	OChannelSelected	= 0x000D
	// 0x0E
	ORelogResponse      = 0x000F
	OSecondPassError	= 0x0010
	// 0x11 - sub_49D89F [1]
	OChooseGender		= 0x0014
	OGenderSet			= 0x0015
	OServerStatus       = 0x0016

	// TWMS hasn't following packets
	OPinOperation       = 0x999
	OAllCharlist        = 0x999
	OGenderDone         = 0x999
	OPinAssigned        = 0x999

	// CwvsContext

	OModifyInventoryItem	= 0x001B
	OUpdateInventorySlot	= 0x001C
	OUpdateStats   = 0x001D
	OGiveBuff	= 0x001E
	OCancelBuff	= 0x001F
	OTemporaryStats	= 0x0020
	OTemporaryStatsReset = 0x0021
	OUpdateSkills = 0x0022
	OSkillUseResult = 0x0023
	OFameResponse = 0x0024
	OShowStatusInfo = 0x0025
	OShowNotes = 0x0026
	OMapTransferResult = 0x0027
	OAntiMacroResult = 0x0028
	OClaimResult = 0x002A
	// 0x2B
	OClaimStatusChanged = 0x002C
	OSetTamingMobInfo = 0x002D
	OShowQuestCompletion = 0x002E
	OEntrustedShopCheckResult = 0x002F

	OUseSkillBook = 0x0031
	//
	OCharacterInfo = 0x0036
	OPartyOperation = 0x0037
	OBuddyList = 0x0038
	//
	OGuildOperation = 0x003A
	OAllianceOperation = 0x003B
	OSpawnPortal = 0x3C
	OServerMessage = 0x003D
	OIncubatorResult = 0x003E
	OShopScannerResult = 0x003F
	OShopLinkResult = 0x0040
	OMarriageRequest = 0x0041
	OMarriageResponse = 0x0042
	//
	OSetWeekEventMessage = 0x0046
	OPotionDiscountRate = 0x0047
	OBladeMobCatchFail = 0x0048
	//
	OImitatedNpcResult = 0x004A
	OImitatedNpcData = 0x004B
	OLimitedNpcDisableInfo = 0x004C
	OMonsterBookAdd = 0x004D
	OMonsterBookChangeCover = 0x004E
	OHourChanged = 0x004F
	OMiniMapOnOff = 0x0050
	OConsultAuthKeyUpdate = 0x0051
	OClassCompetitionAuthKeyUpdate = 0x0052
	OWebBoardAuthKeyUpdate = 0x0053
	OSessionValue = 0x0054
	OBonusExpChanged = 0x0055
	OSendPedigree = 0x0056
	OOpenFamily = 0x0057
	OFamilyMessage = 0x0058
	OFamilyInvite = 0x0059
	OFamilyJunior = 0x005A
	OFamilySeniorMessage = 0x005B
	OFamily = 0x005C
	OFamilyRepIncrease = 0x005D
	OFamilyLoggedIn = 0x005E
	OFamilyBuff = 0x005F
	OFamilyUseRequest = 0x0060
	OLevelUpdate = 0x0061
	OMarrageUpdate = 0x0062
	OJobUpdate = 0x0063
	OSetBuyEquipExt = 0x0064
	OTopMessage = 0x0065
	ODataCrcCheckFailed = 0x0066
	OShowPredictCard = 0x0067
	OBBSOperation = 0x0068
	OFishingBoardUpdate = 0x0069
	OUpdateBeans = 0x006A
	ODonateBeans = 0x006B
	// 0x6B - sub_A2ED85
	// 0x6C - sub_A4689E
	OAvatarMega = 0x006D
	// 0x6E - sub_A46A49
	// 0x6F - sub_A2F42E
	// 0x70 - CWvsContext__OnFakeGMNotice
	// 0x71 - CWvsContext__OnSuccessInUsegachaponBox
	ONewYearCardResponse = 0x0072
	// 0x73 - sub_00A46A65
	// 0x74 - sub_A46A65 (貌似是變更角色名稱)
	// 0x75 - CWvsContext::OnChangeCharNameError []
	// 0x76 - sub_A47710 (儲存個性文字)
	// 0x77 - CWvsContext__OnExpChairSetToZero
	// 0x78 - CWvsContext__OnExpChairClick
	OSkillMacro = 0x007A
	OSetField = 0x007B
	OSetITC = 0x007C
	OSetCashShop = 0x007D
	//
	OSetMapObjectVisible = 0x007F
	OClearBackEffect = 0x0080
	OMapBlocked = 0x0081
	OServerBlocked = 0x0082
	OShowEquipEffect = 0x0083
	OMultipleChat = 0x0084
	OWhisperChat = 0x0085
	OBossEnvironment = 0x0086
	OUpdateEnvironment = 0x0087
	OCashSong = 0x0089
	OGameMasterEffect = 0x008A
	OOXQuiz = 0x008B
	OGameMasterEventInstructions = 0x008C
	OClock = 0x008D
	OBoatEffect = 0x008E
	OBoatUpdate = 0x008F
	//
	OClockStop = 0x0093
	OAriantScoreBoard = 0x0094

	// Rail Station Packet

	OMovePlateform = 0x0096
	OPyramidResult = 0x0097
	OPyramidUpdate = 0x0098


	OSpawnPlayer = 0x0099
	ORemovePlayer = 0x009A
	OChatText = 0x009B
	OChalkBoard = 0x009C
	OUpdateCharacterBox = 0x009D
	//
	OShowScrollEffect = 0x009E
	OFishingCaught = 0x009F
	OSpawnPet = 0x00A2
	//
	OMovePet = 0x00A5
	OPetChat = 0x00A6
	OPetNameChange = 0x00A7
	OPetLoadExceptionList = 0x00A8
	OPetCommand = 0x00A9
	OSpawnSummon = 0x00AA
	ORemoveSummon = 0x00AB
	OMoveSummon = 0x00AC
	OSummonAttack = 0x00AD
	OSummonSkill = 0x00AE
	ODamageSummon = 0x00AF
	//
	OMovePlayer = 0x00B1
	OCloseRangeAttack = 0x00B2
	ORangedAttack = 0x00B3
	OMagicAttack = 0x00B4
	OEnergyAttack = 0x00B5
	OSkillEffect = 0x00B6
	OCancelSkillEffect = 0x000B7
	ODamagePlayer = 0x00B8
	OFacialExpression = 0x00B9
	OShowItemEffect = 0x00BA
	OShowChair = 0x00BD
	OUpdateCharLook = 0x00BE
	OShowForeignEffect = 0x00BF
	OGiveForeignBuff = 0x00C0
	OCancelForeignBuff = 0x00C1
	OUpdatePartyMemberHP = 0x00C2
	OGuildNameChanged = 0x00C3
	OGuildMarkChanged = 0x00C4
	OCancelChair = 0x00C6
	OShowItemGainInChat = 0x00C7
	OCurrentMapWarp = 0x00C8
	OMesoBagSuccess = 0x00CA
	OMesoBagFailure = 0x00CB
	OUpdateQuestInfo = 0x00CC
	// 0xCD - sub_979BF0 [4] 跟精靈商人有關係
	OPetFlagChange = 0x00CE
	OPlayerHint = 0x00CF
	// 0xD0 - sub_97C071 [%s] 貌似音樂的東西
	ORepairWindow = 0x00D5
	OCygnusIntroductionLock = 0x00D6
	OCygnusIntroductionDisableUI = 0x00D7
	OSummonHint = 0x00D8
	OSummonHintMessage = 0x00D9
	OAranCombo = 0x00DA
	// 0xDB - sub_97038D [%s][4] - 提昇廣播音效(?
	// 0xDC - sub_9803FA
	OGamePollReply = 0x00DD
	// 0xDF - sub_98026E [4] 關於雪球
	// 0xE0 - sub_9802FA 訊息：雪球比賽3
	// 0xE1 - sub_952ABC [4]
	// 0xE2 - sub_980326
	OCoolDown = 0x00E3
	OSpawnMonster = 0x00E5
	OKillMonster = 0x00E6
	OSpawnMonsterControl = 0x00E7
	OMoveMonster = 0x00E8
	OMoveMonsterResponse = 0x00E9
	OApplyMonsterStatus = 0x00EB
	OCancelMonsterStatus = 0x00EC
	OMonsterToMonsterDamage = 0x00EE
	ODamageMonster = 0x00EF
	OShowMonsterHp = 0x00F3
	OShowMagnet = 0x00F4
	OCatchMonster = 0x00F5
	OMonsterSpeaking = 0x00F6
	//
	OSpawnNpc = 0x00F9
	ORemoveNpc = 0x00FA
	OSpawnNpcRequestController = 0x00FB
	ONpcAction = 0x00FC
	//
	OSpawnHiredMerchant = 0x0103
	ODestroyHiredMerchant = 0x0104
	OUpdateHiredMerchant = 0x0106
	ODropItemFromHiredMerchant = 0x0107
	ORemoveItemFromMap = 0x0108
	OSpawnKiteError = 0x0109
	ODestroyKite = 0x0100
	OSpawnKite = 0x010A

	OSpawnMist = 0x010C
	ORemoveMist = 0x010D
	OSpawnDoor = 0x010E
	ORemoveDoor = 0x010F
	OHitReactor = 0x0113
	OSpawnReactor = 0x0115
	ODestroyReactor = 0x0116
	ORollSnowBall = 0x0117
	OHitSnowBall = 0x0118
	OSnowBallMessage = 0x0119
	OLeftNockBack = 0x011A
	OHitCoconut = 0x011B
	OCoconutScore = 0x011C
	// 0x11D - sub_57897A
	// 0x11E - sub_5789A9

	OMonsterCarnivalStart = 0x011F
	OMonsterCarnivalObtainedCP = 0x0120
	OMonsterCarnivalPartyCP = 0x0121
	OMonsterCarnivalSummon = 0x0122
	//
	OMonsterCarnivalDead = 0x0124
	OMonsterCarnivalDisconnect = 0x0125
	OMonsterCarnivalEnd = 0x0126

	OChaosHorntailShrine = 0x0128
	OChaosZakumShrine = 0x0129
	OHornTailShrine = 0x012A
	OZakumShrine = 0x012B
	OEnglishQuiz = 0x012C
	//
	ONpcTalk = 0x013C
	OOpenNpcShop = 0x013D
	OConfirmShopTranscation = 0x013E
	// 0x13B - sub_545511
	// 0x13F - sub_425B76
	// 0x140 - sub_425B76
	OOpenStorage = 0x0141
	OMerchItemMessage = 0x0142
	OMerchItemStore = 0x0143
	ORpsGame = 0x0144
	OMessenger = 0x0145
	OPlayerInteraction = 0x0146
	OBeansGameTips = 0x999
	OBeansGame1 = 0x0153
	OBeansGame2 = 0x0154
	// handle at 005A72DD
	// 0x147 - sub_5A7343
	// 0x148 - sub_5A74AE
	// 0x149 - sub_5A7539
	// 0x14A - sub_5A7711
	// 0x14B - 
	// 0x14C - sub_5AD3DC
	// 0x14D - sub_5AD401
	// handle at 00591938
	// 0x14E - sub_591AF9
	// 0x14F - sub_591B0F
	// 0x152 - sub_58E17A
	// 0x153 - sub_58E0A6
	// 0x154 - sub_58E164
	ODuey = 0x0155
	OCashShopWebSite = 0x0156
	OCashShopUpdate = 0x0157
	OCashShopOperation = 0x0158
	OXmsSurprise = 0x015C
	OPetAutoHP = 0x0164
	OPetAutoMP = 0x0165
	OGetMtsTokens = 0x0169
	OMtsOperation = 0x016A

)

// Recv packet headers
const (
	// common
	IPong = 0x0018

	// login server
	ILoginPassword       = 0x0001
	IAfterLogin          = 0x0009
	IServerListRequest   = 0x000B
	IServerListRerequest = 0x0004
	IServerStatusRequest = 0x0006
	IViewAllChar         = 0x000D
	IRelog               = 0x001C
	ICharlistRequest     = 0x0005
	ICharSelect          = 0x0013
	ICheckCharName       = 0x0015
	ICreateChar          = 0x0016
	IDeleteChar          = 0x0017
	IPickAllChar         = 0x000E // unused, idk what this does
	ISetGender           = 0x0008
	IRegisterPin         = 0x000A
	IGuestLogin          = 0x0002
	IUnknownPlsIgnore1   = 0x001A // this gets spammed while on login screen, apparently it means client error
	IUnknownPlsIgnore2   = 0x000F
	IPlayerUpdateIgnore  = 0x00C0 // shouldn't be received by the login server

	// channel server
	ILoadCharacter    = 0x0014
	IPlayerUpdate     = 0x00C0
	IChangeMapSpecial = 0x005C
	IChangeMap        = 0x0023
	IMovePlayer       = 0x0026
)
