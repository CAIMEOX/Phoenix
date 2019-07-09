package packet

const (
	IDLogin = iota + 0x01
	IDPlayStatus
	IDServerToClientHandshake
	IDClientToServerHandshake
	IDDisconnect
	IDResourcePacksInfo
	IDResourcePackStack
	IDResourcePackClientResponse
	IDText
	IDSetTime
	IDStartGame
	IDAddPlayer
	IDAddEntity
	IDRemoveEntity
	IDAddItemEntity
	_
	IDTakeItemEntity
	IDMoveEntityAbsolute
	IDMovePlayer
	IDRiderJump
	IDUpdateBlock
	IDAddPainting
	IDExplode
	_ // IDLevelSoundEvent(1): We don't bother implementing this.
	IDLevelEvent
	IDBlockEvent
	IDEntityEvent
	IDMobEffect
	IDUpdateAttributes
	IDInventoryTransaction
	IDMobEquipment
	IDMobArmourEquipment
	IDInteract
	IDBlockPickRequest
	IDEntityPickRequest
	IDPlayerAction
	IDEntityFall
	IDHurtArmour
	IDSetEntityData
	IDSetEntityMotion
	IDSetEntityLink
	IDSetHealth
	IDSetSpawnPosition
	IDAnimate
	IDRespawn
	IDContainerOpen
	IDContainerClose
	IDPlayerHotBar
	IDInventoryContent
	IDInventorySlot
	IDContainerSetData
	IDCraftingData
	IDCraftingEvent
	IDGUIDataPickItem
	IDAdventureSettings
	IDBlockEntityData
	IDPlayerInput
	IDFullChunkData
	IDSetCommandsEnabled
	IDSetDifficulty
	IDChangeDimension
	IDSetPlayerGameType
	IDPlayerList
	IDSimpleEvent
	IDEvent
	IDSpawnExperienceOrb
	IDClientBoundMapItemData
	IDMapInfoRequest
	IDRequestChunkRadius
	IDChunkRadiusUpdated
	IDItemFrameDropItem
	IDGameRulesChanged
	IDCamera
	IDBossEvent
	IDShowCredits
	IDAvailableCommands
	IDCommandRequest
	IDCommandBlockUpdate
	IDCommandOutput
	IDUpdateTrade
	IDUpdateEquip
	IDResourcePackDataInfo
	IDResourcePackChunkData
	IDResourcePackChunkRequest
	IDTransfer
	IDPlaySound
	IDStopSound
	IDSetTitle
	IDAddBehaviourTree
	IDStructureBlockUpdate
	IDShowStoreOffer
	IDPurchaseReceipt
	IDPlayerSkin
	IDSubClientLogin
	IDAutomationClientConnect
	IDSetLastHurtBy
	IDBookEdit
	IDNPCRequest
	IDPhotoTransfer
	IDModalFormRequest
	IDModalFormResponse
	IDServerSettingsRequest
	IDServerSettingsResponse
	IDShowProfile
	IDSetDefaultGameType
	IDRemoveObjective
	IDSetDisplayObjective
	IDSetScore
	IDLabTable
	IDUpdateBlockSynced
	IDMoveEntityDelta
	IDSetScoreboardIdentity
	IDSetLocalPlayerAsInitialised
	IDUpdateSoftEnum
	IDNetworkStackLatency
	_
	IDScriptCustomEvent
	IDSpawnParticleEffect
	IDAvailableEntityIdentifiers
	_ // IDLevelSoundEvent(2): We don't bother implementing this.
	IDNetworkChunkPublisherUpdate
)
